// A utility that converts a location history export from takeout to a chart.
package main

import (
	"flag"
	"log"
	"math"
	"os"

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
	res, err := processor.DailyDistance(loc1, decoded, math.Min)
	if err != nil {
		log.Fatalf("Error when processing %s: %v", *input, err)
	}

	log.Default().Print(res)

	bar := visualizer.BarChart(res)
	f, err := os.Create("bar.html")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	if bar.Render(f); err != nil {
		log.Fatalf("Error writing rendering: %v", err)
	}
}
