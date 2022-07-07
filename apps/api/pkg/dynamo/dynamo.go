package dynamo

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/imgresize"
)

const (
	idSeperator            = "--"
	collectionRecordPrefix = "collection"
	mediaRecordPrefix      = "media"
	collectionMonthPrefix  = "month"
)

type mediaRecord struct {
	Pk                string  `json:"pk"`
	Sk                string  `json:"sk"`
	Hash              string  `json:"hash"`
	CameraMake        string  `json:"camera_make"`
	CameraModel       string  `json:"camera_model"`
	LocationLat       float64 `json:"location_lat"`
	LocationLng       float64 `json:"location_lng"`
	LocationRegion    string  `json:"location_region"`
	LocationLocality  string  `json:"location_locality"`
	LocationCountryL  string  `json:"location_country_l"`
	LocationCountrySh string  `json:"location_country_sh"`
	Ext               string  `json:"ext"`
	Keywords          string  `json:"keywords"`
	Title             string  `json:"title"`
	MediaKey          string  `json:"media_key"`
	MimeType          string  `json:"mime_type"`
	Width             string  `json:"width"`
	Height            string  `json:"height"`
	Date              string  `json:"date"`
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
	mr.Sk = newMediaRecordPK(mediaMeta)
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

func newMediaFromMediaRecord(mr mediaRecord) app.MediaCollectionItem {
	m := app.MediaCollectionItem{}
	m.ID = app.MediaCollectionID{
		CollectionID: mr.Pk,
		MediaID:      mr.Sk,
	}
	m.MediaSrc = app.MediaSrc{
		Key:    mr.MediaKey,
		Large:  fmt.Sprintf("%s/%s_%s", "thmnb", imgresize.ImgSizeLGPrefix, filepath.Base(mr.MediaKey)),
		Medium: fmt.Sprintf("%s/%s_%s", "thmnb", imgresize.ImgSizeSQMDPrefix, filepath.Base(mr.MediaKey)),
		Small:  fmt.Sprintf("%s/%s_%s", "thmnb", imgresize.ImgSizeSQSMPrefix, filepath.Base(mr.MediaKey)),
	}

	m.MimeType = mr.MimeType
	m.Date = mr.Date
	m.Location = app.Location{
		Region:   mr.LocationRegion,
		Locality: mr.LocationLocality,
		Country: app.Country{
			Short: mr.LocationCountrySh,
			Long:  mr.LocationCountryL,
		},
		Coordinates: app.Coordinates{
			Lat: mr.LocationLat,
			Lng: mr.LocationLng,
		},
	}
	m.Width = mr.Width
	m.Height = mr.Height
	m.Ext = mr.Ext
	m.CameraMake = mr.CameraMake
	m.CameraModel = mr.CameraModel
	m.Hash = mr.Hash
	m.Keywords = mr.Keywords
	m.Title = mr.Title

	return m
}

func newMediaMetaFromRecord(mr mediaDateCollectionRecord) app.MediaMonth {
	dat, _ := time.Parse("2006-01", mr.Date)
	return app.MediaMonth{
		ID:         mr.Date,
		Date:       dat.Format("2006 Jan"),
		MediaCount: mr.MediaCount,
	}
}

func newMediaRecordPK(mediaMeta app.MediaMetadata) string {
	return mediaRecordPrefix + idSeperator + mediaMeta.ID()
}

func newMediaRecordSK(mediaMeta app.MediaMetadata) string {
	return mediaRecordPrefix + "#" + mediaMeta.NewFilename()
}

func newCollectionMediaRecordPK(collectionType, collectionID string) string {
	return fmt.Sprintf(
		"%s%s%s%s%s",
		collectionRecordPrefix,
		idSeperator,
		collectionType,
		idSeperator,
		collectionID,
	)
}

func newCollectionRecordSK(collectionID string) string {
	return fmt.Sprintf(
		"meta%s%s",
		idSeperator,
		collectionID,
	)
}

type collectionMediaRecord struct {
	Pk string `json:"pk"`
	Sk string `json:"sk"`
}

func newCollectionMediaRecord(collectionType, collectionID string, mediaMeta app.MediaMetadata) collectionMediaRecord {
	mdr := collectionMediaRecord{}

	mdr.Pk = newCollectionMediaRecordPK(collectionType, collectionID)
	mdr.Sk = newMediaRecordPK(mediaMeta)

	return mdr
}

func newCollectionRecordKey(collectionID, collectionType string) collectionMediaRecord {
	mdr := collectionMediaRecord{}

	mdr.Pk = newCollectionMediaRecordPK(collectionType, collectionID)
	mdr.Sk = newCollectionRecordSK(collectionID)

	return mdr
}

type collectionRecord struct {
	Gsi1pk     string `json:"gsi1pk,omitempty"`
	Gsi1sk     string `json:"gsi1sk,omitempty"`
	ID         string `json:"collection_id,omitempty"`
	Title      string `json:"collection_title,omitempty"`
	Type       string `json:"collection_type,omitempty"`
	MediaCount int    `json:"media_count,omitempty"`
}

func (r collectionRecord) toCollection() app.Collection {
	return app.Collection{
		ID:         r.ID,
		Title:      r.Title,
		Type:       r.Type,
		MediaCount: r.MediaCount,
	}
}

type collectionRecordUpdate struct {
	Gsi1pk     string `json:":gsi1pk,omitempty"`
	Gsi1sk     string `json:":gsi1sk,omitempty"`
	ID         string `json:":collection_id,omitempty"`
	Title      string `json:":collection_title,omitempty"`
	Type       string `json:":collection_type,omitempty"`
	MediaCount int    `json:":media_count,omitempty"`
}

func newCollectionRecordUpdate(collectionID, collectionType, collectionTitle string) collectionRecordUpdate {
	out := collectionRecordUpdate{}
	out.ID = collectionID
	out.Type = collectionType
	out.Title = collectionTitle
	out.Gsi1pk = collectionRecordPrefix + idSeperator + collectionType
	out.Gsi1sk = "meta" + idSeperator + collectionID
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
		mediaRecord, err := dynamodbattribute.MarshalMap(newMediaRecord(mediaMeta))
		if err != nil {
			return err
		}
		putInput := dynamodb.PutItemInput{
			Item:      mediaRecord,
			TableName: &tableName,
		}
		_, err = client.PutItem(&putInput)
		if err != nil {
			return err
		}

		// collection--timeline_month--2016-04	 media--20160409_223441_23790b8c0f257a9c0789ba144c345a2f
		// media--20160409_223441_23790b8c0f257a9c0789ba144c345a2f
		collectionType := "timeline_month"
		collectionID := mediaMeta.Date.Format("2006-01")
		collectionTitle := mediaMeta.Date.Format("2006 January")

		// // -- save media collection
		collectionMediaRecord, err := dynamodbattribute.MarshalMap(newCollectionMediaRecord(collectionType, collectionID, mediaMeta))
		if err != nil {
			return err
		}
		collectionRecordKey, err := dynamodbattribute.MarshalMap(newCollectionRecordKey(collectionID, collectionType))
		if err != nil {
			return err
		}
		collectionRecordUpdate, err := dynamodbattribute.MarshalMap(newCollectionRecordUpdate(collectionID, collectionType, collectionTitle))
		if err != nil {
			return err
		}

		_, err = client.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName:           aws.String(tableName),
						Item:                collectionMediaRecord,
						ConditionExpression: aws.String("attribute_not_exists(pk)"),
					},
				},
				{
					Update: &dynamodb.Update{
						TableName:                 aws.String(tableName),
						Key:                       collectionRecordKey,
						UpdateExpression:          aws.String("SET collection_id = :collection_id, collection_title = :collection_title, collection_type = :collection_type, gsi1pk = :gsi1pk, gsi1sk = :gsi1sk ADD media_count :media_count"),
						ExpressionAttributeValues: collectionRecordUpdate,
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
			":pk": "collection--timeline_month",
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
			cr := collectionRecord{}
			err = dynamodbattribute.UnmarshalMap(item, &cr)
			if err != nil {
				return timelineView, err
			}

			// -- convert media record to media day
			timelineView.Months = append(timelineView.Months, cr.toCollection())
		}

		return timelineView, nil
	}
}

