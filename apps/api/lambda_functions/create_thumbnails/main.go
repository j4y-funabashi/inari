package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
	"go.uber.org/zap"
)

func NewHandler(logger *zap.SugaredLogger, createThumbnails app.Thumbnailer) func(ctx context.Context, req events.SQSEvent) error {
	return func(ctx context.Context, req events.SQSEvent) error {

		for _, record := range req.Records {
			mediaID := record.Body

			err := createThumbnails(mediaID)
			if err != nil {
				logger.
					Errorw("failed to create thumbnails",
						"err", err,
						"mediaKey", mediaID)
				return err
			}

			logger.
				Infow("created thumbnails",
					"mediaKey", mediaID)

		}

		return nil
	}
}

func main() {

	// conf
	mediaStoreBucket := "inari-mediastore-dev"
	thumbnailStoreBucket := "inari-thumbnailstore-dev"
	region := "eu-central-1"
	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"

	// deps
	zlogger, _ := zap.NewProduction()
	logger := zlogger.Sugar()
	defer logger.Sync()

	// -- create client
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	dynamoClient := dynamodb.New(sess)

	// deps
	fetchMedia := dynamo.NewMediaDetailQuery(mediaStoreTableName, dynamoClient)
	downloader := s3.NewDownloader(mediaStoreBucket, region)
	uploader := s3.NewUploader(thumbnailStoreBucket, region)
	resizer := imgresize.NewResizer()
	createThumbnails := app.NewThumbnailer(fetchMedia, downloader, resizer, uploader)

	lambda.Start(NewHandler(logger, createThumbnails))
}
