package processor

import (
	"log"
	"math"
	"time"

	"github.com/golang/geo/s2"
	"github.com/panmari/locationhistory/internal/reader"
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS_KM = 6371
)

func bucketTimestamp(location reader.Location) (time.Time, error) {
	ts, err := location.ParsedTimestamp()
	if err != nil {
		return time.Time{}, err
	}
	return ts.Round(time.Hour * 24), nil
}

// DailyMinimumDistance measures for each day in the data the minimum distance to the anchor location.
func DailyMinimumDistance(anchor s2.LatLng, locations []reader.Location) ([]float64, error) {
	dist := math.MaxFloat64
	lastBucketTime, err := bucketTimestamp(locations[0])
	if err != nil {
		return nil, err
	}
	res := make([]float64, 0, 0)
	for _, loc := range locations {
		latlng := s2.LatLngFromDegrees(float64(loc.LatitudeE7)/1e7, float64(loc.LongitudeE7)/1e7)
		dist = math.Min(dist, float64(latlng.Distance(anchor))*EARTH_RADIUS_KM)
		// TODO: some smarter bucketing here
		if ts, err := bucketTimestamp(loc); err == nil {
			if ts != lastBucketTime {
				res = append(res, dist)
				dist = math.MaxFloat64
				lastBucketTime = ts
			}
		} else {
			log.Default().Println(err)
		}
	}
	return res, nil
}
