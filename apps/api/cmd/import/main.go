package main

import (
	"os"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.StandardLogger()

	mediaFilename := os.Args[1]
	logger.
		WithField("arg", os.Args).
		WithField("mediaFilename", mediaFilename).
		Info("importing")

	// conf
	bucket := "backup.funabashi"
	mediaStoreBucket := "inari-mediastore-dev"
	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"

	// deps
	downloader := s3.NewDownloader(bucket, region)
	uploader := s3.NewUploader(mediaStoreBucket, region)
	indexer := dynamo.NewIndexer(mediaStoreTableName, region)
	extractMetadata := exiftool.NewExtractor()
	importMedia := app.NewImporter(downloader, extractMetadata, uploader, indexer)

	result, err := importMedia(mediaFilename)
	if err != nil {
		logger.
			WithError(err).
			WithField("mediaFilename", mediaFilename).
			Error("failed to import")
		os.Exit(1)
	}
	logger.WithField("res", result.NewFilename())
}
