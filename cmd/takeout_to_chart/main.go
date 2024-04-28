// A utility that converts a location history export from takeout to a chart.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/golang/geo/s2"
	"github.com/panmari/locationhistory/internal/processor"
	"github.com/panmari/locationhistory/internal/reader"
	"github.com/panmari/locationhistory/internal/visualizer"
)

var (
	input = flag.String("input", "", "Input file from google Takeout, either .zip or .json")
	loc1  = s2.LatLngFromDegrees(46.9570768, 7.4339792)
)

func main() {
	flag.Parse()

	r, err := reader.OpenFile(*input)
	if err != nil {
		log.Fatalf("Error when reading %s: %v", *input, err)
	}
	decoded, err := reader.DecodeJson(r)
	if err != nil {
		log.Fatalf("Error when decoding %s: %v", *input, err)
	}

	page := components.NewPage()
	for year := 2014; year < 2024; year++ {
		first, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-01-01", year))
		last, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-01-01", year+1))

		filter := reader.CreateDateFilter(first, last)
		locations, err := reader.FilterFunc(decoded, filter)
		if err != nil {
			log.Fatalf("Error when filtering for %d: %v", year, err)
		}
		maxDist, _ := processor.DailyDistance(loc1, locations, math.Max)

		for i, res := range [][]processor.DistanceByBucket{maxDist} {
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
	filename := fmt.Sprintf("bar.html")
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filename, err)
	}
	if page.Render(f); err != nil {
		log.Fatalf("Error writing rendering for file %s: %v", filename, err)
	}

}
