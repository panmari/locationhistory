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
			name: "One entry start of day has same value whole day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime,
			}},
			want: []opts.RadarData{
				{Name: "2024-05-03", Value: makeFloatSlice(math.Log(10), 24)},
			},
		}, {
			name: "One entry end of day has same value whole day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime,
			}},
			want: []opts.RadarData{
				{Name: "2024-05-03", Value: makeFloatSlice(math.Log(10), 24)},
			},
		}, {
			name: "Processes one day with multiple values",
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
					Name: "2024-05-03",
					Value: append([]float64{
						math.Log(10),
						math.Log(10),
						math.Log(5),
						math.Log(5),
						math.Log(5),
						math.Log(5),
					}, makeFloatSlice(math.Log(20), 18)...),
				},
			},
		}, {
			name: "Processes multi day data",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime.Add(1 * time.Hour),
			}, {
				Distance: 20,
				Bucket:   fixedTime.Add(26 * time.Hour),
			}},
			want: []opts.RadarData{
				{
					Name:  "2024-05-03",
					Value: makeFloatSlice(math.Log(10), 24),
				}, {
					Name: "2024-05-04",
					Value: append(
						makeFloatSlice(math.Log(10), 2), // From the day before
						makeFloatSlice(math.Log(20), 22)...),
				},
			},
		}, {
			name: "Processes multi day data with gap",
			items: []processor.DistanceByTimeBucket{
				{
					Distance: 10,
					Bucket:   fixedTime.Add(1 * time.Hour),
				},
				// No distances for more than 24h
				{
					Distance: 20,
					Bucket:   fixedTime.Add(50 * time.Hour),
				}},
			want: []opts.RadarData{
				{
					Name:  "2024-05-03",
					Value: makeFloatSlice(math.Log(10), 24),
				}, {
					Name:  "2024-05-04",
					Value: makeFloatSlice(math.Log(10), 24), // From the date not covered.
				}, {
					Name: "2024-05-05",
					Value: append(
						makeFloatSlice(math.Log(10), 2), // From the day before
						makeFloatSlice(math.Log(20), 22)...),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := generateRadarItems(tc.items, Options{})
			if diff := cmp.Diff(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)); diff != "" {
				t.Errorf("generateRadarItems() = %v, want %v. Diff: %v", got, tc.want, diff)
			}
		})
	}

}
