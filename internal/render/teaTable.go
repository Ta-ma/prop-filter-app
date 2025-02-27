/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ta-ma/prop-filter-app/internal/db"
	"github.com/ta-ma/prop-filter-app/internal/models"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table        table.Model
	currentPage  int
	maxPage      int
	pageHeight   int
	queryFilter  string
	calcDistance bool
	distX        string
	distY        string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var tableChanged bool
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.currentPage > 1 {
				m.currentPage--
			}
			tableChanged = true
		case "right":
			if m.currentPage < m.maxPage {
				m.currentPage++
			}
			tableChanged = true
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	if tableChanged {
		rows, err := getTableRows(m.queryFilter, m.pageHeight, m.currentPage, m.calcDistance, m.distX, m.distY)
		if err != nil {
			panic("Error while rebuilding the table! Exiting...")
		}
		m.table.SetRows(rows)
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.getDetails() +
		baseStyle.Render(m.table.View()) + "\n" +
		m.getPageInfo() + "\n" +
		m.getKeysInfo() + "\n"
}

func (m model) getDetails() string {
	var lines []string
	row := m.table.SelectedRow()
	columns := []string{
		"Description", "Price", "Square ft", "Rooms", "Bathrooms", "Lighting", "Location",
	}

	if m.calcDistance {
		columns = append(columns, "Distance")
	}
	columns = append(columns, "Ammenities")

	for i, c := range columns {
		lines = append(lines, fmt.Sprintf("%s: %s", c, row[i]))
	}

	return lipgloss.NewStyle().
		Padding(1).
		Render(strings.Join(lines, "\n")) + "\n"
}

func (m model) getKeysInfo() string {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Render("Up/Down: Move selection   Left/Right: Change page   Q: Exit")
}

func (m model) getPageInfo() string {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Page %d / %d\n", m.currentPage, m.maxPage))
}

func ShowTeaTable(startPageNumber int, pageHeight int, maxPage int, queryFilter string, calcDistance bool, distX string, distY string) {
	rows, err := getTableRows(queryFilter, pageHeight, startPageNumber, calcDistance, distX, distY)
	if err != nil {
		return
	}

	columns := []table.Column{
		{Title: "Description", Width: 30},
		{Title: "Price", Width: 10},
		{Title: "Square ft", Width: 10},
		{Title: "Rooms", Width: 6},
		{Title: "Bathrooms", Width: 10},
		{Title: "Lighting", Width: 10},
		{Title: "Location", Width: 16},
	}

	if calcDistance {
		columns = append(columns, table.Column{Title: "Distance", Width: 10})
	}
	columns = append(columns, table.Column{Title: "Ammenities", Width: 20})

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(pageHeight+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t, currentPage: startPageNumber, maxPage: maxPage, pageHeight: pageHeight,
		queryFilter: queryFilter, calcDistance: calcDistance, distX: distX, distY: distY}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error displaying table:", err)
		return
	}
}

func mapPropertiesToRows(results []models.PropertyViewModel, calcDistance bool) []table.Row {
	var rows []table.Row

	for _, r := range results {
		location := fmt.Sprintf("(%.2f,%.2f)", r.Latitude, r.Longitude)
		price := fmt.Sprintf("$%.2f", r.Price)
		sqft := fmt.Sprintf("%.2f", r.Square_footage)
		rooms := fmt.Sprintf("%d", r.Rooms)
		bathrooms := fmt.Sprintf("%d", r.Bathrooms)
		distance := fmt.Sprintf("%.2f", r.Dist)

		row := table.Row{
			r.Description, price, sqft, rooms, bathrooms, r.Lighting, location,
		}

		if calcDistance {
			row = append(row, distance)
		}
		row = append(row, r.Ammenities)

		rows = append(rows, row)
	}

	return rows
}

func getTableRows(queryFilter string, pageHeight int, pageNumber int, calcDistance bool, distX string, distY string) ([]table.Row, error) {
	props, err := db.QueryProperties(queryFilter, pageHeight, (pageNumber-1)*pageHeight, calcDistance, distX, distY)
	if err != nil {
		return []table.Row{}, err
	}
	rows := mapPropertiesToRows(props, calcDistance)
	return rows, nil
}
