package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
	"go.uber.org/zap"
)

func NewHandler(logger *zap.SugaredLogger, geocodeMedia app.MediaGeocoder) func(ctx context.Context, req events.SQSEvent) error {
	return func(ctx context.Context, req events.SQSEvent) error {
		for _, record := range req.Records {
			_ = record.Body
		}
		return nil
	}
}

func main() {

	// conf
	mediaStoreBucket := "inari-mediastore-dev"
	thumbnailStoreBucket := "inari-thumbnailstore-dev"
	region := "eu-central-1"

	// deps
	zlogger, _ := zap.NewProduction()
	logger := zlogger.Sugar()
	defer logger.Sync()

	// deps
	downloader := s3.NewDownloader(mediaStoreBucket, region)
	uploader := s3.NewUploader(thumbnailStoreBucket, region)
	resizer := imgresize.NewResizer()
	createThumbnails := app.NewThumbnailer(downloader, resizer, uploader)

	lambda.Start(NewHandler(logger, createThumbnails))
}
