package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/nikolaydubina/calendarheatmap/charts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func getQuarterStartDate(currentDate time.Time) time.Time {
	year, month, _ := currentDate.Date()
	quarterStartMonth := time.Month(((int(month)-1)/3)*3 + 1)
	return time.Date(year, quarterStartMonth, 1, 0, 0, 0, 0, currentDate.Location())
}

func getQuarterEndDate(currentDate time.Time) time.Time {
	year, month, _ := currentDate.Date()
	quarterStartMonth := time.Month(((int(month)-1)/3)*3 + 1)
	quarterEndMonth := quarterStartMonth + 2 // The end month is three months after the start month.
	quarterEndMonth %= 12                    // Ensure the month doesn't exceed 12 (December).
	if quarterEndMonth < quarterStartMonth {
		year++
	}
	lastDayOfQuarter := time.Date(year, quarterEndMonth+1, 0, 0, 0, 0, 0, currentDate.Location())
	return lastDayOfQuarter
}

//go:embed assets/fonts/Sunflower-Medium.ttf
var defaultFontFaceBytes []byte

//go:embed assets/colorscales/incidents.csv
var defaultColorScaleBytes []byte

// Taken from https://github.com/nikolaydubina/calendarheatmap/blob/master/main.go.
func plotHeatmap(data []byte, outputFormat string) {

	labels := true
	monthSep := true
	colorScale := "incidents.csv"
	locale := "en_US"

	var colorscale charts.BasicColorScale
	if assetsPath := os.Getenv("CALENDAR_HEATMAP_ASSETS_PATH"); assetsPath != "" {
		var err error
		colorscale, err = charts.NewBasicColorscaleFromCSVFile(path.Join(assetsPath, "colorscales", colorScale))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		if colorScale != "incidents.csv" {
			log.Printf("defaulting to colorscale %s since CALENDAR_HEATMAP_ASSETS_PATH is not set", "incidents.csv")
		}
		colorscale, err = charts.NewBasicColorscaleFromCSV(bytes.NewBuffer(defaultColorScaleBytes))
		if err != nil {
			log.Fatal(err)
		}
	}

	fontFace, err := charts.LoadFontFace(defaultFontFaceBytes, opentype.FaceOptions{
		Size:    36,
		DPI:     280,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	var counts map[string]int
	if err := json.Unmarshal(data, &counts); err != nil {
		log.Fatal(err)
	}
	conf := charts.HeatmapConfig{
		Counts:              counts,
		ColorScale:          colorscale,
		DrawMonthSeparator:  monthSep,
		DrawLabels:          labels,
		Margin:              90,
		BoxSize:             350,
		MonthSeparatorWidth: 15,
		MonthLabelYOffset:   50,
		TextWidthLeft:       300,
		TextHeightTop:       200,
		TextColor:           color.RGBA{100, 100, 100, 255},
		BorderColor:         color.RGBA{200, 200, 200, 255},
		Locale:              locale,
		Format:              outputFormat,
		FontFace:            fontFace,
		ShowWeekdays: map[time.Weekday]bool{
			time.Monday:    true,
			time.Wednesday: true,
			time.Friday:    true,
		},
	}
	charts.WriteHeatmap(conf, os.Stdout)
}

var (
	statusPageEndpoint = flag.String("endpoint", "https://status.verity.gumgum.com", "Status page endpoint to scrape incidents from")
	fromFlag           = flag.String("from", getQuarterStartDate(time.Now()).Format("2006-01-02"), "Releases after this date will be included")
	toFlag             = flag.String("to", getQuarterEndDate(time.Now()).Format("2006-01-02"), "Releases before this date will be included")
	outputFlag         = flag.String("output", "png", "Output format (json, png, jpeg, gif, svg)")
)

type ApiResponse struct {
	Incidents []Incident `json:"incidents"`
}

type Incident struct {
	CreatedAt  time.Time `json:"created_at"`
	ResolvedAt time.Time `json:"resolved_at"`
}

func fetchAndProcessIncidents(url string, startDate, endDate time.Time) (map[string]int, error) {
	// Fetch the data
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create the map
	dateDurationMap := make(map[string]int)

	for _, incident := range apiResponse.Incidents {
		// Filter incidents by date range
		if incident.CreatedAt.After(startDate) && incident.CreatedAt.Before(endDate.AddDate(0, 0, 1)) {
			date := incident.CreatedAt.Format("2006-01-02")
			duration := int(incident.ResolvedAt.Sub(incident.CreatedAt).Minutes())

			dateDurationMap[date] += duration
		}
	}

	return dateDurationMap, nil
}

func main() {
	// Parse flags
	flag.Parse()

	// find all releases within date range
	fromDate, _ := time.Parse("2006-01-02", *fromFlag)
	toDate, _ := time.Parse("2006-01-02", *toFlag)
	url := fmt.Sprintf("%s/api/v2/incidents.json", *statusPageEndpoint)

	// Fetch the data and process it
	incidentsMap, err := fetchAndProcessIncidents(url, fromDate, toDate)
	if err != nil {
		fmt.Printf("Failed to fetch or process data: %v\n", err)
		return
	}

	// Print the result
	incidentsMapJson, _ := json.MarshalIndent(incidentsMap, "", "")

	if *outputFlag == "json" {
		fmt.Println(string(incidentsMapJson))
	} else {
		// Plot heatmap
		plotHeatmap(incidentsMapJson, *outputFlag)
	}
}
