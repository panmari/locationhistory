package processor

import (
	"math"
	"testing"
	"time"

	"github.com/golang/geo/s2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/panmari/locationhistory/internal/reader"
)

func parseDate(t *testing.T, date string) time.Time {
	t.Helper()
	res, err := time.Parse(time.DateOnly, date)
	if err != nil {
		t.Error(err)
	}
	return res
}

func TestDailyDistance(t *testing.T) {
	locations := []reader.Location{
		{
			Timestamp:   "2014-04-01T07:55:51.093Z",
			LatitudeE7:  469287872,
			LongitudeE7: 74171385,
		},
		{
			Timestamp:   "2014-04-01T08:55:51.093Z",
			LatitudeE7:  469281883,
			LongitudeE7: 74156002,
		},
	}
	for _, tc := range []struct {
		name   string
		anchor []Anchor
		want   []DistanceByTimeBucket
	}{
		{
			name: "Anchor at one location gives zero",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(46.9287872, 7.4171385),
			}},
			want: []DistanceByTimeBucket{{0, parseDate(t, "2014-04-01")}},
		},
		{
			name: "Anchor far away gives non-zero",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(50, 5),
			}},
			want: []DistanceByTimeBucket{{385.159008, parseDate(t, "2014-04-01")}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := Options{Anchors: tc.anchor, BucketDuration: time.Hour * 24, Reducer: math.Min}
			got, err := TimeBucketDistance(locations, opts)
			if err != nil || !cmp.Equal(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)) {
				t.Errorf("DailyDistance() = %v, %v, want %v", got, err, tc.want)
			}
		})
	}
}

func TestDailyDistanceMultipleDays(t *testing.T) {
	locations := []reader.Location{
		{
			Timestamp:   "2014-04-01T07:55:51.093Z",
			LatitudeE7:  469287872,
			LongitudeE7: 74171385,
		},
		{
			Timestamp:   "2014-04-01T09:15:51.093Z",
			LatitudeE7:  469287872,
			LongitudeE7: 74171385,
		},
		{
			Timestamp:   "2014-05-02T15:55:51.093Z",
			LatitudeE7:  459281883,
			LongitudeE7: 74156002,
		},
	}
	for _, tc := range []struct {
		name           string
		anchor         []Anchor
		bucketDuration time.Duration
		want           []DistanceByTimeBucket
	}{
		{
			name: "Anchor at one location gives zero",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(46.9287872, 7.4171385),
			}},
			bucketDuration: time.Hour * 24,
			want:           []DistanceByTimeBucket{{0, parseDate(t, "2014-04-01")}, {111.2615837, parseDate(t, "2014-05-03")}},
		},
		{
			name: "Anchor at one location gives zero with hourly bucket",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(46.9287872, 7.4171385),
			}},
			bucketDuration: time.Hour,
			want: []DistanceByTimeBucket{
				{0, time.Date(2014, 4, 1, 8, 0, 0, 0, time.UTC)},
				{0, time.Date(2014, 4, 1, 9, 0, 0, 0, time.UTC)},
				{111.2615837, time.Date(2014, 5, 2, 16, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "Two Anchors give two times zero",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(46.9287872, 7.4171385),
			}, {
				// Expected to be skipped.
				StartTime: time.Date(2014, 04, 13, 0, 0, 0, 0, time.UTC),
				Location:  s2.LatLngFromDegrees(12, 34),
			}, {
				StartTime: time.Date(2014, 05, 01, 0, 0, 0, 0, time.UTC),
				Location:  s2.LatLngFromDegrees(45.9281883, 7.4156002),
			}},
			bucketDuration: time.Hour * 24,
			want:           []DistanceByTimeBucket{{0, parseDate(t, "2014-04-01")}, {0, parseDate(t, "2014-05-03")}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := Options{Anchors: tc.anchor, BucketDuration: tc.bucketDuration, Reducer: math.Min}
			got, err := TimeBucketDistance(locations, opts)
			if err != nil || !cmp.Equal(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)) {
				t.Errorf("DailyDistance() = %v, %v, want %v", got, err, tc.want)
			}
		})
	}
}
