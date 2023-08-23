package index_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/matryer/is"
	"github.com/tkrajina/gpxgo/gpx"
)

func TestIndex(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "it saves media to db",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			is := is.New(t)

			// arrange
			dbFilepath := filepath.Join(os.TempDir(), fmt.Sprintf("inari-test-db-%s", uuid.New().String()))
			db, err := sql.Open("sqlite3", dbFilepath)
			if err != nil {
				t.Fatalf("failed to open sqlite db: %s", err)
			}

			indexMedia := index.NewSqliteIndexer(db)
			getMedia := index.NewQueryMediaDetail(db)

			// act
			media := app.Media{
				MediaMetadata: app.MediaMetadata{
					Hash: "test-hash",
					Date: time.Now(),
				},
			}

			err = index.CreateIndex(db)
			if err != nil {
				t.Fatalf("%s", err)
			}

			iMedia, err := indexMedia(media)
			if err != nil {
				t.Fatalf("%s", err)
			}

			actualMedia, err := getMedia(iMedia.ID)
			if err != nil {
				t.Fatalf("%s", err)
			}

			// assert
			is.Equal(iMedia.ID, actualMedia.ID)
		})
	}
}

func TestFindNearestGPX(t *testing.T) {
	testCases := []struct {
		desc          string
		points        []gpx.GPXPoint
		currentTime   time.Time
		expectedGPX   app.GPXPoint
		hoursBoundary int
	}{
		{
			desc:          "nearest is in the future",
			hoursBoundary: 5,
			points: []gpx.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 13, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 14, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 10, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{
				Timestamp: time.Date(2022, time.January, 28, 13, 0, 0, 0, time.UTC),
				Coordinates: app.Coordinates{
					Lat: 10,
					Lng: 20,
				},
			},
		},
		{
			desc:          "nearest is in the past",
			hoursBoundary: 5,
			points: []gpx.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 15, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 14, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 11, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 10, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{
				Timestamp: time.Date(2022, time.January, 28, 11, 0, 0, 0, time.UTC),
				Coordinates: app.Coordinates{
					Lat: 10,
					Lng: 20,
				},
			},
		},
		{
			desc:          "does not fetch future points over the hours boundary",
			hoursBoundary: 2,
			points: []gpx.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 15, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{},
		},
		{
			desc:          "does not fetch past points over the hours boundary",
			hoursBoundary: 2,
			points: []gpx.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 9, 0, 0, 0, time.UTC),
					Point: gpx.Point{
						Latitude:  10,
						Longitude: 20,
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			is := is.New(t)

			// arrange
			dbFilepath := filepath.Join(os.TempDir(), fmt.Sprintf("inari-test-db-%s", uuid.New().String()))
			db, err := sql.Open("sqlite3", dbFilepath)
			if err != nil {
				t.Fatalf("failed to open sqlite db: %s", err)
			}

			err = index.CreateIndex(db)
			if err != nil {
				t.Fatalf("%s", err)
			}
			fetchNearestPoint := index.NewQueryNearestGPX(db, tC.hoursBoundary)

			// act
			numPoints, err := index.InsertGPXPoints(db, tC.points)
			t.Logf("points inserted: %d", numPoints)
			is.NoErr(err)
			nearestPoint, err := fetchNearestPoint(tC.currentTime)
			is.NoErr(err)

			// assert
			is.Equal(nearestPoint, tC.expectedGPX)
		})
	}
}
