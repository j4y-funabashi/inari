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
			desc:     "photo with keywords and caption",
			filePath: "p20140321_080118.jpg",
			expected: app.Media{
				MediaMetadata: app.MediaMetadata{
					Hash:        "caf73e9785fa54300a051df95cfa2db9",
					Location:    app.Location{},
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

			// act
			filePath := path.Join("./test_data", tC.filePath)
			actual, err := importMedia(filePath)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, tC.expected, actual)

		})
	}
}

func newImporter(testDir string) app.Importer {

	// conf
	dbFileName := "inari-media-db.db"
	dbFilepath := filepath.Join(testDir, filepath.Base(dbFileName))
	mediaStorePath := testDir

	// deps
	logger := app.NewNullLogger()
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		logger.Error("failed to open db",
			"err", err)
		panic(err)
	}
	err = index.CreateIndex(db)
	if err != nil {
		logger.Error("failed to create index",
			"err", err)
		panic(err)
	}
	downloader := storage.NewLocalFSDownloader()
	uploader := storage.NewLocalFSUploader(mediaStorePath)
	indexer := index.NewSqliteIndexer(db)
	extractMetadata := exiftool.NewExtractor("/usr/bin/exiftool")
	notifier := notify.NewNoopNotifier()

	return app.NewImporter(logger, downloader, extractMetadata, uploader, indexer, notifier)
}
