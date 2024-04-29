package processor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/geo/s2"
)

// Anchor is a location that is used for computing distances, starting at a given date.
type Anchor struct {
	StartTime time.Time
	Location  s2.LatLng
}

// ParseAnchors parses a string as anchors. The required format is
// date,lat,lng:date2,lat2,lng2
func ParseAnchors(anchors string) ([]Anchor, error) {
	multidayAnchors := strings.Split(anchors, ":")
	if len(multidayAnchors) == 1 {
		var lat, lng float64
		_, err := fmt.Sscanf(multidayAnchors[0], "%f,%f", &lat, &lng)
		return []Anchor{{StartTime: time.Time{}, Location: s2.LatLngFromDegrees(lat, lng)}}, err
	}
	res := make([]Anchor, 0, len(multidayAnchors))
	for _, s := range multidayAnchors {
		split := strings.Split(s, ",")
		if len(split) != 3 {
			return nil, fmt.Errorf("dated anchor %q does not contain two commas", s)
		}
		t, err := time.Parse(time.DateOnly, split[0])
		if err != nil {
			return nil, fmt.Errorf("failed parsing date from %q: %w", s, err)
		}
		lat, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing lat from %q: %w", s, err)
		}
		lng, err := strconv.ParseFloat(split[2], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing lng from %q: %w", s, err)
		}
		res = append(res, Anchor{
			StartTime: t,
			Location:  s2.LatLngFromDegrees(lat, lng),
		})
	}
	return res, nil
}
