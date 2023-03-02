package notify

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewNotifier(snsClient *sns.SNS, topicARN string) app.Notifier {
	return func(mediaMeta app.Media) error {

		_, err := snsClient.Publish(&sns.PublishInput{
			Message:  aws.String(mediaMeta.ID),
			TopicArn: &topicARN,
		})

		return err
	}
}

func NewNoopNotifier() app.Notifier {
	return func(mediaMeta app.Media) error {
		return nil
	}
}
