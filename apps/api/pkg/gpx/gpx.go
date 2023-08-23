package gpx

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/tkrajina/gpxgo/gpx"
)

func NewGpxImporter(logger app.Logger, db *sql.DB) app.Importer {
	return func(inputFilename string) (app.Media, error) {
		m := app.Media{}
		startTime := time.Now()

		ext := strings.ToLower(filepath.Ext(inputFilename))
		if ext != ".gpx" {
			return m, nil
		}

		gpxBytes, err := os.ReadFile(inputFilename)
		if err != nil {
			logger.Error("failed to open file",
				"err", err,
				"filename", inputFilename)
			return m, err
		}
		gpxFile, err := gpx.ParseBytes(gpxBytes)
		if err != nil {
			logger.Error("failed to parse file",
				"err", err,
				"filename", inputFilename)
			return m, err
		}

		pointCount := 0
		for _, track := range gpxFile.Tracks {
			for _, segment := range track.Segments {

				insertCount, err := index.InsertGPXPoints(db, segment.Points)
				if err != nil {
					logger.Error("failed to save gpx point",
						"err", err,
						"filename", inputFilename)
					return m, err
				}

				pointCount += insertCount
			}
		}

		logger.Info(
			"imported file",
			"time", time.Since(startTime),
			"points", pointCount,
			"filename", inputFilename)

		return m, err
	}
}
