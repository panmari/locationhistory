package visualizer

import (
	"log"
	"math"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/panmari/locationhistory/internal/processor"
)

func generateBarItems(items []processor.DistanceByBucket) []opts.BarData {
	res := make([]opts.BarData, 0, len(items))
	for _, i := range items {
		res = append(res, opts.BarData{Value: math.Log(i.Distance * 1000)})
	}
	return res
}

func generateXAxis(items []processor.DistanceByBucket) []string {
	res := make([]string, 0, len(items))
	for _, i := range items {
		res = append(res, i.Bucket)
	}
	return res
}

func BarChart(items []processor.DistanceByBucket) *charts.Bar {
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: false}),
		charts.WithXAxisOpts(opts.XAxis{Show: false}),
		charts.WithYAxisOpts(opts.YAxis{Show: false}),
	)

	// Put data into instance
	y := generateBarItems(items)
	x := generateXAxis(items)
	log.Default().Print(y)
	bar.SetXAxis(x).AddSeries("Distances", y)
	return bar
}
