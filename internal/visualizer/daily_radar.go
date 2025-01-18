package visualizer

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/panmari/locationhistory/internal/processor"
)

type Options struct {
	Title    string
	TimeZone *time.Location
}

// generateRadarItems creates daily radar items from the given slice of daily vectors.
func generateRadarItems(dailyVectors []dailyVector) []opts.RadarData {
	res := make([]opts.RadarData, 0, len(dailyVectors))
	for _, dv := range dailyVectors {
		radarValues := make([]float64, len(dv.Values))
		for i, v := range dv.Values {
			// Take log to make curves more interesting.
			radarValues[i] = math.Max(math.Log(v), 0)
		}
		res = append(res, opts.RadarData{Name: dv.Day.Format(time.DateOnly), Value: radarValues})
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

func diffFromMeanHeatmap(dailyVectors []dailyVector) *charts.HeatMap {
	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: opts.Bool(true),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Name: "Week",
			// Show: false,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "category",
			Data: weekDays,
			// Show: false,
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: opts.Bool(true),
			Min:        2,
			Max:        6,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#50a3ba", "#eac736", "#d94e5d"},
			},
		}),
	)

	meanDailyVector := mean(dailyVectors)
	series := make([]opts.HeatMapData, 0, len(dailyVectors)+6)
	// time.ISOWeek has undesired behavior at the beginning/end of the year, mapping to either 52 or 1.
	week := 0
	for i, dv := range dailyVectors {
		y := int(dv.Day.Weekday())
		if i == 0 {
			// For the first entry, prepend empty values for the missing days of the week.
			for j := 1; j < y; j++ {
				series = append(series, opts.HeatMapData{Value: [3]interface{}{week, j, "-"}})
			}
		}
		diff := dv.euclideanDistance(meanDailyVector)
		v := math.Log(diff)
		series = append(series, opts.HeatMapData{Name: dv.Day.Format(time.DateOnly), Value: [3]interface{}{week, y, v}})
		// TODO(panmari): This assumes the data is complete and does not have gaps.
		if dv.Day.Weekday() == time.Saturday {
			week++
		}
	}
	hm.AddSeries("diff from mean day", series)
	return hm
}

func DailyRadar(items []processor.DistanceByTimeBucket, options Options) []components.Charter {
	res := make([]components.Charter, 0, 365)
	dailyVectors := computeDailyVectors(items)
	radarSeries := generateRadarItems(dailyVectors)
	radar := newDailyRadar().SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: options.Title}))
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
	res = append(res, diffFromMeanHeatmap(dailyVectors))
	return res
}
