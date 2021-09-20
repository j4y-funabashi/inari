package s3

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewDownloader(bucket, region string) app.Downloader {
	return func(backupFilename string) (string, error) {
		file, err := ioutil.TempFile("", "inari-tmp-")
		if err != nil {
			return "", err
		}
		defer file.Close()
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		downloader := s3manager.NewDownloader(sess)

		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(backupFilename),
			})
		if err != nil {
			return "", err
		}

		return file.Name(), nil
	}
}

func NewUploader(bucket, region string) app.Uploader {
	return func(localFilename, mediaStoreFilename string) error {
		file, err := os.Open(localFilename)
		if err != nil {
			return err
		}
		defer file.Close()
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		uploader := s3manager.NewUploader(sess)
		s3Client := s3.New(sess)

		// check object exists
		headRes, err := s3Client.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(mediaStoreFilename),
		})
		if headRes.ETag != nil {
			return nil
		}

		_, err = uploader.Upload(
			&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(mediaStoreFilename),
				Body:   file,
			})
		if err != nil {
			return err
		}

		return nil
	}
}
