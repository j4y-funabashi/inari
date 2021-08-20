package dynamo

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

type mediaDateRecord struct {
	Pk       string `json:"pk"`
	Sk       string `json:"sk"`
	Gsi1pk   string `json:"gsi1pk"`
	Gsi1sk   string `json:"gsi1sk"`
	MediaKey string `json:"media_key"`
	Date     string `json:"date"`
}

func newMediaDateRecord(mediaMeta app.MediaMetadata) mediaDateRecord {
	mdr := mediaDateRecord{}

	mdr.Pk = newMediaRecordPK(mediaMeta)
	mdr.Sk = newMediaDateRecordSK(mediaMeta)
	mdr.Gsi1pk = "mediaDate#" + mediaMeta.Date.Format("2006")
	mdr.Gsi1sk = newMediaDateRecordSK(mediaMeta)
	mdr.MediaKey = mediaMeta.NewFilename()
	mdr.Date = mediaMeta.Date.Format(time.RFC3339)

	return mdr
}

func newMediaDateRecordSK(mediaMeta app.MediaMetadata) string {
	return "mediaDate" + "#" + mediaMeta.Date.Format(time.RFC3339) + "#" + mediaMeta.Hash
}

func NewIndexer(tableName, region string) app.Indexer {
	return func(mediaMeta app.MediaMetadata) error {

		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		client := dynamodb.New(sess)

		// save media record
		mr := newMediaRecord(mediaMeta)
		mrItem, err := dynamodbattribute.MarshalMap(mr)
		if err != nil {
			return err
		}
		fmt.Printf("%+v", mrItem)
		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      mrItem,
		})
		if err != nil {
			return err
		}

		// save media date record
		mdr := newMediaDateRecord(mediaMeta)
		mdrItem, err := dynamodbattribute.MarshalMap(mdr)
		if err != nil {
			return err
		}
		fmt.Printf("%+v", mdrItem)
		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      mdrItem,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
