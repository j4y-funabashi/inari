package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/inconshreveable/log15"
)

const (
	CollectionTypeInbox         = "inbox"
	CollectionTypeCamera        = "camera"
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
type DeleteMedia = func(mediaID string) error

type CollectionLister func(collectionType string) ([]Collection, error)
type CollectionDetailQuery = func(collectionID string) (CollectionDetail, error)
type Resizer = func(in, out string) (MediaSrc, error)
type Downloader = func(backupFilename string) (string, error)
type Uploader = func(localFilename, mediaStoreFilename string) error
type Indexer = func(media Media) (Media, error)
type Notifier = func(mediaMeta Media) error
type FileLister = func() ([]string, error)
type MetadataExtractor = func(mediaFile string) (MediaMetadata, error)
type MediaDetailQuery = func(mediaID string) (MediaDetailView, error)
type Geocoder = func(lat, lng float64) (Location, error)
type MediaGeocoder = func(mediaID string) (Location, error)
type LocationPutter = func(mediaID string, location Location) error

type Media struct {
	ID            string `json:"id,omitempty"`
	FilePath      string `json:"file_path,omitempty"`
	MediaMetadata `json:"media_metadata,omitempty"`
	Thumbnails    MediaSrc     `json:"thumbnails,omitempty"`
	Location      Location     `json:"location,omitempty"`
	Collections   []Collection `json:"collections,omitempty"`
	FormattedDate string       `json:"date,omitempty"`
	Caption       string       `json:"caption,omitempty"`
}

// Collection types can be TIMELINE_MONTH
type Collection struct {
	ID         string `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Type       string `json:"type,omitempty"`
	MediaCount int    `json:"media_count,omitempty"`
}

type CollectionDetail struct {
	CollectionMeta Collection `json:"collection_meta"`
	Media          []Media    `json:"media"`
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
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}
type Location struct {
	Country     Country `json:"country,omitempty"`
	Region      string  `json:"region,omitempty"`
	Locality    string  `json:"locality,omitempty"`
	Coordinates `json:"coordinates,omitempty"`
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

		var wg sync.WaitGroup
		workerCount := 50
		jobs := make(chan string)
		results := make(chan Media)
		for w := 1; w <= workerCount; w++ {
			wg.Add(1)
			w := w
			go func() {
				defer wg.Done()
				worker(importFile, logger, w, jobs, results)
			}()
		}
		go processResults(logger, results)

		filepath.Walk(
			backupFilename,
			func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				jobs <- path
				return nil
			})
		close(jobs)

		wg.Wait()

		return nil
	}
}

func worker(importFile Importer, logger log.Logger, id int, jobs <-chan string, results chan<- Media) {
	for path := range jobs {
		media, iErr := importFile(path)
		if iErr != nil {
			logger.Error("failed to import file", "err", iErr, "path", path)
		}
		results <- media
	}
}

func processResults(logger log.Logger, results <-chan Media) {
	for media := range results {
		logger.Info("imported media",
			"media", media.ID,
			"newFilename", media.NewFilename())
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
		media.Caption = mediaMeta.Title

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

		return media, nil
	}
}
