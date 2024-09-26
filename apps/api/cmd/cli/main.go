package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/inconshreveable/log15"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	appconfig "github.com/j4y_funabashi/inari/apps/api/pkg/app_config"
	"github.com/j4y_funabashi/inari/apps/api/pkg/geo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/gpx"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func main() {

	// conf
	baseDir := filepath.Join(os.TempDir(), "inari")
	dbFilepath := filepath.Join(baseDir, "inari-media-db.db")
	geo2tzBaseURL := "http://localhost:2004"
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

	////////////////////

	lookupTimezone := geo.NewTZAPILookupTimezone(geo2tzBaseURL)

	importMedia := app.ImportDir(appconfig.NewMediaImporter(os.TempDir()), logger)

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