func NewTimelineMonthQuery(tableName string, client *dynamodb.DynamoDB) app.TimelineMonthQuery {
	return func(monthID string) (app.TimelineMonthView, error) {

		timelineView, err := fetchMediaRecords(client, tableName, monthID)
		if err != nil {
			return timelineView, err
		}
		monthMeta, err := fetchMonthMeta(client, tableName, monthID)
		if err != nil {
			return timelineView, err
		}

		timelineView.CollectionMeta = monthMeta
		return timelineView, err
	}
}

func fetchMonthMeta(client *dynamodb.DynamoDB, tableName, monthID string) (app.MediaMonth, error) {
	meta := app.MediaMonth{}

	eavalues, err := dynamodbattribute.MarshalMap(map[string]string{
		":pk": "month#" + monthID,
		":sk": "META#",
	})
	if err != nil {
		return meta, err
	}
	res, err := client.Query(
		&dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			KeyConditionExpression:    aws.String("pk = :pk and begins_with(sk, :sk)"),
			ExpressionAttributeValues: eavalues,
			ScanIndexForward:          aws.Bool(false),
		},
	)
	if err != nil {
		return meta, err
	}

	for _, item := range res.Items {
		mdr := mediaDateCollectionRecord{}
		err = dynamodbattribute.UnmarshalMap(item, &mdr)
		if err != nil {
			return meta, err
		}

		// -- convert media record to media
		meta := newMediaMetaFromRecord(mdr)
		return meta, nil
	}
	return meta, nil

}

