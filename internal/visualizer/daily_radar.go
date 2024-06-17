package visualizer

import (
	"fmt"
	"math"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/panmari/locationhistory/internal/processor"
)

// generateRadarItems creates daily radar items from the given items.
// * If a time range does not have a value, the last available data point is used. This also applies for data without coverage.
// * For the last day, distances without values have 0
// Assumes that items are ordered by time ascendingly.
// TODO(panmari): Allow shifting to a time zone.
func generateRadarItems(items []processor.DistanceByTimeBucket) []opts.RadarData {
	if len(items) == 0 {
		return nil
	}
	res := make([]opts.RadarData, 0)
	dayCount := 0
	i := 0
	day := items[i].Bucket.Round(time.Hour * 24)
	for {
		distances := make([]float64, 24)
		t := day
		for j := range distances {
			if i+1 < len(items) {
				// Move to next item if its at or after t. In other words, we always take the distance from the closest timestamp that is bigger than t.
				if nextBucket := items[i+1].Bucket; nextBucket.Equal(t) || t.After(nextBucket) {
					i++
				}
			}
			// Take distance from last item by default.
			// To make graph more engaging, apply Log.
			distances[j] = math.Log(items[i].Distance)
			t = t.Add(time.Hour)
		}
		res = append(res, opts.RadarData{Value: distances})
		dayCount++
		day = day.AddDate(0, 0, 1)
		if i >= len(items)-1 {
			// Last item was processed, finish computation.
			break
		}
	}
	return res

}

func indicators() []*opts.Indicator {
	res := make([]*opts.Indicator, 24)
	for i := range res {
		// TODO(panmari): Consider setting Max.
		res[i] = &opts.Indicator{Name: fmt.Sprintf("H%02d", i)}
	}
	return res
}

func color(i, numSeries int) string {
	// Distribute colors evenly in hue space.
	h := 360 / numSeries * i
	return fmt.Sprintf("hsla(%d, 100%%, 50%%, 50%%)", h)
}

func DailyRadar(items []processor.DistanceByTimeBucket) []components.Charter {
	res := make([]components.Charter, 0, 365)
	radarSeries := generateRadarItems(items)
	for i, s := range radarSeries {
		radar := charts.NewRadar()
		radar.SetGlobalOptions(
			charts.WithRadarComponentOpts(opts.RadarComponent{
				Indicator: indicators(),
				Shape:     "circle",
				// SplitNumber: 24,
				SplitLine: &opts.SplitLine{
					Show: true,
					LineStyle: &opts.LineStyle{
						Opacity: 0.1,
					},
				},
			}),
			charts.WithLegendOpts(opts.Legend{
				Show: false,
			}),
		)
		c := color(i, len(radarSeries))
		radar.AddSeries(fmt.Sprintf("Day %d", i), []opts.RadarData{s}, charts.WithItemStyleOpts(opts.ItemStyle{Color: c})).
			SetSeriesOptions(
				charts.WithLineStyleOpts(opts.LineStyle{
					Width:   1,
					Opacity: 0.5,
				}),
			)
		res = append(res, radar)
	}
	return res
}
