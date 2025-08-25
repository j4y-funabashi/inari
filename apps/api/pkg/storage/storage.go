package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewNullDownloader() app.Downloader {
	return func(srcFilename string) (string, error) {
		return "test-temp-filename.jpg", nil
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

func NewNullUploader() app.Uploader {
	return func(srcFilename, dstFilename string) error {
		return nil
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

func NewUploader(bucket string, uploader *manager.Uploader, s3Client *s3.Client) app.UploaderB {
	return func(sourceData []byte, mediaStoreFilename string, contentType string) error {
		file := bytes.NewReader(sourceData)

		_, err := uploader.Upload(
			context.Background(),
			&s3.PutObjectInput{
				Bucket:      aws.String(bucket),
				Key:         aws.String(mediaStoreFilename),
				Body:        file,
				ContentType: aws.String(contentType),
			})
		if err != nil {
			return err
		}

		return nil
	}
}
