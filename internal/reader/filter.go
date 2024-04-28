package reader

import "time"

func FilterFunc(locations []Location, filter func(Location) bool) ([]Location, error) {
	res := make([]Location, 0)
	for _, l := range locations {
		if filter(l) {
			res = append(res, l)
		}
	}
	return res, nil
}

func CreateDateFilter(first, last time.Time) func(Location) bool {
	return func(loc Location) bool {
		if t, err := loc.ParsedTimestamp(); err == nil {
			return t.After(first) && t.Before(last)
		}
		// Potentially log error.
		return false
	}
}
