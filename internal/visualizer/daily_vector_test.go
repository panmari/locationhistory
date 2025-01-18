package visualizer

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/panmari/locationhistory/internal/processor"
)

func makeFloatArray(value float64, repeats int) [24]float64 {
	res := [24]float64{}
	for i := range res {
		res[i] = value
	}
	return res
}

func TestComputeDailyVectors(t *testing.T) {
	fixedTime := time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name  string
		items []processor.DistanceByTimeBucket
		want  []dailyVector
	}{
		{
			name: "One entry start of day has same value whole day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime,
			}},
			want: []dailyVector{
				{Day: fixedTime, Values: makeFloatArray(10, 24)},
			},
		}, {
			name: "One entry end of day has same value whole day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime,
			}},
			want: []dailyVector{
				{Day: fixedTime, Values: makeFloatArray(10, 24)},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := computeDailyVectors(tc.items)
			if diff := cmp.Diff(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)); diff != "" {
				t.Errorf("generateRadarItems() = %v, want %v. Diff: %v", got, tc.want, diff)
			}
		})
	}

}
