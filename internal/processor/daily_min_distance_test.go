package processor

import (
	"testing"

	"github.com/golang/geo/s2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/panmari/locationhistory/internal/reader"
)

func TestDailyMinimumDistance(t *testing.T) {
	locations := []reader.Location{
		{
			Timestamp:   "2014-04-01T07:55:51.093Z",
			LatitudeE7:  469287872,
			LongitudeE7: 74171385,
		},
		{
			Timestamp:   "2014-04-01T15:55:51.093Z",
			LatitudeE7:  469281883,
			LongitudeE7: 74156002,
		},
	}
	for _, tc := range []struct {
		name   string
		anchor s2.LatLng
		want   []float64
	}{
		{
			name:   "Anchor at one location gives zero",
			anchor: s2.LatLngFromDegrees(46.9287872, 7.4171385),
			want:   []float64{0},
		},
		{
			name:   "Anchor far away gives non-zero",
			anchor: s2.LatLngFromDegrees(50, 5),
			want:   []float64{385.159008},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DailyMinimumDistance(tc.anchor, locations)
			if err != nil || !cmp.Equal(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)) {
				t.Errorf("DailyMinimumDistance() = %f, %v, want %f", got, err, tc.want)
			}
		})
	}
}

func TestDailyMinimumDistanceMultipleDays(t *testing.T) {
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
		want   []float64
	}{
		{
			name:   "Anchor at one location gives zero",
			anchor: s2.LatLngFromDegrees(46.9287872, 7.4171385),
			want:   []float64{0, 0.0001},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DailyMinimumDistance(tc.anchor, locations)
			if err != nil || !cmp.Equal(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)) {
				t.Errorf("DailyMinimumDistance() = %f, %v, want %f", got, err, tc.want)
			}
		})
	}
}
