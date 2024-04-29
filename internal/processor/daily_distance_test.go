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
		want   []DistanceByBucket
	}{
		{
			name: "Anchor at one location gives zero",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(46.9287872, 7.4171385),
			}},
			want: []DistanceByBucket{{0, "2014-04-01"}},
		},
		{
			name: "Anchor far away gives non-zero",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(50, 5),
			}},
			want: []DistanceByBucket{{385.159008, "2014-04-01"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DailyDistance(tc.anchor, locations, math.Min)
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
			Timestamp:   "2014-05-02T15:55:51.093Z",
			LatitudeE7:  459281883,
			LongitudeE7: 74156002,
		},
	}
	for _, tc := range []struct {
		name   string
		anchor []Anchor
		want   []DistanceByBucket
	}{
		{
			name: "Anchor at one location gives zero",
			anchor: []Anchor{{
				StartTime: time.Time{},
				Location:  s2.LatLngFromDegrees(46.9287872, 7.4171385),
			}},
			want: []DistanceByBucket{{0, "2014-04-01"}, {111.2615837, "2014-05-03"}},
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
			want: []DistanceByBucket{{0, "2014-04-01"}, {0, "2014-05-03"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DailyDistance(tc.anchor, locations, math.Min)
			if err != nil || !cmp.Equal(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)) {
				t.Errorf("DailyDistance() = %v, %v, want %v", got, err, tc.want)
			}
		})
	}
}
