package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"go.uber.org/zap"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {

	// parse user claims
	userName := req.RequestContext.Authorizer["jwt"].(map[string]interface{})["claims"].(map[string]interface{})["cognito:username"].(string)
	userEmail := req.RequestContext.Authorizer["jwt"].(map[string]interface{})["claims"].(map[string]interface{})["email"].(string)

	monthID := req.PathParameters["monthID"]

	zlogger, _ := zap.NewDevelopment()
	logger := zlogger.Sugar()
	defer logger.Sync()
	logger.
		Infow("timeline month!",
			"username", userName,
			"userEmail", userEmail,
		)

	// ----

	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"

	// -- create client
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	dynamoClient := dynamodb.New(sess)

	// deps
	timelineQuery := dynamo.NewTimelineMonthQuery(mediaStoreTableName, dynamoClient)
	viewTimeline := app.NewTimelineMonthView(timelineQuery)

	timelineView, err := viewTimeline(monthID)
	if err != nil {
		logger.Errorw(
			"failed to fetch timeline",
			"err", err,
		)
		return Response{
			StatusCode:      500,
			IsBase64Encoded: false,
		}, err
	}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(timelineView)
	if err != nil {
		logger.Errorw(
			"failed to encode to json",
			"err", err,
		)
		return Response{
			StatusCode:      500,
			IsBase64Encoded: false,
		}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
