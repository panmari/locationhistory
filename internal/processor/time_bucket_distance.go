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

type Options struct {
	Anchors        []Anchor
	BucketDuration time.Duration
	Reducer        func(a, b float64) float64
}

// TimeBucketDistance measures the distance of each data point to the anchor location for each duration and reduces
// it to a single value using the given reducer fuction.
func TimeBucketDistance(locations []reader.Location, opts Options) ([]DistanceByTimeBucket, error) {
	minDistanceByDate := make(map[time.Time]float64, 365)

	for _, loc := range locations {
		ts, err := bucketTimestamp(opts.BucketDuration, loc)
		if err != nil {
			log.Default().Println(err)
			continue
		}
		// TODO(panmari): Consider validating that ts is not before StartTime.
		for len(opts.Anchors) > 1 && ts.After(opts.Anchors[1].StartTime) {
			opts.Anchors = opts.Anchors[1:]
		}
		latlng := s2.LatLngFromDegrees(float64(loc.LatitudeE7)/1e7, float64(loc.LongitudeE7)/1e7)
		dist := float64(latlng.Distance(opts.Anchors[0].Location)) * EARTH_RADIUS_KM
		d, ok := minDistanceByDate[ts]
		if !ok {
			minDistanceByDate[ts] = dist
			continue
		}
		minDistanceByDate[ts] = opts.Reducer(d, dist)
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
