package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
	"go.uber.org/zap"
)

func NewHandler(logger *zap.SugaredLogger, importMedia app.Importer) func(ctx context.Context, req events.SQSEvent) error {
	return func(ctx context.Context, req events.SQSEvent) error {

		for _, record := range req.Records {
			mediaFilename := record.Body
			logger.
				Infow("importing",
					"mediaFilename", mediaFilename)

			err := importMedia(mediaFilename)
			if err != nil {
				logger.
					Errorw("failed to import",
						"err", err,
						"mediaFilename", mediaFilename)
				return err
			}

		}

		return nil
	}
}

func listOptDir(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	// conf
	bucket := "backup.funabashi"
	mediaStoreBucket := "inari-mediastore-dev"
	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"
	exiftoolPath := "/opt/bin/perl /opt/bin/exiftool"
	listOptDir("/opt/")

	// deps
	zlogger, _ := zap.NewProduction()
	logger := zlogger.Sugar()
	defer logger.Sync()

	downloader := s3.NewDownloader(bucket, region)
	uploader := s3.NewUploader(mediaStoreBucket, region)
	indexer := dynamo.NewIndexer(mediaStoreTableName, region)
	extractMetadata := exiftool.NewExtractor(exiftoolPath)
	importMedia := app.NewImporter(logger, downloader, extractMetadata, uploader, indexer)

	lambda.Start(NewHandler(logger, importMedia))
}
