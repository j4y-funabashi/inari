package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"go.uber.org/zap"
)

func main() {

	zlogger, _ := zap.NewDevelopment()
	logger := zlogger.Sugar()
	defer logger.Sync()

	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"

	// -- create client
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	dynamoClient := dynamodb.New(sess)

	mediaID := os.Args[1]

	// deps
	viewDetail := dynamo.NewMediaDetailQuery(mediaStoreTableName, dynamoClient)
	view, err := viewDetail(mediaID)
	if err != nil {
		logger.Errorw("failed to fetch media detail",
			"err", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(view)
}
