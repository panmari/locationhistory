// A utility that converts a location history export from takeout to a chart.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/panmari/locationhistory/internal/processor"
	"github.com/panmari/locationhistory/internal/reader"
	"github.com/panmari/locationhistory/internal/visualizer"
)

var (
	input         = flag.String("input", "", "Input file from google Takeout, either .zip or .json")
	anchorsString = flag.String("anchors", "", "Anchor location which are used to compute distance. either in the format lat,lng or date,lat,lng:date2,lat,lng")
	timeZone      = flag.String("timezone", "", "Time zone to use for displaying data")
)

func yearlyCharts(anchors []processor.Anchor, decoded []reader.Location) *components.Page {
	page := components.NewPage()
	page.PageTitle = "Yearly plots from timeline"
	bucketOpts := processor.Options{Anchors: anchors, BucketDuration: time.Hour * 24, Reducer: math.Max}
	for year := 2014; year < 2024; year++ {
		first, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-01-01", year))
		last, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-12-31", year))

		filter := reader.CreateDateFilter(first, last)
		locations, err := reader.FilterFunc(decoded, filter)
		if err != nil {
			log.Fatalf("Error when filtering for %d: %v", year, err)
		}
		if len(locations) == 0 {
			continue
		}
		maxDist, err := processor.TimeBucketDistance(locations, bucketOpts)
		if err != nil {
			log.Fatalf("Error when bucketing for %d: %v", year, err)
		}
		for i, res := range [][]processor.DistanceByTimeBucket{maxDist} {
			bar := visualizer.BarChart(res)
			bar.SetGlobalOptions(
				charts.WithTitleOpts(opts.Title{
					Title: fmt.Sprintf("Year: %d, %d", year, i),
				}),
			)
			page.AddCharts(bar)
			page.AddCharts(visualizer.Heatmap(res))
		}
	}
	return page
}

func dailyCharts(anchors []processor.Anchor, decoded []reader.Location) *components.Page {
	page := components.NewPage().SetLayout(components.PageFlexLayout)
	page.PageTitle = "Daily plots from timeline"
	// TODO(panmari): Move concept of timezone to anchor, so far moves are easier to account for.
	tz, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		log.Fatalf("Error when parsing time zone %q: %v", *timeZone, err)
	}
	bucketOpts := processor.Options{Anchors: anchors, BucketDuration: time.Hour, Reducer: math.Max}
	for year := 2014; year < 2024; year++ {
		first, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-01-01", year))
		last, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-12-31", year))

		filter := reader.CreateDateFilter(first, last)
		locations, err := reader.FilterFunc(decoded, filter)
		if err != nil {
			log.Fatalf("Error when filtering for %d: %v", year, err)
		}
		if len(locations) == 0 {
			continue
		}
		maxDist, err := processor.TimeBucketDistance(locations, bucketOpts)
		if err != nil {
			log.Fatalf("Error when bucketing for %d: %v", year, err)
		}
		radars := visualizer.DailyRadar(maxDist, visualizer.Options{Title: fmt.Sprintf("Year %d", year), TimeZone: tz})
		page.AddCharts(radars...)

	}
	return page
}

func main() {
	flag.Parse()

	anchors, err := processor.ParseAnchors(*anchorsString)
	if err != nil {
		log.Fatalf("Error parsing --anchors argument %q: %v", *anchorsString, err)
	}

	r, err := reader.OpenFile(*input)
	if err != nil {
		log.Fatalf("Error when reading %s: %v", *input, err)
	}
	decoded, err := reader.DecodeJson(r)
	if err != nil {
		log.Fatalf("Error when decoding %s: %v", *input, err)
	}

	yearlyPage := yearlyCharts(anchors, decoded)
	filename := fmt.Sprintf("yearly.html")
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filename, err)
	}
	if yearlyPage.Render(f); err != nil {
		log.Fatalf("Error writing rendering for file %s: %v", filename, err)
	}

	dailyPage := dailyCharts(anchors, decoded)
	filename = fmt.Sprintf("daily.html")
	f, err = os.Create(filename)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filename, err)
	}
	if dailyPage.Render(f); err != nil {
		log.Fatalf("Error writing rendering for file %s: %v", filename, err)
	}

}
