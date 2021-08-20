package dynamo

import (
	"fmt"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

var mediaRecordName = "media"

type mediaRecord struct {
	Pk          string  `json:"pk"`
	Sk          string  `json:"sk"`
	MediaKey    string  `json:"media_key"`
	Date        string  `json:"date"`
	Width       string  `json:"width"`
	Height      string  `json:"height"`
	CameraMake  string  `json:"camera_make"`
	CameraModel string  `json:"camera_model"`
	LocationLat float64 `json:"location_lat"`
	LocationLng float64 `json:"location_lng"`
}

func newMediaRecord(mediaMeta app.MediaMetadata) mediaRecord {
	mr := mediaRecord{}

	mr.Pk = newMediaRecordPK(mediaMeta)
	mr.Sk = newMediaRecordSK(mediaMeta)
	mr.MediaKey = mediaMeta.NewFilename()
	mr.Date = mediaMeta.Date.Format(time.RFC3339)
	mr.Width = mediaMeta.Width
	mr.Height = mediaMeta.Height
	mr.CameraMake = mediaMeta.CameraMake
	mr.CameraModel = mediaMeta.CameraModel
	mr.LocationLat = mediaMeta.Location.Coordinates.Lat
	mr.LocationLng = mediaMeta.Location.Coordinates.Lng

	return mr
}

func newMediaRecordPK(mediaMeta app.MediaMetadata) string {
	return mediaRecordName + "#" + mediaMeta.Hash
}

func newMediaRecordSK(mediaMeta app.MediaMetadata) string {
	return mediaRecordName
}

func NewIndexer(tableName, region string) app.Indexer {
	return func(mediaMeta app.MediaMetadata) error {
		mr := newMediaRecord(mediaMeta)

		fmt.Printf("%+v", mr)
		return nil
	}
}
