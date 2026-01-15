package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	CollectionTypeInbox         = "inbox"
	CollectionTypeCamera        = "camera"
	CollectionTypeTimelineMonth = "timeline_month"
	CollectionTypeTimelineDay   = "timeline_day"
	CollectionTypePlacesCountry = "places_country"
	CollectionTypePlacesRegion  = "places_region"
	CollectionTypeHashTag       = "hashtag"
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

type (
	Importer                = func(backupFilename string) (Media, error)
	Exporter                = func(mediaID string) error
	Thumbnailer             = func(mediastoreKey string) error
	QueryMediaDetail        = func(mediaID string) (Media, error)
	DeleteMedia             = func(mediaID string) error
	ExportMedia             = func(mediaID string) error
	UpdateMediaTextProperty = func(mediaID, caption string) error
	QueryNearestGPX         = func(cTime time.Time) (GPXPoint, error)
)

type (
	CollectionLister      func(collectionType string) ([]Collection, error)
	CollectionDetailQuery = func(collectionID string) (CollectionDetail, error)
	Resizer               = func(in, out string) (MediaSrc, error)
	Downloader            = func(backupFilename string) (string, error)
	Uploader              = func(localFilename, mediaStoreFilename string) error
	UploaderB             = func(sourceData []byte, mediaStoreFilename, contentType string) error
	Indexer               = func(media Media) (Media, error)
	Notifier              = func(mediaMeta Media) error
	FileLister            = func() ([]string, error)
	MetadataExtractor     = func(mediaFile string) (MediaMetadata, error)
	MediaDetailQuery      = func(mediaID string) (MediaDetailView, error)
	Geocoder              = func(lat, lng float64, cTime time.Time) (Location, error)
	LookupTimezone        = func(lat, lng float64, cTime time.Time) (string, error)
	MediaGeocoder         = func(mediaID string) (Location, error)
	LocationPutter        = func(mediaID string, location Location) error
	SaveGPXPoints         = func(points []GPXPoint) error
)

type Media struct {
	ID            string `json:"id,omitempty"`
	FilePath      string `json:"file_path,omitempty"`
	MediaMetadata `json:"media_metadata"`
	Thumbnails    MediaSrc     `json:"thumbnails"`
	Location      Location     `json:"location"`
	Collections   []Collection `json:"collections,omitempty"`
	FormattedDate string       `json:"date,omitempty"`
	Caption       string       `json:"caption,omitempty"`
	IsExported    bool         `json:"is_exported,omitempty"`
}

func (m Media) ToMicroformat() Microformat {
	category := []any{}
	for _, cat := range m.Collections {
		if cat.Type == CollectionTypeHashTag {
			category = append(category, cat.Title)
		}
	}

	photoURL := "https://media.funabashi.co.uk/" + m.Thumbnails.Large

	return Microformat{
		Type: []string{"h-entry"},
		Properties: map[string][]any{
			"uid":       {m.ID},
			"published": {m.FormattedDate},
			"content":   {m.Caption},
			"photo":     {photoURL},
			"category":  category,
			"location": {
				Microformat{
					Type: []string{"h-adr"},
					Properties: map[string][]any{
						"locality":     {m.Location.Locality},
						"region":       {m.Location.Region},
						"country-name": {m.Location.Country.Long},
						"geo": {
							Microformat{
								Type: []string{"h-geo"},
								Properties: map[string][]any{
									"latitude":  {strconv.FormatFloat(m.Location.Coordinates.Lat, 'f', -1, 64)},
									"longitude": {strconv.FormatFloat(m.Location.Coordinates.Lng, 'f', -1, 64)},
								},
							},
						},
					},
				},
			},
		},
	}
}

type Microformat struct {
	Type       []string         `json:"type"`
	Properties map[string][]any `json:"properties"`
}

