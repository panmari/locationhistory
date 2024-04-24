// A utility that converts a location history export from takeout to a chart.
package main

import (
	"flag"
	"log"

	"github.com/golang/geo/s2"
	"github.com/panmari/locationhistory/internal/reader"
)

var (
	input = flag.String("input", "", "Input file from google Takeout, either .zip or .json")
	loc1  = s2.LatLngFromDegrees(46.9570768, 7.4339792)
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS = 6371
)

func main() {
	flag.Parse()

	r, err := reader.OpenFile(*input)
	if err != nil {
		log.Fatalf("Error when reading %s: %v", *input, err)
	}
	res, err := reader.DecodeJson(r)
	if err != nil {
		log.Fatalf("Error when decoding %s: %v", *input, err)
	}

	for _, loc := range res {
		latlng := s2.LatLngFromDegrees(float64(loc.LatitudeE7)/1e7, float64(loc.LongitudeE7)/1e7)
		dist := latlng.Distance(loc1) * EARTH_RADIUS
	}

	log.Default().Print(res)
}
