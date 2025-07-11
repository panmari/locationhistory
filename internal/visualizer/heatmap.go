package visualizer

import (
	"math"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/panmari/locationhistory/internal/processor"
)

var (
	weekDays = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
)

func transformToHeatMapData(items []processor.DistanceByTimeBucket) []opts.HeatMapData {
	if len(items) == 0 {
		return nil
	}
	res := make([]opts.HeatMapData, 0, len(items)+6)
	// time.ISOWeek has undesired behavior at the beginning/end of the year, mapping to either 52 or 1.
	week := 0
	for i, dbb := range items {
		t := dbb.Bucket
		y := int(t.Weekday())
		if i == 0 {
			// For the first entry, prepend empty values for the missing days of the week.
			for j := 1; j < y; j++ {
				res = append(res, opts.HeatMapData{Value: [3]interface{}{week, j, "-"}})
			}
		}
		v := math.Log(dbb.Distance.Kilometers())
		res = append(res, opts.HeatMapData{Name: t.Format(time.DateOnly), Value: [3]interface{}{week, y, v}})
		// TODO(panmari): This assumes the data is complete and does not have gaps.
		if t.Weekday() == time.Saturday {
			week++
		}
	}
	return res
}

func Heatmap(items []processor.DistanceByTimeBucket) *charts.HeatMap {
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
			Min:        0,
			Max:        6,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#50a3ba", "#eac736", "#d94e5d"},
			},
		}),
	)

	hm.AddSeries("distance from anchor", transformToHeatMapData(items))
	return hm

}
