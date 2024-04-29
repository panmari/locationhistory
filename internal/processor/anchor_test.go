package processor

import (
	"testing"
	"time"

	"github.com/golang/geo/s2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseAnchors(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  []Anchor
	}{
		{
			name:  "One anchor without date",
			input: "10,15",
			want: []Anchor{
				{
					StartTime: time.Time{},
					Location:  s2.LatLngFromDegrees(10, 15),
				},
			},
		}, {
			name:  "Two anchors",
			input: "2007-01-31,10.0,15.0:2007-02-12,11.0,16.0",
			want: []Anchor{
				{
					StartTime: time.Date(2007, 1, 31, 0, 0, 0, 0, time.UTC),
					Location:  s2.LatLngFromDegrees(10, 15),
				}, {
					StartTime: time.Date(2007, 2, 12, 0, 0, 0, 0, time.UTC),
					Location:  s2.LatLngFromDegrees(11, 16),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseAnchors(tc.input)
			if err != nil || !cmp.Equal(got, tc.want, cmpopts.EquateApprox(0.001, 0.001)) {
				t.Errorf("ParseAnchors() = %v, %v, want %v", got, err, tc.want)
			}
		})
	}

}
