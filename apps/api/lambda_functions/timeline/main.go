package main

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {

	var buf bytes.Buffer

	// parse user claims
	userName := req.RequestContext.Authorizer["jwt"].(map[string]interface{})["claims"].(map[string]interface{})["cognito:username"].(string)
	userEmail := req.RequestContext.Authorizer["jwt"].(map[string]interface{})["claims"].(map[string]interface{})["email"].(string)

	logger := log.StandardLogger()
	logger.
		WithField("username", userName).
		WithField("userEmail", userEmail).
		Info("timeline!")

	now := time.Now()
	body, err := json.Marshal(map[string]interface{}{
		"message": "HELLCHICKEN!" + req.PathParameters["date"] + " :: " + userName + " :: " + userEmail + " :: " + now.String(),
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

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
