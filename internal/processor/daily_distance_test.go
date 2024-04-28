package processor

import (
	"math"
	"testing"

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
		anchor s2.LatLng
		want   []DistanceByBucket
	}{
		{
			name:   "Anchor at one location gives zero",
			anchor: s2.LatLngFromDegrees(46.9287872, 7.4171385),
			want:   []DistanceByBucket{{0, "2014-04-01"}},
		},
		{
			name:   "Anchor far away gives non-zero",
			anchor: s2.LatLngFromDegrees(50, 5),
			want:   []DistanceByBucket{{385.159008, "2014-04-01"}},
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
			Timestamp:   "2014-05-01T15:55:51.093Z",
			LatitudeE7:  459281883,
			LongitudeE7: 74156002,
		},
	}
	for _, tc := range []struct {
		name   string
		anchor s2.LatLng
		want   []DistanceByBucket
	}{
		{
			name:   "Anchor at one location gives zero",
			anchor: s2.LatLngFromDegrees(46.9287872, 7.4171385),
			want:   []DistanceByBucket{{0, "2014-04-01"}, {111.2615837, "2014-05-02"}},
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