func fetchMediaRecords(client *dynamodb.DynamoDB, tableName, monthID string) (app.TimelineMonthView, error) {
	timelineView := app.TimelineMonthView{}
	eavalues, err := dynamodbattribute.MarshalMap(map[string]string{
		":pk": "month#" + monthID,
		":sk": "media#",
	})
	if err != nil {
		return timelineView, err
	}
	err = client.QueryPages(
		&dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			KeyConditionExpression:    aws.String("pk = :pk and begins_with(sk, :sk)"),
			ExpressionAttributeValues: eavalues,
			ScanIndexForward:          aws.Bool(true),
		},
		func(res *dynamodb.QueryOutput, isLastPg bool) bool {
			for _, item := range res.Items {
				mdr := mediaRecord{}
				err = dynamodbattribute.UnmarshalMap(item, &mdr)
				if err != nil {
					return false
				}

				// -- convert media record to media
				media := newMediaFromMediaRecord(mdr)
				timelineView.Media = append(timelineView.Media, media)
			}

			return true
		})
	return timelineView, err
}

func NewMediaDetailQuery(tableName string, client *dynamodb.DynamoDB) app.MediaDetailQuery {
	return func(mediaID app.MediaCollectionID) (app.MediaDetailView, error) {

		view := app.MediaDetailView{}

		// -- query dynamo
		keyValue, err := dynamodbattribute.MarshalMap(
			map[string]string{
				"pk": mediaID.CollectionID,
				"sk": mediaID.MediaID,
			},
		)
		if err != nil {
			return view, err
		}
		res, err := client.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       keyValue,
		})
		if err != nil {
			return view, err
		}

		mdr := mediaRecord{}
		err = dynamodbattribute.UnmarshalMap(res.Item, &mdr)
		if err != nil {
			return view, err
		}

		// -- convert media record to media day
		media := newMediaFromMediaRecord(mdr)
		view.Media = media

		return view, nil
	}
}

func NewPutLocation(tableName string, client *dynamodb.DynamoDB) app.LocationPutter {
	return func(mediaID app.MediaCollectionID, location app.Location) error {
		keyValue, err := dynamodbattribute.MarshalMap(
			map[string]string{
				"pk": mediaID.CollectionID,
				"sk": mediaID.MediaID,
			},
		)
		if err != nil {
			return err
		}
		updateValues, err := dynamodbattribute.MarshalMap(
			map[string]string{
				":region":    location.Region,
				":locality":  location.Locality,
				":country_s": location.Country.Short,
				":country_l": location.Country.Long,
			},
		)
		if err != nil {
			return err
		}

		_, err = client.UpdateItem(
			&dynamodb.UpdateItemInput{
				TableName:                 aws.String(tableName),
				Key:                       keyValue,
				UpdateExpression:          aws.String("SET location_region=:region, location_locality=:locality, location_country_sh=:country_s, location_country_l=:country_l"),
				ExpressionAttributeValues: updateValues,
			},
		)
		return err
	}
}
