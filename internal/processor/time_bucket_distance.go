package processor

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/golang/geo/s2"
	"github.com/panmari/locationhistory/internal/reader"
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS_KM = 6371
)

func bucketTimestamp(bucketDuration time.Duration, location reader.Location) (time.Time, error) {
	// TODO(panmari): Accept offset to better handle timezones in bucketing.
	ts, err := location.ParsedTimestamp()
	if err != nil {
		return time.Time{}, err
	}
	return ts.Round(bucketDuration), nil
}

type DistanceByTimeBucket struct {
	Distance float64
	Bucket   time.Time
}

func (d DistanceByTimeBucket) String() string {
	return fmt.Sprintf("Dist: %f, Bucket: %s", d.Distance, d.Bucket.Format(time.DateTime))
}

// TimeBucketDistance measures the distance of each data point to the anchor location for each duration and reduces
// it to a single value using the given reducer fuction.
func TimeBucketDistance(anchors []Anchor, locations []reader.Location, bucketDuration time.Duration, reducer func(a, b float64) float64) ([]DistanceByTimeBucket, error) {
	minDistanceByDate := make(map[time.Time]float64, 365)

	for _, loc := range locations {
		ts, err := bucketTimestamp(bucketDuration, loc)
		if err != nil {
			log.Default().Println(err)
			continue
		}
		// TODO(panmari): Consider validating that ts is not before StartTime.
		for len(anchors) > 1 && ts.After(anchors[1].StartTime) {
			anchors = anchors[1:]
		}
		latlng := s2.LatLngFromDegrees(float64(loc.LatitudeE7)/1e7, float64(loc.LongitudeE7)/1e7)
		dist := float64(latlng.Distance(anchors[0].Location)) * EARTH_RADIUS_KM
		d, ok := minDistanceByDate[ts]
		if !ok {
			minDistanceByDate[ts] = dist
			continue
		}
		minDistanceByDate[ts] = reducer(d, dist)
	}
	res := make([]DistanceByTimeBucket, 0, len(minDistanceByDate))
	for date, distance := range minDistanceByDate {
		res = append(res, DistanceByTimeBucket{
			Distance: distance,
			Bucket:   date,
		})
	}
	slices.SortFunc(res, func(a, b DistanceByTimeBucket) int {
		return a.Bucket.Compare(b.Bucket)
	})
	return res, nil
}
