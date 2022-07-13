package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/notify"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
	"go.uber.org/zap"
)

func main() {
	zlogger, _ := zap.NewDevelopment()
	logger := zlogger.Sugar()
	defer logger.Sync()

	// conf
	backupBucket := "backup.funabashi"
	mediaStoreBucket := "inari-mediastore-dev"
	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	topicARN := "arn:aws:sns:eu-central-1:725941804651:PostImportTopic"
	region := "eu-central-1"

	// aws clients
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	dynamoClient := dynamodb.New(sess)
	s3Downloader := s3manager.NewDownloader(sess)
	snsClient := sns.New(sess)
	s3Uploader := s3manager.NewUploader(sess)
	s3Client := s3.New(sess)

	// deps
	downloader := storage.NewDownloader(backupBucket, s3Downloader)
	uploader := storage.NewUploader(mediaStoreBucket, s3Uploader, s3Client)
	indexer := dynamo.NewIndexer(mediaStoreTableName, dynamoClient)
	extractMetadata := exiftool.NewExtractor("/usr/bin/exiftool")
	notifier := notify.NewNotifier(snsClient, topicARN)
	importMedia := app.NewImporter(logger, downloader, extractMetadata, uploader, indexer, notifier)

	if len(os.Args) > 1 {
		inputFilename := os.Args[1]
		if inputFilename != "" {
			err := importMedia(inputFilename)
			if err != nil {
				logger.Errorw("failed to import",
					"error", err,
					"inputFilename", inputFilename)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	listFiles := storage.NewLister(backupBucket, region, "jayr")
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
