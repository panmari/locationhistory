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

func TestGenerateRadarItems(t *testing.T) {
	t.Skip() // TODO(panmari): This is not working yet.
	fixedTime := time.Date(2024, 5, 3, 3, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name  string
		items []processor.DistanceByTimeBucket
		want  []opts.RadarData
	}{
		{
			name: "Works all items on same day",
			items: []processor.DistanceByTimeBucket{{
				Distance: 10,
				Bucket:   fixedTime,
			}, {
				Distance: 5,
				Bucket:   fixedTime.Add(2 * time.Minute),
			}, {
				Distance: 20,
				Bucket:   fixedTime.Add(5 * time.Hour),
			}},
			want: []opts.RadarData{
				{Value: math.Log(10)},
				{Value: math.Log(10)},
				{Value: math.Log(5)},
				{Value: math.Log(5)},
				{Value: math.Log(5)},
				{Value: math.Log(20)},
				{Value: math.Log(20)},
				{Value: math.Log(20)},
				{Value: math.Log(20)},
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
