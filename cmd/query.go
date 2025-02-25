/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
	"github.com/ta-ma/prop-filter-app/internal/db"
	"github.com/ta-ma/prop-filter-app/internal/models"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the available properties data using specific parameters and operators.",
	Long: `Queries the properties data using the parameters and operators passed as arguments.
The resulting filtered data will be printed in a table, limited by the max amount of entries
per page.

Example: prop-filter-app query -p 10 -n 2`,
	Run: func(cmd *cobra.Command, args []string) {
		filterExpr, _ := cmd.Flags().GetString("filter")
		propertiesData, err := db.QueryProperties(filterExpr)
		if err != nil {
			fmt.Println("Properties could not be queried:", err)
			return
		}

		pageSize, _ := cmd.Flags().GetInt("page-size")
		pageNumber, _ := cmd.Flags().GetInt("page")
		lenData := len(propertiesData)

		if lenData == 0 {
			fmt.Println("There is no properties data available to display.")
			return
		}

		if pageSize < 1 {
			fmt.Println("ERROR: Page size parameter should be 1 or greater.")
			return
		}

		maxPage := lenData / pageSize
		if lenData%pageSize > 0 {
			maxPage++
		}

		if pageNumber > maxPage || pageNumber < 1 {
			fmt.Println("ERROR: Page number must be at least 1 and lesser than the max amount of pages given the page size.")
			return
		}

		startLoop(pageNumber, pageSize, maxPage, propertiesData)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().IntP("page-size", "p", 20, "page size, value of 0 will retrieve all filtered entries")
	queryCmd.Flags().IntP("page", "n", 1, "page number, will display entries of that specified page, amount of pages depends on page-size")
	queryCmd.Flags().StringP("filter", "f", "", "conditional SQL expression which will be used to filter queried properties")
}

func calculateLimits(pageNumber int, pageSize int, maxLen int) (int, int) {
	lower := (pageNumber - 1) * pageSize
	upper := lower + pageSize
	if upper > maxLen {
		upper = maxLen
	}

	return lower, upper
}

func printTable(propertiesData []models.Property, lowerLimit int, upperLimit int) {
	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintf(tw, "Description\tPrice\tSquare Footage\tRooms\tBathrooms\tLighting\tLocation\tAmmenities\n")
	fmt.Fprintf(tw, "-----\t-----\t-----\t-----\t-----\t-----\t-----\t-----\t\n")

	for i := lowerLimit; i < upperLimit; i++ {
		p := propertiesData[i]

		ammenities := ""
		for _, a := range p.Ammenities {
			ammenities += a.Description + ", "
		}
		ammenities = strings.TrimSuffix(ammenities, ", ")

		fmt.Fprintf(tw, "%s\t%.2f\t%.2f\t%d\t%d\t%s\t(%.2f, %.2f)\t%s\n", trimString(p.Description),
			p.Price, p.SquareFootage, p.Rooms, p.Bathrooms, p.Lighting.Description, p.Latitude,
			p.Longitude, trimString(ammenities))
	}

	tw.Flush()
}

func startLoop(startPageNumber int, pageSize int, maxPage int, properties []models.Property) {
	pageNumber := startPageNumber
	invalidKeyPressed := false

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		if !invalidKeyPressed {
			lowerLimit, upperLimit := calculateLimits(pageNumber, pageSize, len(properties))
			printTable(properties, lowerLimit, upperLimit)

			fmt.Println()
			if pageNumber != 1 {
				fmt.Print("LeftArrow -> previous page  ")
			}
			if pageNumber != maxPage {
				fmt.Print("RightArrow -> next page  ")
			}
			fmt.Print("ESC and ENTER -> exit\n")
		} else {
			invalidKeyPressed = false
		}

		_, key, err := keyboard.GetKey()

		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			return
		} else if key == keyboard.KeyArrowLeft && pageNumber != 1 {
			pageNumber--
		} else if key == keyboard.KeyArrowRight && pageNumber != maxPage {
			pageNumber++
		} else {
			invalidKeyPressed = true
			fmt.Println("Invalid key pressed")
		}
	}
}

func trimString(str string) string {
	if cfg.TrimLength != 0 && len(str) > cfg.TrimLength {
		return str[0:cfg.TrimLength-3] + "..."
	}

	return str
}
