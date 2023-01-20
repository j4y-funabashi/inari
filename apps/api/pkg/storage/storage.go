package storage

import (
	"fmt"
	"io"
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

func NewLocalFSDownloader() app.Downloader {
	return func(srcFilename string) (string, error) {
		dstFilename := filepath.Join(os.TempDir(), filepath.Base(srcFilename))
		dstFile, err := os.Create(dstFilename)
		if err != nil {
			return "", err
		}
		defer dstFile.Close()

		srcFile, err := os.Open(srcFilename)
		if err != nil {
			return "", err
		}
		defer srcFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return "", err
		}

		return dstFile.Name(), nil
	}
}

func NewLocalFSUploader(dstRootDir string) app.Uploader {
	return func(srcFilename, dstFilename string) error {

		dstFilename = filepath.Join(dstRootDir, dstFilename)

		if _, err := os.Stat(dstFilename); !os.IsNotExist(err) {
			return nil
		}

		err := os.MkdirAll(filepath.Dir(dstFilename), os.ModePerm)
		if err != nil {
			return err
		}

		srcFile, err := os.Open(srcFilename)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstFilename)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}

		return nil
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
		seen := map[string]bool{}
		totalSize := 0

		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(keyPrefix),
		}

		err := client.ListObjectsV2Pages(input,
			func(page *s3.ListObjectsV2Output, lastPage bool) bool {
				for _, f := range page.Contents {

					// filter to media ext
					key := *f.Key
					if _, ok := seen[*f.ETag]; ok {
						continue
					}
					ext := strings.ToLower(filepath.Ext(key))
					if mediaExt[ext] {
						files = append(files, key)
						seen[*f.ETag] = true
						totalSize += int(*f.Size)
					}
				}
				return true
			})
		if err != nil {
			return files, err
		}

		fmt.Println("totalSize:: ", totalSize)

		return files, nil
	}
}
