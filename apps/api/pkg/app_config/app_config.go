package appconfig

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	log "github.com/inconshreveable/log15"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/geo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"github.com/j4y_funabashi/inari/apps/api/pkg/gpx"
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
	logger := log.New()
	db := newDB(baseDir)
	mediaDetail := index.NewQueryMediaDetail(db)
	downloader := storage.NewLocalFSDownloader()
	uploader := storage.NewLocalFSUploader(mediaStorePath)
	indexer := index.NewSqliteIndexer(db)
	extractMetadata := exiftool.NewExtractor()
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

func NewListCollections(baseDir string) app.CollectionLister {
	db := newDB(baseDir)
	return index.NewSqliteCollectionLister(db)
}

func NewImportGPX(baseDir string) app.Importer {
	logger := log.New()
	db := newDB(baseDir)
	geo2tzBaseURL := "http://localhost:2004"
	lookupTimezone := geo.NewTZAPILookupTimezone(geo2tzBaseURL)

	return gpx.NewGpxImporter(
		gpx.NewAddLocationToGPXPoints(lookupTimezone),
		index.NewSaveGPXPoints(db),
		logger,
	)
}

func NewMediaDetail(baseDir string) app.QueryMediaDetail {
	db := newDB(baseDir)
	return index.NewQueryMediaDetail(db)
}

func NewCollectionDetail(baseDir string) app.CollectionDetailQuery {
	db := newDB(baseDir)
	return index.NewSqliteCollectionDetail(db)
}

func NewDeleteMedia(baseDir string) app.DeleteMedia {
	db := newDB(baseDir)
	return index.NewDeleteMedia(db)
}

func NewUpdateMediaCaption(baseDir string) app.UpdateMediaTextProperty {
	db := newDB(baseDir)
	return index.NewUpdateMediaCaption(db)
}

func NewUpdateMediaHashtag(baseDir string) app.UpdateMediaTextProperty {
	db := newDB(baseDir)
	return index.NewUpdateMediaTag(db)
}

func NewExporter(logger app.Logger, queryMediaDetail app.QueryMediaDetail, mediaUploader, postUploader app.UploaderB, baseDir string) app.Exporter {
	return func(mediaID string) error {
		// fetch media
		media, err := queryMediaDetail(mediaID)
		if err != nil {
			return err
		}

		// convert media to microformat
		mf := media.ToMicroformat()
		logger.Info("mf!", "mf", mf)

		// save microformat to s3
		mfJson, err := json.Marshal(mf)
		if err != nil {
			return err
		}
		postFilePath := filepath.Join("posts", media.PostFilename())
		err = postUploader(mfJson, postFilePath, "application/json")
		if err != nil {
			return err
		}

		// save image to media bucket
		thumbnailPath := filepath.Join(baseDir, "thumbnails", media.Thumbnails.Large)
		thumbnailData, err := os.ReadFile(thumbnailPath)
		if err != nil {
			return err
		}
		err = mediaUploader(thumbnailData, media.Thumbnails.Large, "image/jpeg")
		if err != nil {
			return err
		}

		logger.Info("exported media", "media", media.Thumbnails.Large, "post", media.PostFilename())

		// mark media as exported
		// call github workflow URL
		return nil
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
