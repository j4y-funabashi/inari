package storage

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewDownloader(bucket string, downloader *s3manager.Downloader) app.Downloader {
	return func(backupFilename string) (string, error) {
		tmpFilename := filepath.Join(os.TempDir(), filepath.Base(backupFilename))
		file, err := os.Create(tmpFilename)
		if err != nil {
			return "", err
		}
		defer file.Close()

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

func NewUploader(bucket string, uploader *s3manager.Uploader, s3Client *s3.S3) app.Uploader {
	return func(localFilename, mediaStoreFilename string) error {
		file, err := os.Open(localFilename)
		if err != nil {
			return err
		}
		defer file.Close()

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

func NewLister(bucket, region, keyPrefix string) app.FileLister {

	mediaExt := map[string]bool{
		".jpg": true,
		".mov": true,
		".mp4": true,
		".avi": true,
	}

	return func() ([]string, error) {
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		client := s3.New(sess)

		files := []string{}

		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(keyPrefix),
		}

		err := client.ListObjectsV2Pages(input,
			func(page *s3.ListObjectsV2Output, lastPage bool) bool {
				for _, f := range page.Contents {

					// filter to media ext
					key := *f.Key
					ext := strings.ToLower(filepath.Ext(key))
					if mediaExt[ext] {
						files = append(files, key)
					}
				}
				return true
			})
		if err != nil {
			return files, err
		}

		return files, nil
	}
}
