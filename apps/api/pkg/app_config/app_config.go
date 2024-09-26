package appconfig

import (
	"database/sql"
	"os"
	"path/filepath"

	log "github.com/inconshreveable/log15"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/geo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/j4y_funabashi/inari/apps/api/pkg/notify"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
)

func NewMediaImporter(baseDirectory string, c ...func(*app.MediaImporterConfig)) app.Importer {
	baseDir := filepath.Join(baseDirectory)
	mediaStorePath := filepath.Join(baseDir, "media")
	thumbnailsPath := filepath.Join(baseDir, "thumbnails")

	err := os.MkdirAll(baseDir, 0700)
	if err != nil {
		panic("failed to create root dir: " + err.Error())
	}

	// deps
	db := newDB(baseDir)
	mediaDetail := index.NewQueryMediaDetail(db)
	logger := log.New()
	downloader := storage.NewLocalFSDownloader()
	uploader := storage.NewLocalFSUploader(mediaStorePath)
	indexer := index.NewSqliteIndexer(db)
	extractMetadata := exiftool.NewExtractor("/usr/bin/exiftool")
	notifier := notify.NewNoopNotifier()
	createThumbnails := imgresize.NewResizer(thumbnailsPath)

	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	geo2tzBaseURL := "http://localhost:2004"
	lookupTimezone := geo.NewTZAPILookupTimezone(geo2tzBaseURL)
	googleGeocodeURL := "https://maps.googleapis.com/maps/api/geocode/json"
	queryNearestGPX := index.NewQueryNearestGPX(db, 8)
	mediaGeocoder := google.NewMediaGeocoder(queryNearestGPX, lookupTimezone, logger, googleAPIKey, googleGeocodeURL)

	config := app.MediaImporterConfig{
		FetchMediaDetail:   mediaDetail,
		Logger:             logger,
		DownloadFromBackup: downloader,
		ExtractMetadata:    extractMetadata,
		UploadToMediaStore: uploader,
		IndexMedia:         indexer,
		CreateThumbnails:   createThumbnails,
		Geocode:            mediaGeocoder,
		NotifyDownstream:   notifier,
	}

	for _, nc := range c {
		nc(&config)
	}

	return app.NewImporter(config)
}

func WithNullLogger() func(*app.MediaImporterConfig) {
	return func(c *app.MediaImporterConfig) {
		c.Logger = app.NewNullLogger()
	}
}

func WithNullGeocoder() func(*app.MediaImporterConfig) {
	return func(c *app.MediaImporterConfig) {
		c.Geocode = google.NewNullGeocoder()
	}
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
