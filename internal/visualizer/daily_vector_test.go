package visualizer

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/go-units/unit"
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
				Distance: 10 * unit.Kilometer,
				Bucket:   fixedTime,
			}},
			want: []dailyVector{
				{Day: fixedTime, Values: makeFloatArray(10, 24)},
			},
		}, {
			name: "One entry end of day has same value whole day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10 * unit.Kilometer,
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

func TestCosineSimilarity(t *testing.T) {
	for _, tc := range []struct {
		name string
		a, b dailyVector
		want float64
	}{
		{
			name: "Empty vectors",
			a:    dailyVector{Values: [24]float64{}},
			b:    dailyVector{Values: [24]float64{}},
			want: 0,
		},
		{
			name: "Same direction, same length",
			a:    dailyVector{Values: makeFloatArray(10, 24)},
			b:    dailyVector{Values: makeFloatArray(10, 24)},
			want: 1,
		},
		{
			name: "Opposite direction",
			a:    dailyVector{Values: makeFloatArray(10, 24)},
			b:    dailyVector{Values: makeFloatArray(-10, 24)},
			want: -1,
		}, {
			name: "Same direction, different length",
			a:    dailyVector{Values: makeFloatArray(10, 24)},
			b:    dailyVector{Values: makeFloatArray(2, 24)},
			want: 1,
		}, {
			name: "Different direction",
			a:    dailyVector{Values: makeFloatArray(1, 24)},
			b:    dailyVector{Values: [24]float64{1}},
			want: 0.204,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.cosineSimilarity(tc.b.Values)
			if diff := cmp.Diff(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)); diff != "" {
				t.Errorf("cosineSimilarity() = %v, want %v.", got, tc.want)
			}
		})
	}
}
