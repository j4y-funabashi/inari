package geo_test

import (
	"testing"

	"github.com/j4y_funabashi/inari/apps/api/pkg/geo"
)

func TestLookupTimezone(t *testing.T) {
	testCases := []struct {
		desc string
		lat float64
		lng float64
		cTime time.Time
	}{
		{
			desc: "it returns timezone",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			apiBaseURL := ""
			lookupTimezone := geo.NewTZAPILookupTimezone(apiBaseURL)

			actual := lookupTimezone()
		})
	}
}
