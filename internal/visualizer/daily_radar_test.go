package visualizer

import (
	"math"
	"testing"
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/panmari/locationhistory/internal/processor"
)

func makeFloatSlice(value float64, repeats int) []float64 {
	res := make([]float64, repeats)
	for i := range res {
		res[i] = value
	}
	return res
}

func TestGenerateRadarItems(t *testing.T) {
	fixedTime := time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name  string
		items []processor.DistanceByTimeBucket
		want  []opts.RadarData
	}{
		{
			name: "One entry start of day has zeros after time",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime,
			}},
			want: []opts.RadarData{
				{Value: append([]float64{math.Log(10)}, makeFloatSlice(0, 23)...)},
			},
		},
		{
			name: "One entry end of day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime,
			}},
			want: []opts.RadarData{
				{Value: append([]float64{math.Log(10)}, makeFloatSlice(0, 23)...)},
			},
		},
		{
			name: "Works all items on same day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime.Add(1 * time.Hour),
			}, {
				Distance: 5,
				Bucket:   fixedTime.Add(2 * time.Hour),
			}, {
				Distance: 20,
				Bucket:   fixedTime.Add(6 * time.Hour),
			}},
			want: []opts.RadarData{
				{
					Value: append([]float64{
						math.Log(10),
						math.Log(10),
						math.Log(5),
						// TODO(panmari): Distance 5 should be present some more.
						// math.Log(5),
						// math.Log(5),
						math.Log(20),
						math.Log(20),
						math.Log(20),
						math.Log(20),
					}, makeFloatSlice(0, 17)...),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := generateRadarItems(tc.items)
			if diff := cmp.Diff(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)); diff != "" {
				t.Errorf("generateRadarItems() = %v, want %v. Diff: %v", got, tc.want, diff)
			}
		})
	}

}
