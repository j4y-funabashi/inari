package geo_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/geo"
	"gotest.tools/v3/assert"
)

func TestLookupTimezone(t *testing.T) {
	testCases := []struct {
		desc     string
		lat      float64
		lng      float64
		cTime    time.Time
		expected string
	}{
		{
			desc:     "it returns timezone",
			lat:      66.666,
			lng:      66.666,
			cTime:    time.Now(),
			expected: "Europe/London",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockAPI := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						resBody := `{
							  "coords": {
							    "lat": 51.47781,
							    "lon": 0
							  },
							  "tz": "Europe/London"
							}`
						w.Header().Add("Content-Type", "application/json")
						w.Write([]byte(resBody))
					},
				),
			)
			defer mockAPI.Close()

			lookupTimezone := geo.NewTZAPILookupTimezone(mockAPI.URL)

			actual, err := lookupTimezone(tC.lat, tC.lng, tC.cTime)
			assert.NilError(t, err)
			assert.DeepEqual(t, actual, tC.expected)
		})
	}
}
