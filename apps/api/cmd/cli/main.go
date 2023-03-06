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

	downloader := storage.NewLocalFSDownloader()
	uploader := storage.NewLocalFSUploader(mediaStorePath)
	indexer := index.NewSqliteIndexer(db)
	extractMetadata := exiftool.NewExtractor("/usr/bin/exiftool")
	notifier := notify.NewNoopNotifier()
	resizer := imgresize.NewResizer(thumbnailsPath)

	importMedia := app.ImportDir(app.NewImporter(logger, downloader, extractMetadata, uploader, indexer, resizer, notifier), logger)
	listCollections := index.NewSqliteCollectionLister(db)

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
				Name:  "collection",
				Usage: "list collections",
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
