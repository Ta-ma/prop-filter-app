/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
	"github.com/ta-ma/prop-filter-app/internal/db"
	"github.com/ta-ma/prop-filter-app/internal/filter"
	"github.com/ta-ma/prop-filter-app/internal/models"
	"github.com/ta-ma/prop-filter-app/internal/render"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the available properties data using specific parameters and operators.",
	Long: `Queries the properties data using the parameters and operators passed as arguments.
The resulting filtered data will be printed in a table, limited by the max amount of entries
per page.

Example: prop-filter-app query -w 10 -n 2 -p "<700000"`,
	Run: func(cmd *cobra.Command, args []string) {
		pageHeight, _ := cmd.Flags().GetInt("page-size")
		pageNumber, _ := cmd.Flags().GetInt("page")
		priceExpr, _ := cmd.Flags().GetString("price")
		roomsExpr, _ := cmd.Flags().GetString("rooms")
		bathroomsExpr, _ := cmd.Flags().GetString("bathrooms")
		latitudeExpr, _ := cmd.Flags().GetString("latitude")
		longitudeExpr, _ := cmd.Flags().GetString("longitude")
		sqftExpr, _ := cmd.Flags().GetString("sqft")
		descExpr, _ := cmd.Flags().GetString("description")
		amenitiesExpr, _ := cmd.Flags().GetString("amenities")
		lightingExpr, _ := cmd.Flags().GetString("lighting")
		distanceExpr, _ := cmd.Flags().GetString("distance")

		translator := filter.Translator{}
		translator.Init()
		translator.Translate("p.price", priceExpr, filter.Num)
		translator.Translate("p.rooms", roomsExpr, filter.Num)
		translator.Translate("p.bathrooms", bathroomsExpr, filter.Num)
		translator.Translate("p.latitude", latitudeExpr, filter.Num)
		translator.Translate("p.longitude", longitudeExpr, filter.Num)
		translator.Translate("p.square_footage", sqftExpr, filter.Num)
		translator.Translate("p.description", descExpr, filter.Str)
		translator.Translate("l.description", lightingExpr, filter.Lighting)
		translator.Translate("a.description", amenitiesExpr, filter.Amenity)

		var calcDistance bool
		var distanceData filter.DistanceFilterData
		if distanceExpr != "" {
			calcDistance = true
			distanceData = translator.TranslateDistanceExpr("d.dist", distanceExpr)
		}

		if translator.Err != nil {
			fmt.Println("Failed to parse filter parameters:", translator.Err)
			return
		}
		sqlFilter := translator.GetSqlTranslation()

		propsCount, err := db.GetPropertiesCount(sqlFilter, calcDistance, distanceData.X, distanceData.Y)
		if err != nil {
			fmt.Println("Properties could not be counted:", err)
			return
		}

		if propsCount == 0 {
			fmt.Println("There is no properties data available to display.")
			return
		}

		if pageHeight < 1 {
			fmt.Println("ERROR: Page size parameter should be 1 or greater.")
			return
		}

		maxPage := propsCount / pageHeight
		if propsCount%pageHeight > 0 {
			maxPage++
		}

		if pageNumber > maxPage || pageNumber < 1 {
			fmt.Println("ERROR: Page number must be at least 1 and lesser than the max amount of pages given the page size.")
			return
		}

		if cfg.UseOldRender {
			startLoop(pageNumber, pageHeight, maxPage, sqlFilter, calcDistance, distanceData.X, distanceData.Y)
		} else {
			render.ShowTeaTable(pageNumber, pageHeight, maxPage, sqlFilter, calcDistance, distanceData.X, distanceData.Y)
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().IntP("page-size", "w", 15, "page size, amount of entries listed in each page")
	queryCmd.Flags().IntP("page", "n", 1, "page number, will display entries of that specified page, amount of pages depends on page-size")
	queryCmd.Flags().StringP("price", "p", "", "Expression to filter entries by the Price field")
	queryCmd.Flags().StringP("rooms", "r", "", "Expression to filter entries by the Rooms field")
	queryCmd.Flags().StringP("bathrooms", "b", "", "Expression to filter entries by the Bathrooms field")
	queryCmd.Flags().StringP("latitude", "x", "", "Expression to filter entries by the Latitude field")
	queryCmd.Flags().StringP("longitude", "y", "", "Expression to filter entries by the Longitude field")
	queryCmd.Flags().StringP("sqft", "s", "", "Expression to filter entries by the Square ft field")
	queryCmd.Flags().StringP("description", "d", "", "Expression to filter entries by the Description field")
	queryCmd.Flags().StringP("amenities", "a", "", "Expression to filter entries by the Amenities field")
	queryCmd.Flags().StringP("lighting", "l", "", "Expression to filter entries by the Lighting field")
	queryCmd.Flags().StringP("distance", "k", "", "Expression to filter entries by the Description field")
}

func printTable(result []models.PropertyViewModel, calcDist bool) {
	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	if calcDist {
		fmt.Fprintf(tw, "Description\tPrice\tSquare ft\tRooms\tBathrooms\tLighting\tLocation\tDistance\tAmenities\n")
		fmt.Fprintf(tw, "-----\t-----\t-----\t-----\t-----\t-----\t-----\t-----\t-----\t\n")
		for _, r := range result {
			fmt.Fprintf(tw, "%s\t%.2f\t%.2f\t%d\t%d\t%s\t(%.2f, %.2f)\t%.2f\t%s\n", trimString(r.Description),
				r.Price, r.Square_footage, r.Rooms, r.Bathrooms, r.Lighting, r.Latitude, r.Longitude,
				r.Dist, r.Amenities)
		}
	} else {
		fmt.Fprintf(tw, "Description\tPrice\tSquare ft\tRooms\tBathrooms\tLighting\tLocation\tAmenities\n")
		fmt.Fprintf(tw, "-----\t-----\t-----\t-----\t-----\t-----\t-----\t-----\t\n")
		for _, r := range result {
			fmt.Fprintf(tw, "%s\t%.2f\t%.2f\t%d\t%d\t%s\t(%.2f, %.2f)\t%s\n", trimString(r.Description),
				r.Price, r.Square_footage, r.Rooms, r.Bathrooms, r.Lighting, r.Latitude,
				r.Longitude, r.Amenities)
		}
	}

	tw.Flush()
}

func startLoop(startPageNumber int, pageHeight int, maxPage int, queryFilter string, calcDistance bool, distX string, distY string) {
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
			properties, err := db.QueryProperties(queryFilter, pageHeight, (pageNumber-1)*pageHeight, calcDistance, distX, distY)
			if err != nil {
				fmt.Println("Properties could not be queried:", err)
				return
			}

			fmt.Println()
			printTable(properties, calcDistance)
			fmt.Println()
			fmt.Printf("Page %d / %d\n", pageNumber, maxPage)
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
