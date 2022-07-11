package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"go.uber.org/zap"
)

func NewHandler(logger *zap.SugaredLogger, geocodeMedia app.MediaGeocoder) func(ctx context.Context, req events.SNSEvent) error {
	return func(ctx context.Context, req events.SNSEvent) error {
		for _, record := range req.Records {
			mediaID := record.SNS.Message

			location, err := geocodeMedia(mediaID)
			if err != nil {
				logger.
					Errorw("failed to geocode media",
						"err", err,
						"mediaKey", mediaID)
				return err
			}

			logger.
				Infow("geocoded media",
					"mediaKey", mediaID,
					"location", location,
				)

		}

		return nil
	}
}

func main() {

	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"

	// -- create client
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	dynamoClient := dynamodb.New(sess)

	zlogger, _ := zap.NewDevelopment()
	logger := zlogger.Sugar()
	defer logger.Sync()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"

	geocoder := google.NewGeocoder(apiKey, baseURL)
	fetchMedia := dynamo.NewMediaDetailQuery(mediaStoreTableName, dynamoClient)
	saveLocation := dynamo.NewPutLocation(mediaStoreTableName, dynamoClient)
	reverseGeocode := app.NewGeocoder(logger, geocoder, fetchMedia, saveLocation)

	lambda.Start(NewHandler(logger, reverseGeocode))
}
