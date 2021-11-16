package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	log "github.com/sirupsen/logrus"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {

	buf := new(bytes.Buffer)

	// parse user claims
	userName := req.RequestContext.Authorizer["jwt"].(map[string]interface{})["claims"].(map[string]interface{})["cognito:username"].(string)
	userEmail := req.RequestContext.Authorizer["jwt"].(map[string]interface{})["claims"].(map[string]interface{})["email"].(string)

	logger := log.StandardLogger()
	logger.
		WithField("username", userName).
		WithField("userEmail", userEmail).
		Info("timeline!")

	// ----

	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"

	// deps
	timelineQuery := dynamo.NewTimelineQuery(mediaStoreTableName, region)
	viewTimeline := app.NewTimelineView(timelineQuery)

	timelineView, err := viewTimeline()
	if err != nil {
		logger.WithError(err).Error("failed to fetch timeline")
		return Response{
			StatusCode:      500,
			IsBase64Encoded: false,
		}, err
	}

	err = json.NewEncoder(buf).Encode(timelineView)
	if err != nil {
		logger.WithError(err).Error("failed to encode to json")
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
