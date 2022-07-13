package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.StandardLogger()

	mediaID := os.Args[1]
	logger.
		WithField("arg", os.Args).
		WithField("mediaKey", mediaID).
		Info("creating thumbnails for")

	region := "eu-central-1"
	// -- create client
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	dynamoClient := dynamodb.New(sess)
	s3Downloader := s3manager.NewDownloader(sess)
	s3Uploader := s3manager.NewUploader(sess)
	s3Client := s3.New(sess)

	// conf
	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	mediaStoreBucket := "inari-mediastore-dev"
	thumbnailStoreBucket := "inari-thumbnailstore-dev"

	// deps
	fetchMedia := dynamo.NewMediaDetailQuery(mediaStoreTableName, dynamoClient)
	downloader := storage.NewDownloader(mediaStoreBucket, s3Downloader)
	uploader := storage.NewUploader(thumbnailStoreBucket, s3Uploader, s3Client)
	resizer := imgresize.NewResizer()

	createThumbnails := app.NewThumbnailer(fetchMedia, downloader, resizer, uploader)
	err := createThumbnails(mediaID)
	if err != nil {
		logger.WithError(err).Error("failed to create thumbnails")
	}
}
