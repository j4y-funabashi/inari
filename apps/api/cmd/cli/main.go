package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/inconshreveable/log15"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"github.com/j4y_funabashi/inari/apps/api/pkg/gpx"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/j4y_funabashi/inari/apps/api/pkg/notify"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func main() {

	// conf
	baseDir := filepath.Join(os.TempDir(), "inari")
	mediaStorePath := filepath.Join(baseDir, "media")
	thumbnailsPath := filepath.Join(baseDir, "thumbnails")
	dbFilepath := filepath.Join(baseDir, "inari-media-db.db")
	apiKey := os.Getenv("GOOGLE_API_KEY")
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"

	err := os.MkdirAll(baseDir, 0700)
	if err != nil {
		panic("failed to create root dir: " + err.Error())
	}

	// deps
	logger := log.New()

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

	mediaDetail := index.NewQueryMediaDetail(db)
	downloader := storage.NewLocalFSDownloader()
	uploader := storage.NewLocalFSUploader(mediaStorePath)
	indexer := index.NewSqliteIndexer(db)
	extractMetadata := exiftool.NewExtractor("/usr/bin/exiftool")
	notifier := notify.NewNoopNotifier()
	resizer := imgresize.NewResizer(thumbnailsPath)
	queryNearestGPX := index.NewQueryNearestGPX(db, 8)
	lookupTimezone := google.NewLookupTimezone(apiKey)
	mediaGeocoder := google.NewMediaGeocoder(queryNearestGPX, lookupTimezone, logger, apiKey, baseURL)

	importMedia := app.ImportDir(app.NewImporter(mediaDetail, logger, downloader, extractMetadata, uploader, indexer, resizer, mediaGeocoder, notifier), logger)
	listCollections := index.NewSqliteCollectionLister(db)
	importGPX := app.ImportDir(gpx.NewGpxImporter(
		gpx.NewAddLocationToGPXPoints(lookupTimezone),
		index.NewSaveGPXPoints(db),
		logger,
	), logger)

	// app commands
	app := &cli.App{
		Name:  "inari",
		Usage: "photo organiser",
		Commands: []*cli.Command{
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "import media",
				Action: func(cCtx *cli.Context) error {
					inputFilename := cCtx.Args().First()
					err := importMedia(inputFilename)
					return err
				},
			},
			{
				Name:  "igpx",
				Usage: "import gpx data",
				Action: func(cCtx *cli.Context) error {
					inputFilename := cCtx.Args().First()
					err := importGPX(inputFilename)
					return err
				},
			},
			{
				Name:    "collection",
				Aliases: []string{"lsc"},
				Usage:   "list collections",
				Action: func(cCtx *cli.Context) error {
					collectionType := cCtx.Args().First()
					cols, err := listCollections(collectionType)
					out, _ := json.Marshal(cols)
					fmt.Printf("%s", string(out))
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error("failed to run cli app",
			"err", err)
	}

}
