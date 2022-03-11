package main

import (
	"os"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
	"go.uber.org/zap"
)

func main() {
	zlogger, _ := zap.NewProduction()
	logger := zlogger.Sugar()
	defer logger.Sync()

	// conf
	backupBucket := "backup.funabashi"
	mediaStoreBucket := "inari-mediastore-dev"
	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"

	// deps
	downloader := s3.NewDownloader(backupBucket, region)
	uploader := s3.NewUploader(mediaStoreBucket, region)
	indexer := dynamo.NewIndexer(mediaStoreTableName, region)
	extractMetadata := exiftool.NewExtractor("/usr/bin/exiftool")
	importMedia := app.NewImporter(logger, downloader, extractMetadata, uploader, indexer)
	listFiles := s3.NewLister(backupBucket, region, "jayr")
	files, err := listFiles()
	if err != nil {
		logger.Errorw("failed to list files",
			"error", err)
	}

	logger.Infow("listed files",
		"files", len(files))

	for _, mediaFilename := range files {
		err = importMedia(mediaFilename)
		if err != nil {
			logger.Errorw("failed to import",
				"error", err,
				"mediaFilename", mediaFilename)
			os.Exit(1)
		}
	}
}
