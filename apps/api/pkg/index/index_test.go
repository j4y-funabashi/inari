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
	"gotest.tools/v3/assert"
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
			assert.Equal(t, iMedia.ID, actualMedia.ID)
		})
	}
}

func TestFindNearestGPX(t *testing.T) {
	testCases := []struct {
		desc          string
		points        []app.GPXPoint
		currentTime   time.Time
		expectedGPX   app.GPXPoint
		hoursBoundary int
	}{
		{
			desc:          "nearest is in the future",
			hoursBoundary: 5,
			points: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 13, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 54.26073423586785793304443359375,
							Lng: -4.46283244527876377105712890625,
						},
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 14, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 10, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{
				Timestamp: time.Date(2022, time.January, 28, 13, 0, 0, 0, time.UTC),
				Location: app.Location{
					Coordinates: app.Coordinates{
						Lat: 54.26073423586786,
						Lng: -4.462832445278764,
					},
				},
			},
		},
		{
			desc:          "nearest is in the past",
			hoursBoundary: 5,
			points: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 15, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 14, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 11, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
				{
					Timestamp: time.Date(2022, time.January, 28, 10, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{
				Timestamp: time.Date(2022, time.January, 28, 11, 0, 0, 0, time.UTC),
				Location: app.Location{
					Coordinates: app.Coordinates{
						Lat: 10,
						Lng: 20,
					},
				},
			},
		},
		{
			desc:          "does not fetch future points over the hours boundary",
			hoursBoundary: 2,
			points: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 15, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{},
		},
		{
			desc:          "does not fetch past points over the hours boundary",
			hoursBoundary: 2,
			points: []app.GPXPoint{
				{
					Timestamp: time.Date(2022, time.January, 28, 9, 0, 0, 0, time.UTC),
					Location: app.Location{
						Coordinates: app.Coordinates{
							Lat: 10,
							Lng: 20,
						},
					},
				},
			},
			currentTime: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedGPX: app.GPXPoint{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

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
			saveGPXPoints := index.NewSaveGPXPoints(db)

			// act
			err = saveGPXPoints(tC.points)
			assert.NilError(t, err)
			nearestPoint, err := fetchNearestPoint(tC.currentTime)
			assert.NilError(t, err)

			// assert
			assert.Equal(t, nearestPoint, tC.expectedGPX)
		})
	}
}

func TestUpdateMediaTags(t *testing.T) {
	testCases := []struct {
		desc          string
		media         app.Media
		expectedMedia app.Media
		newTags       []string
	}{
		{
			desc: "it adds multiple slugified tags",
			media: app.Media{
				ID: "test-id-1",
				MediaMetadata: app.MediaMetadata{
					Date: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
				},
			},
			newTags: []string{"tEst tag 1", "tesTtag*2"},
			expectedMedia: app.Media{
				ID: "test-id-1",
				MediaMetadata: app.MediaMetadata{
					Date: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
				},
				Collections: []app.Collection{
					{ID: "hashtag__test-tag-1", Title: "tEst tag 1", Type: "hashtag"},
					{ID: "hashtag__testtag-2", Title: "tesTtag*2", Type: "hashtag"},
				},
			},
		},
		{
			desc: "it de dupes tags",
			media: app.Media{
				ID: "test-id-1",
				MediaMetadata: app.MediaMetadata{
					Date: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
				},
			},
			newTags: []string{"tEst tag 1", "tEst tag 1"},
			expectedMedia: app.Media{
				ID: "test-id-1",
				MediaMetadata: app.MediaMetadata{
					Date: time.Date(2022, time.January, 28, 12, 0, 0, 0, time.UTC),
				},
				Collections: []app.Collection{
					{ID: "hashtag__test-tag-1", Title: "tEst tag 1", Type: "hashtag"},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

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

			updateTags := index.NewUpdateMediaTag(db)
			getMedia := index.NewQueryMediaDetail(db)

			// act
			_, err = index.InsertMedia(db, tC.media)
			assert.NilError(t, err)

			for _, tag := range tC.newTags {
				err := updateTags(tC.media.ID, tag)
				assert.NilError(t, err)
			}

			newMedia, err := getMedia(tC.media.ID)
			assert.NilError(t, err)
			assert.DeepEqual(t, tC.expectedMedia, newMedia)

		})
	}
}
