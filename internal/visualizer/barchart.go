package visualizer

import (
	"fmt"
	"math"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/panmari/locationhistory/internal/processor"
)

func generateBarItems(items []processor.DistanceByTimeBucket) []opts.BarData {
	res := make([]opts.BarData, 0, len(items))
	for _, i := range items {
		v := math.Max(math.Log(i.Distance.Kilometers()*1000)-3.3, 0)
		res = append(res, opts.BarData{Value: v})
	}
	return res
}

func generateXAxis(items []processor.DistanceByTimeBucket) []string {
	res := make([]string, 0, len(items))
	for _, i := range items {
		res = append(res, i.Bucket.Format(time.DateOnly))
	}
	return res
}

func BarChart(items []processor.DistanceByTimeBucket) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
		charts.WithXAxisOpts(opts.XAxis{Show: opts.Bool(false),
			AxisTick:  &opts.AxisTick{Show: opts.Bool(false)},
			AxisLabel: &opts.AxisLabel{Show: opts.Bool(false)}}),
		charts.WithYAxisOpts(opts.YAxis{Show: opts.Bool(false),
			AxisLabel:   &opts.AxisLabel{Show: opts.Bool(false)},
			AxisPointer: &opts.AxisPointer{Show: opts.Bool(false)}}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", 365*5+20),
			Height: "600px",
		}),
	)

	// Put data into instance
	y := generateBarItems(items)
	x := generateXAxis(items)
	bar.SetXAxis(x).AddSeries("Distances", y)
	return bar
}
