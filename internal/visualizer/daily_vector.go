package visualizer

import (
	"math"
	"time"

	"github.com/panmari/locationhistory/internal/processor"
)

type dailyVector struct {
	day time.Time
	// One value for each hour of the day.
	values [24]float64
}

// DailyVector converts a given list of processor.DistanceByTimeBucket to a list of
// distances with a measurement for each hour, grouped by day.
// * If a time range does not have a value, the last available data point is used. This also applies for dates without coverage.
// * For the last day, distances without values have 0
// Assumes that items are ordered by time ascendingly.
func computeDailyVectors(items []processor.DistanceByTimeBucket) []dailyVector {
	if len(items) == 0 {
		return nil
	}
	res := make([]dailyVector, 0)
	dayCount := 0
	i := 0
	// TODO(panmari): Make use of options.TimeZone for shifting the start of the day before rounding.
	day := items[i].Bucket.Round(time.Hour * 24)
	for {
		distances := [24]float64{}
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
			distances[j] = math.Max(math.Log(items[i].Distance), 0)
			t = t.Add(time.Hour)
		}
		res = append(res, dailyVector{day: day, values: distances})
		dayCount++
		day = day.AddDate(0, 0, 1)
		if i >= len(items)-1 {
			// Last item was processed, finish computation.
			break
		}
	}
	return res
}

func mean(items []dailyVector) [24]float64 {
	res := [24]float64{}
	for _, dv := range items {
		for i := range dv.values {
			res[i] += dv.values[i]
		}
	}
	length := float64(len(items))
	for i := range res {
		res[i] = res[i] / length
	}
	return res
}

func (a dailyVector) euclideanDistance(b dailyVector) float64 {
	if len(a.values) != len(b.values) {
		panic("Vectors must have the same length")
	}

	sum := 0.0
	for i := range a.values {
		diff := a.values[i] - b.values[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}
