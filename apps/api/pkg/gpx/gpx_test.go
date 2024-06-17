package gpx_test

import (
	"testing"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"github.com/j4y_funabashi/inari/apps/api/pkg/gpx"
	"gotest.tools/v3/assert"
)

func TestImporter(t *testing.T) {

	testCases := []struct {
		desc      string
		gpxPoints []app.GPXPoint
		expected  []app.GPXPoint
	}{
		{
			desc: "first and last points are in the same timezone",
			gpxPoints: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.SpainCoordinates,
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 19, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.SpainCoordinates,
					},
				},
			},
			expected: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 13, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.SpainCoordinates,
						Timezone:    "Europe/Madrid",
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 20, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.SpainCoordinates,
						Timezone:    "Europe/Madrid",
					},
				},
			},
		},
		{
			desc: "first and last points are in different timezones",
			gpxPoints: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.LusakaCoordinates,
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 19, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.SpainCoordinates,
					},
				},
			},
			expected: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 14, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.LusakaCoordinates,
						Timezone:    "Africa/Lusaka",
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 20, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.SpainCoordinates,
						Timezone:    "Europe/Madrid",
					},
				},
			},
		},
		{
			desc: "it applies daylight savings to points",
			gpxPoints: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.June, 10, 14, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.UKCoordinates,
					},
				},
			},
			expected: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.June, 10, 15, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: google.UKCoordinates,
						Timezone: "Europe/London",
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			addLocationToGPXPoints := gpx.NewAddLocationToGPXPoints(google.NewNullLookupTimezone())

			actual, err := addLocationToGPXPoints(tC.gpxPoints)
			assert.NilError(t, err)
			assert.DeepEqual(t, actual, tC.expected)
		})
	}
}
