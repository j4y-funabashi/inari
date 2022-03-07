package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
	"github.com/sirupsen/logrus"
)

func NewHandler(importMedia app.Importer) func(ctx context.Context, req events.SQSEvent) error {
	return func(ctx context.Context, req events.SQSEvent) error {
		logger := logrus.StandardLogger()

		for _, record := range req.Records {
			mediaFilename := record.Body
			logger.
				WithField("mediaFilename", mediaFilename).
				Info("importing")

			err := importMedia(mediaFilename)
			if err != nil {
				logger.
					WithError(err).
					WithField("mediaFilename", mediaFilename).
					Error("failed to import")
				return err
			}

		}

		return nil
	}
}

func main() {

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

	lambda.Start(NewHandler(importMedia))
}
