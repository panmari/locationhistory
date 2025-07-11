package visualizer

import (
	"math"
	"time"

	"github.com/panmari/locationhistory/internal/processor"
)

type dailyVector struct {
	Day time.Time
	// One value for each hour of the day.
	Values [24]float64
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
			distances[j] = items[i].Distance.Kilometers()
			t = t.Add(time.Hour)
		}
		res = append(res, dailyVector{Day: day, Values: distances})
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
		for i := range dv.Values {
			res[i] += dv.Values[i]
		}
	}
	length := float64(len(items))
	for i := range res {
		res[i] = res[i] / length
	}
	return res
}

func (a dailyVector) euclideanDistance(b [24]float64) float64 {
	sum := 0.0
	for i := range a.Values {
		diff := a.Values[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

// See https://en.wikipedia.org/wiki/Cosine_similarity.
func (a dailyVector) cosineSimilarity(b [24]float64) float64 {
	num := float64(0)
	for i, a := range a.Values {
		num += a * b[i]
	}
	lengthA := float64(0)
	for _, a := range a.Values {
		lengthA += a * a
	}
	lengthA = math.Sqrt(lengthA)

	lengthB := float64(0)
	for _, b := range b {
		lengthB += b * b
	}
	lengthB = math.Sqrt(lengthB)
	return num / (lengthA * lengthB)
}
