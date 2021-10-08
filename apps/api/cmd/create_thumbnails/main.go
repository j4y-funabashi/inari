package main

import (
	"os"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.StandardLogger()

	mediaKey := os.Args[1]
	logger.
		WithField("arg", os.Args).
		WithField("mediaKey", mediaKey).
		Info("creating thumbnails for")

	// conf
	mediaStoreBucket := "inari-mediastore-dev"
	thumbnailStoreBucket := "inari-thumbnailstore-dev"
	region := "eu-central-1"

	// deps
	downloader := s3.NewDownloader(mediaStoreBucket, region)
	uploader := s3.NewUploader(thumbnailStoreBucket, region)
	resizer := imgresize.NewResizer()

	createThumbnails := app.NewThumbnailer(downloader, resizer, uploader)
	err := createThumbnails(mediaKey)
	if err != nil {
		logger.WithError(err).Error("failed to create thumbnails")
	}
}
