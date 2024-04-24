// A utility that converts a location history export from takeout to a chart.
package main

import (
	"flag"
	"log"

	"github.com/panmari/locationhistory/internal/reader"
)

var (
	input = flag.String("input", "", "Input file from google Takeout, either .zip or .json")
)

func main() {
	flag.Parse()

	r, err := reader.OpenFile(*input)
	if err != nil {
		log.Fatalf("Error when reading %s: %v", *input, err)
	}
	log.Default().Print(reader.DecodeJson(r))
}
