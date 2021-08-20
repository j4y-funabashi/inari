package app

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type importer = func(backupFilename string) (MediaMetadata, error)
type Downloader = func(backupFilename string) (string, error)
type Uploader = func(localFilename, mediaStoreFilename string) error
type Indexer = func(mediaMeta MediaMetadata) error
type MetadataExtractor = func(mediaFile string) (MediaMetadata, error)

type Coordinates struct {
	Lat float64
	Lng float64
}
type Location struct {
	Coordinates
}
type MediaMetadata struct {
	Hash        string
	Date        time.Time
	Location    Location
	Ext         string
	MimeType    string
	Width       string
	Height      string
	CameraMake  string
	CameraModel string
}

func (mm MediaMetadata) NewFilename() string {
	return fmt.Sprintf(
		"%s_%s.%s",
		mm.Date.Format("2006/20060102_150405"),
		mm.Hash,
		mm.Ext,
	)
}

func NewImporter(downloadFromBackup Downloader, extractMetadata MetadataExtractor, uploadToMediaStore Uploader, indexMedia Indexer) importer {
	return func(backupFilename string) (MediaMetadata, error) {

		// download file from backup storage
		downloadedFilename, err := downloadFromBackup(backupFilename)
		if err != nil {
			return MediaMetadata{}, err
		}
		logrus.
			WithField("backupFilename", backupFilename).
			WithField("filename", downloadedFilename).
			Info("downloaded media file from backup")

		// extract metadata
		mediaMeta, err := extractMetadata(downloadedFilename)
		if err != nil {
			return MediaMetadata{}, err
		}
		logrus.
			WithField("meta", mediaMeta).
			WithField("newFilename", mediaMeta.NewFilename()).
			Info("extracted metadata")

		// upload renamed file to media storage
		err = uploadToMediaStore(downloadedFilename, mediaMeta.NewFilename())
		if err != nil {
			return MediaMetadata{}, err
		}
		logrus.
			WithField("newFilename", mediaMeta.NewFilename()).
			Info("uploaded to mediastore")

		// index metadata in datastore
		err = indexMedia(mediaMeta)
		if err != nil {
			return MediaMetadata{}, err
		}
		return mediaMeta, nil
	}
}
