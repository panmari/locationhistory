package reader

import (
	"strings"
	"testing"
)

// Run using
//
// go test -fuzz=FuzzDecodeJson
func FuzzDecodeJson(f *testing.F) {
	testcases := []string{"Hello, world", " ", `
	{
		"locations": [{
		"latitudeE7": 123,
		"longitudeE7": 456,
		"accuracy": 19,
		"source": "WIFI",
		"timestamp": "2014-01-01T00:00:50.307Z"
	  },
	`}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, input string) {
		r := strings.NewReader(input)
		// There's no good way to validate the output, so just not segfaulting is good enough.
		DecodeJson(r)
	})
}
