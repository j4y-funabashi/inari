package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/inconshreveable/log15"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	appconfig "github.com/j4y_funabashi/inari/apps/api/pkg/app_config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func main() {

	baseDir := filepath.Join(os.TempDir(), "inari")
	logger := log.New()

	////////////////////

	importMedia := app.ImportDir(appconfig.NewMediaImporter(os.TempDir()), logger)
	listCollections := appconfig.NewListCollections(baseDir)
	importGPX := app.ImportDir(appconfig.NewImportGPX(baseDir), logger)

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