// Collection types can be TIMELINE_MONTH
type Collection struct {
	ID            string `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	Type          string `json:"type,omitempty"`
	MediaCount    int    `json:"media_count,omitempty"`
	ExportedCount int    `json:"exported_count,omitempty"`
}

type GPXPoint struct {
	Timestamp time.Time
	Location
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
	Timezone    string `json:"timezone,omitempty"`
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

func (mm MediaMetadata) PostFilename() string {
	return fmt.Sprintf(
		"%s_%s.%s",
		mm.Date.Format("20060102_150405"),
		mm.Hash,
		"json",
	)
}

// ImportDir will check if backupFilename is a directory
// if it is a directory we will import all files with media extensions
func ImportDir(importFile Importer, logger Logger) func(backupFilename string) error {
	return func(backupFilename string) error {
		fInfo, err := os.Lstat(backupFilename)
		if err != nil {
			return err
		}
		if !fInfo.IsDir() {
			_, err := importFile(backupFilename)
			return err
		}

		// approx count files
		fileCount := 0
		filepath.Walk(
			backupFilename,
			func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				fileCount++
				return nil
			})
		logger.Info("importing files",
			"dir", backupFilename,
			"approx-count", fileCount,
		)

		filepath.Walk(
			backupFilename,
			func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				_, iErr := importFile(path)
				if iErr != nil {
					logger.Error("failed to import file", "err", iErr, "path", path)
				}
				return nil
			})

		return nil
	}
}

type MediaImporterConfig struct {
	FetchMediaDetail   QueryMediaDetail
	Logger             Logger
	DownloadFromBackup Downloader
	ExtractMetadata    MetadataExtractor
	UploadToMediaStore Uploader
	IndexMedia         Indexer
	CreateThumbnails   Resizer
	Geocode            Geocoder
	NotifyDownstream   Notifier
}

func NewImporter(config MediaImporterConfig) Importer {
	return func(inputFilename string) (Media, error) {
		startTime := time.Now()
		media := Media{}

		ext := strings.ToLower(filepath.Ext(inputFilename))
		if _, extValid := mediaExtensions[ext]; !extValid {
			return media, nil
		}

		// check media exists
		hash, err := parseHash(inputFilename)
		if err != nil {
			return media, fmt.Errorf("failed to parse media hash: %w", err)
		}
		existingMedia, _ := config.FetchMediaDetail(hash)
		if existingMedia.Hash == hash {
			return existingMedia, nil
		}

		// download file from backup storage
		tmpFilename, err := config.DownloadFromBackup(inputFilename)
		if err != nil {
			return media, fmt.Errorf("failed to download media from backup: %w", err)
		}
		defer os.Remove(tmpFilename)

		// extract metadata
		mediaMeta, err := config.ExtractMetadata(tmpFilename)
		if err != nil {
			return media, fmt.Errorf("failed to extract media metadata: %w", err)
		}
		media.MediaMetadata = mediaMeta
		media.Caption = mediaMeta.Title

		// upload renamed file to media storage
		err = config.UploadToMediaStore(tmpFilename, media.NewFilename())
		if err != nil {
			return media, fmt.Errorf("failed to upload to media store: %w", err)
		}
		media.FilePath = media.NewFilename()

		// create thumbnails
		thumbnails, err := config.CreateThumbnails(tmpFilename, media.NewFilename())
		if err != nil {
			return media, fmt.Errorf("failed to create thumbnails: %w", err)
		}
		media.Thumbnails = thumbnails

		// geocode
		loc, err := config.Geocode(media.Coordinates.Lat, media.Coordinates.Lng, media.Date)
		if err != nil {
			return media, fmt.Errorf("failed to geocode: %w", err)
		}
		media.Location = loc

		// index metadata in datastore
		media, err = config.IndexMedia(media)
		if err != nil {
			return media, fmt.Errorf("failed to index media metadata: %w", err)
		}

		// add to queue
		err = config.NotifyDownstream(media)
		if err != nil {
			config.Logger.Error("failed to notify downstream",
				"err", err,
				"backupFilename", inputFilename)
			return media, nil
		}

		config.Logger.Info("imported media",
			"path", inputFilename,
			"elapsedTime", time.Since(startTime),
		)

		return media, nil
	}
}

func parseHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
