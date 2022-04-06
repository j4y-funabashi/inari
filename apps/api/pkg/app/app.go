package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Importer = func(backupFilename string) error
type Thumbnailer = func(mediastoreKey string) error
type ViewTimeline = func() (TimelineView, error)
type ViewTimelineMonth = func(monthID string) (TimelineMonthView, error)
type Resizer = func(imgFilename string) ([]string, error)
type Downloader = func(backupFilename string) (string, error)
type Uploader = func(localFilename, mediaStoreFilename string) error
type Indexer = func(mediaMeta MediaMetadata) error
type FileLister = func() ([]string, error)
type MetadataExtractor = func(mediaFile string) (MediaMetadata, error)
type TimelineQuery = func() (TimelineView, error)
type TimelineMonthQuery = func(monthID string) (TimelineMonthView, error)
type Geocoder = func(lat, lng float64) (Location, error)

type MediaMonth struct {
	ID         string
	Date       string `json:"date"`
	MediaCount int    `json:"media_count"`
}
type TimelineView struct {
	Months []MediaMonth `json:"months"`
}
type TimelineMonthView struct {
	CollectionMeta MediaMonth            `json:"collection_meta"`
	Media          []MediaCollectionItem `json:"media"`
}
type MediaCollectionItem struct {
	ID       string   `json:"id"`
	MimeType string   `json:"mime_type"`
	Date     string   `json:"date"`
	MediaSrc string   `json:"media_src"`
	Location Location `json:"location"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Location struct {
	Country     Country `json:"country"`
	Region      string  `json:"region"`
	Locality    string  `json:"locality"`
	Coordinates `json:"coordinates"`
}

type Country struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}
type MediaMetadata struct {
	Hash        string    `json:"hash"`
	Date        time.Time `json:"date"`
	Location    Location  `json:"location"`
	Ext         string    `json:"ext"`
	MimeType    string    `json:"mime_type"`
	Width       string    `json:"width"`
	Height      string    `json:"height"`
	CameraMake  string    `json:"camera_make"`
	CameraModel string    `json:"camera_model"`
	Keywords    string    `json:"keywords"`
	Title       string    `json:"title"`
}

func (mm MediaMetadata) NewFilename() string {
	return fmt.Sprintf(
		"%s_%s.%s",
		mm.Date.Format("2006/20060102_150405"),
		mm.Hash,
		mm.Ext,
	)
}

func (mm MediaMetadata) ThumbnailKey() string {
	return fmt.Sprintf(
		"%s_%s.%s",
		mm.Date.Format("20060102_150405"),
		mm.Hash,
		mm.Ext,
	)
}

func NewImporter(logger *zap.SugaredLogger, downloadFromBackup Downloader, extractMetadata MetadataExtractor, uploadToMediaStore Uploader, indexMedia Indexer) Importer {
	return func(backupFilename string) error {

		// download file from backup storage
		downloadedFilename, err := downloadFromBackup(backupFilename)
		if err != nil {
			return fmt.Errorf("failed to download media from backup: %w", err)
		}
		logger.Infow("downloaded media from backup",
			"backupFilename", backupFilename,
			"filename", downloadedFilename)
		defer os.Remove(downloadedFilename)

		// extract metadata
		mediaMeta, err := extractMetadata(downloadedFilename)
		if err != nil {
			return fmt.Errorf("failed to extract media metadata: %w", err)
		}
		logger.Infow("extracted metadata",
			"meta", mediaMeta,
			"newFilename", mediaMeta.NewFilename())

		// upload renamed file to media storage
		err = uploadToMediaStore(downloadedFilename, mediaMeta.NewFilename())
		if err != nil {
			return fmt.Errorf("failed to upload to media store: %w", err)
		}
		logger.Infow("uploaded to media store",
			"newFilename", mediaMeta.NewFilename())

		// index metadata in datastore
		err = indexMedia(mediaMeta)
		if err != nil {
			return fmt.Errorf("failed to index media metadata: %w", err)
		}
		return nil
	}
}

func NewThumbnailer(downloadFromMediaStore Downloader, resizeImage Resizer, uploadToThumbnailStore Uploader) Thumbnailer {
	return func(mediastoreKey string) error {
		// download file from media store
		downloadedFilename, err := downloadFromMediaStore(mediastoreKey)
		if err != nil {
			return err
		}
		logrus.
			WithField("mediastoreKey", mediastoreKey).
			WithField("filename", downloadedFilename).
			Info("downloaded media file from media store")

		thumbnailFiles, err := resizeImage(downloadedFilename)
		if err != nil {
			return err
		}
		logrus.
			WithField("mediastoreKey", mediastoreKey).
			WithField("filename", downloadedFilename).
			WithField("thumbnailFiles", thumbnailFiles).
			Info("resized image")

		for _, thumbnailFile := range thumbnailFiles {
			err := uploadToThumbnailStore(thumbnailFile, "thmnb/"+filepath.Base(thumbnailFile))
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func NewTimelineView(timelineQuery TimelineQuery) ViewTimeline {
	return func() (TimelineView, error) {
		timelineView, err := timelineQuery()
		if err != nil {
			return TimelineView{}, err
		}
		return timelineView, nil
	}
}

func NewTimelineMonthView(timelineQuery TimelineMonthQuery) ViewTimelineMonth {
	return func(monthID string) (TimelineMonthView, error) {
		timelineView, err := timelineQuery(monthID)
		if err != nil {
			return TimelineMonthView{}, nil
		}
		return timelineView, nil
	}
}

func NewGeocoder(logger *zap.SugaredLogger, reverseGeocode Geocoder) Geocoder {
	return func(lat, lng float64) (Location, error) {
		loc, err := reverseGeocode(lat, lng)
		if err != nil {
			logger.Errorw(
				"failed to geocode location",
				"lat", lat,
				"lng", lng,
			)
		}
		return loc, err
	}
}
