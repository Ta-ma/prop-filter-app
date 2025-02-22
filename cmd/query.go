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
)

type Property struct {
	squareFootage float32
	lighting      string
	price         float32
	rooms         int
	bathrooms     int
	location      [2]float64
	description   string
	ammenities    map[string]bool
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the available properties data using specific parameters and operators.",
	Long: `Queries the properties data using the parameters and operators passed as arguments.
The resulting filtered data will be printed in a table, limited by the max amount of entries
per page.

Example: prop-filter-app query -p 10 -n 2`,
	Run: func(cmd *cobra.Command, args []string) {

		propertiesData := []Property{
			{
				squareFootage: 500,
				lighting:      "low",
				price:         600000,
				rooms:         6,
				bathrooms:     2,
				location:      [2]float64{150, 250},
				description:   "Ample place",
				ammenities:    map[string]bool{"garage": true, "yard": true, "pool": true},
			},
			{
				squareFootage: 300,
				lighting:      "high",
				price:         450700,
				rooms:         4,
				bathrooms:     1,
				location:      [2]float64{300, 800},
				description:   "Comfy",
				ammenities:    map[string]bool{"garage": false, "yard": true, "pool": false},
			},
			{
				squareFootage: 200,
				lighting:      "low",
				price:         300000,
				rooms:         3,
				bathrooms:     1,
				location:      [2]float64{65.9, 75.7},
				description:   "Haunted",
				ammenities:    map[string]bool{"garage": false, "yard": false, "pool": false},
			},
			{
				squareFootage: 675,
				lighting:      "low",
				price:         78050.5,
				rooms:         3,
				bathrooms:     1,
				location:      [2]float64{500.2, 600},
				description:   "Nice place",
				ammenities:    map[string]bool{"garage": true, "yard": false, "pool": true},
			},
			{
				squareFootage: 333,
				lighting:      "low",
				price:         190532.976,
				rooms:         3,
				bathrooms:     1,
				location:      [2]float64{90, 40},
				description:   "Could be better",
				ammenities:    map[string]bool{"garage": true, "yard": true, "pool": false},
			},
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

		startQueryLoop(pageNumber, pageSize, maxPage, propertiesData)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().IntP("page-size", "p", 20, "page size, value of 0 will retrieve all filtered entries")
	queryCmd.Flags().IntP("page", "n", 1, "page number, will display entries of that specified page, amount of pages depends on page-size")
}

func calculateLimits(pageNumber int, pageSize int, maxLen int) (int, int) {
	lower := (pageNumber - 1) * pageSize
	upper := lower + pageSize
	if upper > maxLen {
		upper = maxLen
	}

	return lower, upper
}

func printTable(propertiesData []Property, lowerLimit int, upperLimit int) {
	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
	fmt.Fprintf(tw, "Description\tPrice\tSquare Footage\tRooms\tBathrooms\tLighting\tLocation\tAmmenities\n")
	fmt.Fprintf(tw, "-----\t-----\t-----\t-----\t-----\t-----\t-----\t-----\t\n")

	for i := lowerLimit; i < upperLimit; i++ {
		p := propertiesData[i]

		ammenities := ""
		for k, v := range p.ammenities {
			if v {
				ammenities += k + ", "
			}
		}
		ammenities = strings.TrimSuffix(ammenities, ", ")

		fmt.Fprintf(tw, "%s\t%.2f\t%.2f\t%d\t%d\t%s\t(%.2f, %.2f)\t%s\n", p.description, p.price,
			p.squareFootage, p.rooms, p.bathrooms, p.lighting, p.location[0], p.location[1], ammenities)
	}

	tw.Flush()
}

func startQueryLoop(startPageNumber int, pageSize int, maxPage int, properties []Property) {
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
