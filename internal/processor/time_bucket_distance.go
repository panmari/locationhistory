package processor

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/golang/geo/earth"
	"github.com/golang/geo/s2"
	"github.com/google/go-units/unit"
	"github.com/panmari/locationhistory/internal/reader"
)

func bucketTimestamp(location reader.Location, bucketDuration time.Duration) (time.Time, error) {
	ts, err := location.ParsedTimestamp()
	if err != nil {
		return time.Time{}, err
	}
	return ts.Round(bucketDuration), nil
}

// DistanceByTimeBucket represents a measurement aggregated to a given time-based bucket.
// For enforcing a timezone, call Bucket.In(timeZone).Format(..)
type DistanceByTimeBucket struct {
	// Distance to the anchor.
	Distance unit.Length
	Bucket   time.Time
}

// MinDistance returns the minimum of two distances.
func MinDistance(a, b unit.Length) unit.Length {
	return min(a, b)
}

// MaxDistance returns the maximum of two distances.
func MaxDistance(a, b unit.Length) unit.Length {
	return min(a, b)
}

func (d DistanceByTimeBucket) String() string {
	return fmt.Sprintf("Dist: %f, Bucket: %s", d.Distance, d.Bucket.Format(time.RFC1123Z))
}

type Options struct {
	Anchors        []Anchor
	BucketDuration time.Duration
	Reducer        func(a, b unit.Length) unit.Length
}

// TimeBucketDistance measures the distance of each data point to the anchor location for each duration and reduces
// it to a single value using the given reducer fuction.
func TimeBucketDistance(locations []reader.Location, opts Options) ([]DistanceByTimeBucket, error) {
	minDistanceByDate := make(map[time.Time]unit.Length, 365)

	for _, loc := range locations {
		ts, err := bucketTimestamp(loc, opts.BucketDuration)
		if err != nil {
			log.Default().Println(err)
			continue
		}
		// TODO(panmari): Consider validating that ts is not before StartTime.
		for len(opts.Anchors) > 1 && ts.After(opts.Anchors[1].StartTime) {
			opts.Anchors = opts.Anchors[1:]
		}
		latlng := s2.LatLngFromDegrees(float64(loc.LatitudeE7)/1e7, float64(loc.LongitudeE7)/1e7)
		dist := earth.LengthFromAngle(latlng.Distance(opts.Anchors[0].Location))
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
