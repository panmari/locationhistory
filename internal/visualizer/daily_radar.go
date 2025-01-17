package visualizer

import (
	"fmt"
	"slices"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/panmari/locationhistory/internal/processor"
)

type Options struct {
	TimeZone *time.Location
}

// generateRadarItems creates daily radar items from the given daily vectors.
func generateRadarItems(dailyVectors []dailyVector) []opts.RadarData {
	res := make([]opts.RadarData, len(dailyVectors))
	for _, dv := range dailyVectors {
		res = append(res, opts.RadarData{Name: dv.day.Format(time.DateOnly), Value: dv.values[:]})
	}
	return res

}

func indicators() []*opts.Indicator {
	res := make([]*opts.Indicator, 24)
	for i := range res {
		// TODO(panmari): Set max according to data.
		res[i] = &opts.Indicator{Name: fmt.Sprintf("%02d:00", i), Max: 6}
	}
	// Reverse indicators to make hours appear clockwise in radar.
	slices.Reverse(res)
	return res
}

func color(i, numSeries int) string {
	// Distribute colors evenly in hue space.
	h := 360 / numSeries * i
	return fmt.Sprintf("hsla(%d, 100%%, 50%%, 50%%)", h)
}

// newDailyRadar creates an empty, initialized radar plot for visualizing data on 24h.
func newDailyRadar() *charts.Radar {
	return charts.NewRadar().SetGlobalOptions(
		charts.WithRadarComponentOpts(opts.RadarComponent{
			Indicator: indicators(),
			Shape:     "circle",
			// SplitNumber: 24,
			SplitLine: &opts.SplitLine{
				Show: opts.Bool(true),
				LineStyle: &opts.LineStyle{
					Opacity: 0.1,
				},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Formatter: "{a}", // Prints as name of the series, see docs for `Formatter`.
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)
}

func DailyRadar(items []processor.DistanceByTimeBucket, options Options) []components.Charter {
	res := make([]components.Charter, 0, 365)
	dailyVectors := computeDailyVectors(items)
	radarSeries := generateRadarItems(dailyVectors)
	radar := newDailyRadar()
	for i, s := range radarSeries {
		// In order to make radar appear clockwise, reverse distances here.
		slices.Reverse(s.Value.([]float64))
		c := color(i, len(radarSeries))
		radar.AddSeries(s.Name, []opts.RadarData{s},
			charts.WithItemStyleOpts(opts.ItemStyle{Color: c}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Width:   1,
				Opacity: 0.5,
			}))
	}
	res = append(res, radar)
	return res
}
