package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"go.uber.org/zap"
)

func main() {

	mediaID := app.MediaCollectionID{
		CollectionID: "month#2018-05",
		MediaID:      "media#2018/20180527_211329_c436eb8941ec3979e8e9ea74ccea8139.JPG",
	}

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
	location, err := reverseGeocode(mediaID)
	if err != nil {
		logger.Fatal(
			"failed to reverse geocode",
			"mediaID", mediaID,
		)
	}

	logger.Infow(
		"geocode completed",
		"location", location,
	)
}
