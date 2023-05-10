package app_test

import (
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/j4y_funabashi/inari/apps/api/pkg/notify"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	testCases := []struct {
		desc     string
		filePath string
		expected app.Media
	}{
		{
			desc:     "it imports photo with keywords and caption",
			filePath: "p20140321_080118.jpg",
			expected: app.Media{
				ID: "caf73e9785fa54300a051df95cfa2db9",
				MediaMetadata: app.MediaMetadata{
					Hash:        "caf73e9785fa54300a051df95cfa2db9",
					Coordinates: app.Coordinates{},
					Ext:         "jpg",
					MimeType:    "image/jpeg",
					Width:       "2448",
					Height:      "3264",
					CameraMake:  "Samsung",
					CameraModel: "GT-I9100",
					Keywords:    "holiday",
					Title:       "Ferry to Rotterdam",
					Date:        time.Date(2014, time.March, 21, 8, 1, 18, 0, time.UTC),
				},
				FilePath: "2014/20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
				Thumbnails: app.MediaSrc{
					Small:  "sqsm_20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
					Medium: "sqmd_20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
					Large:  "lg_20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
				},
				Collections: []app.Collection{
					{
						ID:    "inbox_2014-03",
						Title: "inbox_2014-03",
						Type:  app.CollectionTypeInbox,
					},
					{
						ID:    "2014-03",
						Title: "2014 March",
						Type:  app.CollectionTypeTimelineMonth,
					},
					{
						ID:    "2014-03-21",
						Title: "Fri, 21 Mar 2014",
						Type:  app.CollectionTypeTimelineDay,
					},
				},
			},
		},
		{
			desc:     "photo with location",
			filePath: "IMG_20220103_134540.jpg",
			expected: app.Media{
				ID: "9b3f4e51bd961cb321ca234a0b4703f9",
				MediaMetadata: app.MediaMetadata{
					Hash: "9b3f4e51bd961cb321ca234a0b4703f9",
					Coordinates: app.Coordinates{
						Lat: 53.8700189722222,
						Lng: -1.561703,
					},
					Ext:         "jpg",
					MimeType:    "image/jpeg",
					Width:       "100",
					Height:      "133",
					Date:        time.Date(2022, time.January, 3, 13, 45, 40, 0, time.UTC),
					CameraMake:  "Fairphone",
					CameraModel: "FP3",
				},
				FilePath: "2022/20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
				Thumbnails: app.MediaSrc{
					Small:  "sqsm_20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
					Medium: "sqmd_20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
					Large:  "lg_20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
				},
				Collections: []app.Collection{
					{
						ID:    "inbox_2022-01",
						Title: "inbox_2022-01",
						Type:  app.CollectionTypeInbox,
					},
					{
						ID:    "2022-01",
						Title: "2022 January",
						Type:  app.CollectionTypeTimelineMonth,
					},
					{
						ID:    "2022-01-03",
						Title: "Mon, 03 Jan 2022",
						Type:  app.CollectionTypeTimelineDay,
					},
				},
				Location: app.Location{
					Country: app.Country{
						Long:  "United Kingdom",
						Short: "GB",
					},
					Region:   "West Yorkshire",
					Locality: "Leeds",
					Coordinates: app.Coordinates{
						Lat: 53.8700189722222,
						Lng: -1.561703,
					},
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			// create uniq test dir to store db + output files
			testID := "inari-test-" + uuid.New().String()
			testDir := filepath.Join(os.TempDir(), testID)
			err := os.MkdirAll(testDir, 0700)
			if err != nil {
				t.Fatalf("failed to create test dir: %s %s", testDir, err)
			}

			// arrange
			importMedia := newImporter(testDir)
			queryMediaDetail := newMediaDetailQuery(testDir)

			// act
			filePath := path.Join("./test_data", tC.filePath)
			iMedia, err := importMedia(filePath)
			assert.NoError(t, err)
			actual, err := queryMediaDetail(iMedia.ID)
			assert.NoError(t, err)

			// assert
			assert.Equal(t, tC.expected, actual)

		})
	}
}

func newImporter(baseDir string) app.Importer {

	// conf
	mediaStorePath := baseDir

	// deps
	db := newDB(baseDir)
	logger := app.NewNullLogger()
	downloader := storage.NewLocalFSDownloader()
	uploader := storage.NewLocalFSUploader(mediaStorePath)
	indexer := index.NewSqliteIndexer(db)
	extractMetadata := exiftool.NewExtractor("/usr/bin/exiftool")
	notifier := notify.NewNoopNotifier()
	createThumbnails := imgresize.NewResizer(filepath.Join(baseDir, "thumbnails"))

	geo := google.NewNullGeocoder()

	return app.NewImporter(logger, downloader, extractMetadata, uploader, indexer, createThumbnails, geo, notifier)
}

func newMediaDetailQuery(testDir string) app.QueryMediaDetail {
	db := newDB(testDir)

	return index.NewQueryMediaDetail(db)
}

func newDB(testDir string) *sql.DB {
	dbFileName := "inari-media-db.db"
	dbFilepath := filepath.Join(testDir, filepath.Base(dbFileName))

	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		panic(err)
	}
	err = index.CreateIndex(db)
	if err != nil {
		panic(err)
	}

	return db
}
