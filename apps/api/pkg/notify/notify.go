package notify

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewNotifier(region, topicARN string) app.Notifier {
	return func(mediaMeta app.MediaMetadata) error {
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		svc := sns.New(sess)

		_, err := svc.Publish(&sns.PublishInput{
			Message:  aws.String(mediaMeta.ID()),
			TopicArn: &topicARN,
		})

		return err
	}
}
