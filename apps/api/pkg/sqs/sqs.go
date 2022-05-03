package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewNotifier(region string) app.Notifier {
	return func(mediaMeta app.MediaMetadata) error {
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		svc := sqs.New(sess)
		queueURL := "https://sqs.eu-central-1.amazonaws.com/725941804651/funabashi-photos-dev-CreateThumbnailQueue-0fuB2xmTMHuQ"

		_, err := svc.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(mediaMeta.NewFilename()),
			QueueUrl:    aws.String(queueURL),
		})

		return err
	}
}
