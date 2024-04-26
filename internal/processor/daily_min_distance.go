package processor

import (
	"log"
	"slices"
	"strings"
	"time"

	"github.com/golang/geo/s2"
	"github.com/panmari/locationhistory/internal/reader"
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS_KM = 6371
)

func bucketTimestamp(location reader.Location) (time.Time, error) {
	// TODO(panmari): Accept offset to better handle timezones in bucketing.
	ts, err := location.ParsedTimestamp()
	if err != nil {
		return time.Time{}, err
	}
	return ts.Round(time.Hour * 24), nil
}

type DistanceByBucket struct {
	Distance float64
	Bucket   string
}

// DailyMinimumDistance measures for each day in the data the minimum distance to the anchor location.
func DailyMinimumDistance(anchor s2.LatLng, locations []reader.Location) ([]DistanceByBucket, error) {
	minDistanceByDate := make(map[time.Time]float64, 365)
	for _, loc := range locations {
		latlng := s2.LatLngFromDegrees(float64(loc.LatitudeE7)/1e7, float64(loc.LongitudeE7)/1e7)
		dist := float64(latlng.Distance(anchor)) * EARTH_RADIUS_KM
		ts, err := bucketTimestamp(loc)
		if err != nil {
			log.Default().Println(err)
			continue
		}
		if d, ok := minDistanceByDate[ts]; !ok || dist < d {
			minDistanceByDate[ts] = dist
		}
	}
	res := make([]DistanceByBucket, 0, len(minDistanceByDate))
	for date, distance := range minDistanceByDate {
		res = append(res, DistanceByBucket{
			Distance: distance,
			Bucket:   date.Format(time.DateOnly),
		})
	}
	slices.SortFunc(res, func(a, b DistanceByBucket) int {
		return strings.Compare(a.Bucket, b.Bucket)
	})
	return res, nil
}
