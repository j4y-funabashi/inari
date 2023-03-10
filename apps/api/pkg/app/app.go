package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
)

const (
	CollectionTypeInbox         = "inbox"
	CollectionTypeTimelineMonth = "timeline_month"
	CollectionTypeTimelineDay   = "timeline_day"
	CollectionTypePlacesCountry = "places_country"
	CollectionTypePlacesRegion  = "places_region"
)

type Logger interface {
	Info(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}
type NullLogger struct{}

func (NullLogger) Info(msg string, ctx ...interface{})  {}
func (NullLogger) Error(msg string, ctx ...interface{}) {}

func NewNullLogger() Logger {
	return NullLogger{}
}

type Importer = func(backupFilename string) (Media, error)
type Thumbnailer = func(mediastoreKey string) error
type QueryMediaDetail = func(mediaID string) (Media, error)

type CollectionLister func(collectionType string) ([]Collection, error)
type ViewTimelineMonth = func(monthID string) (TimelineMonthView, error)
type Resizer = func(in, out string) (MediaSrc, error)
type Downloader = func(backupFilename string) (string, error)
type Uploader = func(localFilename, mediaStoreFilename string) error
type Indexer = func(media Media) (Media, error)
type Notifier = func(mediaMeta Media) error
type FileLister = func() ([]string, error)
type MetadataExtractor = func(mediaFile string) (MediaMetadata, error)
type TimelineMonthQuery = func(monthID string) (TimelineMonthView, error)
type MediaDetailQuery = func(mediaID string) (MediaDetailView, error)
type Geocoder = func(lat, lng float64) (Location, error)
type MediaGeocoder = func(mediaID string) (Location, error)
type LocationPutter = func(mediaID string, location Location) error

type Media struct {
	ID       string
	FilePath string
	MediaMetadata
	Thumbnails  MediaSrc
	Location    Location
	Collections []Collection
}

// Collection types can be TIMELINE_MONTH
type Collection struct {
	ID         string `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Type       string `json:"type,omitempty"`
	MediaCount int    `json:"media_count,omitempty"`
}

type MediaMonth struct {
	ID         string
	Date       string `json:"date"`
	MediaCount int    `json:"media_count"`
}
type TimelineView struct {
	Months []Collection `json:"months"`
}
type TimelineMonthView struct {
	CollectionMeta Collection            `json:"collection_meta"`
	Media          []MediaCollectionItem `json:"media"`
}

type MediaDetailView struct {
	Media MediaCollectionItem `json:"media"`
}

type MediaSrc struct {
	Key    string `json:"key"`
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
}

type MediaCollectionItem struct {
	ID       string   `json:"id"`
	Date     string   `json:"date"`
	MediaSrc MediaSrc `json:"media_src"`
	Caption  string   `json:"caption"`
	MediaMetadata
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
	Hash        string      `json:"hash"`
	Date        time.Time   `json:"date"`
	Coordinates Coordinates `json:"coordinates"`
	Ext         string      `json:"ext"`
	MimeType    string      `json:"mime_type"`
	Width       string      `json:"width"`
	Height      string      `json:"height"`
	CameraMake  string      `json:"camera_make"`
	CameraModel string      `json:"camera_model"`
	Keywords    string      `json:"keywords"`
	Title       string      `json:"title"`
}

// file extensions inari will import
var mediaExtensions = map[string]bool{
	".jpg": true,
	".mov": true,
	".mp4": true,
	".avi": true,
}

func (mm MediaMetadata) NewFilename() string {
	return fmt.Sprintf(
		"%s/%s_%s.%s",
		mm.Date.Format("2006"),
		mm.Date.Format("20060102_150405"),
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

// ImportDir will check if backupFilename is a directory
// if it is a directory we will import all files with media extensions
func ImportDir(importFile Importer, logger log.Logger) func(backupFilename string) error {
	return func(backupFilename string) error {
		fInfo, err := os.Lstat(backupFilename)
		if err != nil {
			return err
		}
		if !fInfo.IsDir() {
			_, err := importFile(backupFilename)
			return err
		}

		filepath.Walk(
			backupFilename,
			func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				_, iErr := importFile(path)
				if iErr != nil {
					logger.Error("failed to import file", "err", iErr)
				}
				return nil
			})
		return nil
	}
}

func NewImporter(logger Logger, downloadFromBackup Downloader, extractMetadata MetadataExtractor, uploadToMediaStore Uploader, indexMedia Indexer, createThumbnails Resizer, geocode Geocoder, notifyDownstream Notifier) Importer {
	return func(inputFilename string) (Media, error) {
		media := Media{}

		ext := strings.ToLower(filepath.Ext(inputFilename))
		if _, extValid := mediaExtensions[ext]; !extValid {
			return media, nil
		}

		// download file from backup storage
		tmpFilename, err := downloadFromBackup(inputFilename)
		if err != nil {
			return media, fmt.Errorf("failed to download media from backup: %w", err)
		}
		defer os.Remove(tmpFilename)

		// extract metadata
		mediaMeta, err := extractMetadata(tmpFilename)
		if err != nil {
			return media, fmt.Errorf("failed to extract media metadata: %w", err)
		}
		media.MediaMetadata = mediaMeta

		// upload renamed file to media storage
		err = uploadToMediaStore(tmpFilename, media.NewFilename())
		if err != nil {
			return media, fmt.Errorf("failed to upload to media store: %w", err)
		}
		media.FilePath = media.NewFilename()

		// create thumbnails
		thumbnails, err := createThumbnails(tmpFilename, media.NewFilename())
		if err != nil {
			return media, fmt.Errorf("failed to create thumbnails: %w", err)
		}
		media.Thumbnails = thumbnails

		// geocode
		loc, err := geocode(media.Coordinates.Lat, media.Coordinates.Lng)
		if err != nil {
			return media, fmt.Errorf("failed to geocode: %w", err)
		}
		media.Location = loc

		// index metadata in datastore
		media, err = indexMedia(media)
		if err != nil {
			return media, fmt.Errorf("failed to index media metadata: %w", err)
		}

		// add to queue
		err = notifyDownstream(media)
		if err != nil {
			logger.Error("failed to notify downstream",
				"err", err,
				"backupFilename", inputFilename)
			return media, nil
		}

		logger.Info("imported media",
			"backupFilename", inputFilename,
			"downloadedFilename", tmpFilename,
			"mediaMeta", media,
			"newFilename", media.NewFilename())

		return media, nil
	}
}

func NewTimelineMonthView(timelineQuery TimelineMonthQuery) ViewTimelineMonth {
	return func(monthID string) (TimelineMonthView, error) {
		timelineView, err := timelineQuery(monthID)
		if err != nil {
			return TimelineMonthView{}, err
		}
		return timelineView, nil
	}
}
