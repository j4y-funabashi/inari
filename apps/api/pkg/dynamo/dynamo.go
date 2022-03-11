package dynamo

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

var mediaRecordName = "media"
var mediaDateRecordName = "mediaDate"
var collectionMediaDayRecordName = "collectionMediaDay"
var collectionMonthPrefix = "month"

type baseMediaRecordMeta struct {
	MediaKey string `json:"media_key"`
	MimeType string `json:"mime_type"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Date     string `json:"date"`
}

type mediaRecord struct {
	Pk          string  `json:"pk"`
	Sk          string  `json:"sk"`
	Hash        string  `json:"hash"`
	CameraMake  string  `json:"camera_make"`
	CameraModel string  `json:"camera_model"`
	LocationLat float64 `json:"location_lat"`
	LocationLng float64 `json:"location_lng"`
	Ext         string  `json:"ext"`
	Keywords    string  `json:"keywords"`
	Title       string  `json:"title"`
	baseMediaRecordMeta
}

type mediaDateCollectionRecord struct {
	Pk         string `json:"pk"`
	Sk         string `json:"sk"`
	Date       string `json:"media_date"`
	Gsi1pk     string `json:"gsi1pk"`
	Gsi1sk     string `json:"gsi1sk"`
	MediaCount int    `json:"media_count"`
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
	mr.MimeType = mediaMeta.MimeType
	mr.Hash = mediaMeta.Hash
	mr.Ext = mediaMeta.Ext
	mr.Keywords = mediaMeta.Keywords
	mr.Title = mediaMeta.Title

	return mr
}

func newMediaRecordPK(mediaMeta app.MediaMetadata) string {
	return collectionMonthPrefix + "#" + mediaMeta.Date.Format("2006-01")
}

func newMediaRecordSK(mediaMeta app.MediaMetadata) string {
	return mediaRecordName + "#" + mediaMeta.NewFilename()
}

func newMediaDateCollectionRecord(mediaMeta app.MediaMetadata) mediaDateCollectionRecord {
	mdr := mediaDateCollectionRecord{}

	mdr.Pk = newMediaRecordPK(mediaMeta)
	mdr.Sk = newMediaDateCollectionRecordSK(mediaMeta)
	mdr.Date = mediaMeta.Date.Format("2006-01")
	mdr.Gsi1pk = "monthCollection"
	mdr.Gsi1sk = newMediaDateCollectionRecordSK(mediaMeta)

	return mdr
}

type mediaDateCollectionKey struct {
	Pk string `json:"pk"`
	Sk string `json:"sk"`
}

func newMediaDateCollectionKey(mediaMeta app.MediaMetadata) mediaDateCollectionKey {
	mdr := mediaDateCollectionKey{}

	mdr.Pk = newMediaRecordPK(mediaMeta)
	mdr.Sk = newMediaDateCollectionRecordSK(mediaMeta)

	return mdr
}

type mediaDateCollectionUpdate struct {
	Date       string `json:":media_date"`
	Gsi1pk     string `json:":gsi1pk"`
	Gsi1sk     string `json:":gsi1sk"`
	MediaCount int    `json:":media_count"`
}

func newMediaDateCollectionUpdate(mediaMeta app.MediaMetadata) mediaDateCollectionUpdate {
	out := mediaDateCollectionUpdate{}
	out.Date = mediaMeta.Date.Format("2006-01")
	out.Gsi1pk = "monthCollection"
	out.Gsi1sk = newMediaDateCollectionRecordSK(mediaMeta)
	out.MediaCount = 1
	return out
}

func newMediaDateCollectionRecordSK(mediaMeta app.MediaMetadata) string {
	return "META#" + mediaMeta.Date.Format("2006-01")
}

func NewIndexer(tableName, region string) app.Indexer {
	return func(mediaMeta app.MediaMetadata) error {

		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		client := dynamodb.New(sess)

		// -- save media record
		mr := newMediaRecord(mediaMeta)
		mrItem, err := dynamodbattribute.MarshalMap(mr)
		if err != nil {
			return err
		}

		// -- save media date collection
		mdckeyItem, err := dynamodbattribute.MarshalMap(newMediaDateCollectionKey(mediaMeta))
		if err != nil {
			return err
		}
		updateValues, err := dynamodbattribute.MarshalMap(newMediaDateCollectionUpdate(mediaMeta))
		if err != nil {
			return err
		}

		_, err = client.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName:           aws.String(tableName),
						Item:                mrItem,
						ConditionExpression: aws.String("attribute_not_exists(pk)"),
					},
				},
				{
					Update: &dynamodb.Update{
						TableName:                 aws.String(tableName),
						Key:                       mdckeyItem,
						UpdateExpression:          aws.String("SET media_date = :media_date, gsi1pk = :gsi1pk, gsi1sk = :gsi1sk ADD media_count :media_count"),
						ExpressionAttributeValues: updateValues,
					},
				},
			},
		})
		if err != nil {
			switch t := err.(type) {
			case *dynamodb.TransactionCanceledException:
				for _, r := range t.CancellationReasons {
					if *r.Code == "ConditionalCheckFailed" {
						return nil
					}
				}
			default:
				return err
			}
		}

		return err
	}
}

func NewTimelineQuery(tableName, region string) app.TimelineQuery {
	return func() (app.TimelineView, error) {

		// -- create client
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		client := dynamodb.New(sess)

		timelineView := app.TimelineView{}

		// -- query dynamo
		keyValues := map[string]string{
			":pk": "monthCollection",
		}
		eavalues, err := dynamodbattribute.MarshalMap(keyValues)
		if err != nil {
			return timelineView, err
		}
		res, err := client.Query(&dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			KeyConditionExpression:    aws.String("gsi1pk = :pk"),
			ExpressionAttributeValues: eavalues,
			ScanIndexForward:          aws.Bool(false),
			IndexName:                 aws.String("GSI1"),
		})
		if err != nil {
			return timelineView, err
		}

		for _, item := range res.Items {
			mdr := mediaDateCollectionRecord{}
			err = dynamodbattribute.UnmarshalMap(item, &mdr)
			if err != nil {
				return timelineView, err
			}

			// -- convert media record to media day
			dat, err := time.Parse("2006-01", mdr.Date)
			if err != nil {
				return timelineView, err
			}
			mediaMonth := app.MediaMonth{
				ID:   mdr.Date,
				Date: dat.Format("2006 Jan"),
			}
			timelineView.Months = append(timelineView.Months, mediaMonth)
		}

		return timelineView, nil
	}
}
