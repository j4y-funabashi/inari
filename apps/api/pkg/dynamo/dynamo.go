package dynamo

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

var mediaRecordName = "media"
var mediaDateRecordName = "mediaDate"
var collectionMediaDayRecordName = "collectionMediaDay"

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

type mediaDateRecord struct {
	Pk     string `json:"pk"`
	Sk     string `json:"sk"`
	Gsi1pk string `json:"gsi1pk"`
	Gsi1sk string `json:"gsi1sk"`
	baseMediaRecordMeta
}

type mediaDateCollectionRecord struct {
	Pk        string   `json:"pk"`
	Sk        string   `json:"sk"`
	Date      string   `json:"date"`
	MediaList []string `json:"media_list"`
}

type mediaDateCollectionUpdate struct {
	MediaKey  string   `json:":mk"`
	MediaList []string `json:":mkl"`
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
	return mediaRecordName + "#" + mediaMeta.Hash
}

func newMediaRecordSK(mediaMeta app.MediaMetadata) string {
	return mediaRecordName
}

func newMediaDateRecord(mediaMeta app.MediaMetadata) mediaDateRecord {
	mdr := mediaDateRecord{}

	mdr.Pk = newMediaRecordPK(mediaMeta)
	mdr.Sk = mediaDateRecordName
	mdr.Gsi1pk = mediaDateRecordName + "#" + mediaMeta.Date.Format("2006")
	mdr.Gsi1sk = newMediaDateRecordSK(mediaMeta)
	mdr.MediaKey = mediaMeta.NewFilename()
	mdr.MimeType = mediaMeta.MimeType
	mdr.Width = mediaMeta.Width
	mdr.Height = mediaMeta.Height
	mdr.Date = mediaMeta.Date.Format(time.RFC3339)

	return mdr
}

func newMediaDateRecordSK(mediaMeta app.MediaMetadata) string {
	return mediaDateRecordName + "#" + mediaMeta.Date.Format(time.RFC3339) + "#" + mediaMeta.Hash
}

func newMediaDateCollectionRecord(mediaMeta app.MediaMetadata) mediaDateCollectionRecord {
	mdr := mediaDateCollectionRecord{}

	mdr.Pk = collectionMediaDayRecordName
	mdr.Sk = newMediaDateCollectionRecordSK(mediaMeta)
	mdr.Date = mediaMeta.Date.Format("2006-01-02")
	mdr.MediaList = append(mdr.MediaList, newCollectionMediaListItem(mediaMeta))

	return mdr
}

func newMediaDateCollectionUpdate(mediaMeta app.MediaMetadata) mediaDateCollectionUpdate {
	mdr := mediaDateCollectionUpdate{}

	mdr.MediaKey = newCollectionMediaListItem(mediaMeta)
	mdr.MediaList = append(mdr.MediaList, newCollectionMediaListItem(mediaMeta))

	return mdr
}

func newCollectionMediaListItem(mediaMeta app.MediaMetadata) string {
	return fmt.Sprintf(
		"%s##%s##%s##%s",
		mediaMeta.Hash,
		mediaMeta.MimeType,
		mediaMeta.Date.Format(time.RFC3339),
		mediaMeta.ThumbnailKey(),
	)
}

func convertMediaListItemToMediaCollectionItem(mediaListItem string) app.MediaCollectionItem {
	split := strings.Split(mediaListItem, "##")
	item := app.MediaCollectionItem{
		ID:       split[0],
		MimeType: split[1],
		Date:     split[2],
		MediaSrc: split[3],
	}

	return item
}

func newMediaDateCollectionRecordSK(mediaMeta app.MediaMetadata) string {
	return collectionMediaDayRecordName + "#" + mediaMeta.Date.Format("2006-01-02")
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
		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      mrItem,
		})
		if err != nil {
			return err
		}

		// -- save media date record
		mdr := newMediaDateRecord(mediaMeta)
		mdrItem, err := dynamodbattribute.MarshalMap(mdr)
		if err != nil {
			return err
		}
		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      mdrItem,
		})
		if err != nil {
			return err
		}

		// -- save media date collection if it does not exist
		mdrcoll := newMediaDateCollectionRecord(mediaMeta)
		mdrcollItem, err := dynamodbattribute.MarshalMap(mdrcoll)
		if err != nil {
			return err
		}
		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName:           aws.String(tableName),
			Item:                mdrcollItem,
			ConditionExpression: aws.String("attribute_not_exists(pk)"),
		})
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() != "ConditionalCheckFailedException" {
					return err
				}
			} else {
				return err
			}
		}

		// -- update media date collection mediaList
		mdrcollUpdate, err := dynamodbattribute.MarshalMap(newMediaDateCollectionUpdate(mediaMeta))
		mdrcollUpdateKey, err := dynamodbattribute.MarshalMap(
			struct {
				Pk string `json:"pk"`
				Sk string `json:"sk"`
			}{
				Pk: collectionMediaDayRecordName,
				Sk: newMediaDateCollectionRecordSK(mediaMeta),
			},
		)
		if err != nil {
			return err
		}
		_, err = client.UpdateItem(&dynamodb.UpdateItemInput{
			ExpressionAttributeValues: mdrcollUpdate,
			TableName:                 aws.String(tableName),
			UpdateExpression:          aws.String("SET media_list = list_append(media_list, :mkl)"),
			ConditionExpression:       aws.String("not contains(media_list, :mk)"),
			Key:                       mdrcollUpdateKey,
		})
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() != "ConditionalCheckFailedException" {
					return err
				}
			} else {
				return err
			}
		}

		return nil
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
			":pk": "collectionMediaDay",
		}
		eavalues, err := dynamodbattribute.MarshalMap(keyValues)
		if err != nil {
			return timelineView, err
		}
		res, err := client.Query(&dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			KeyConditionExpression:    aws.String("pk = :pk"),
			ExpressionAttributeValues: eavalues,
			ScanIndexForward:          aws.Bool(false),
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

			// build media list
			mediaList := []app.MediaCollectionItem{}
			for _, mli := range mdr.MediaList {
				mediaList = append(mediaList, convertMediaListItemToMediaCollectionItem(mli))
			}
			// -- convert media record to media day
			mediaDay := app.MediaDay{
				Date:  mdr.Date,
				Media: mediaList,
			}
			timelineView.Days = append(timelineView.Days, mediaDay)
		}

		return timelineView, nil
	}
}
