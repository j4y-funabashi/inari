package exiftool

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	exiftoolz "github.com/barasher/go-exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewExtractor(exiftoolPath string) app.MetadataExtractor {
	return func(mediaFile string) (app.MediaMetadata, error) {
		mediaMetadata := app.MediaMetadata{}
		//et, err := exiftoolz.NewExiftool(exiftoolz.NoPrintConversion(), exiftoolz.SetExiftoolBinaryPath(exiftoolPath))
		et, err := exiftoolz.NewExiftool(exiftoolz.NoPrintConversion())
		if err != nil {
			return mediaMetadata, err
		}
		defer et.Close()
		fileInfos := et.ExtractMetadata(mediaFile)
		fileInfo := fileInfos[0]
		if fileInfo.Err != nil {
			return mediaMetadata, fileInfo.Err
		}

		date, err := parseDate(fileInfo)
		if err != nil {
			return mediaMetadata, err
		}
		// listKeys(fileInfo)
		coordinates := parseGPS(fileInfo)
		ext := parseExt(fileInfo)
		mimeType := parseMimeType(fileInfo)
		cameraModel := parseCameraModel(fileInfo)
		cameraMake := parseCameraMake(fileInfo)
		width := parseImageWidth(fileInfo)
		height := parseImageHeight(fileInfo)
		keywords := parseKeywords(fileInfo)
		title := parseTitle(fileInfo)

		hash, err := parseHash(mediaFile)
		if err != nil {
			return mediaMetadata, err
		}

		mediaMetadata.Location.Coordinates = coordinates
		mediaMetadata.Date = date
		mediaMetadata.Hash = hash
		mediaMetadata.Ext = ext
		mediaMetadata.MimeType = mimeType
		mediaMetadata.CameraModel = cameraModel
		mediaMetadata.CameraMake = cameraMake
		mediaMetadata.Width = width
		mediaMetadata.Height = height
		mediaMetadata.Keywords = keywords
		mediaMetadata.Title = title

		return mediaMetadata, nil
	}
}

func listKeys(fileInfo exiftoolz.FileMetadata) {
	for k, v := range fileInfo.Fields {
		fmt.Printf("%s: %v", k, v)
	}
}

func parseHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func parseDate(fileInfo exiftoolz.FileMetadata) (time.Time, error) {
	// DateTimeOriginal -> 2019:02:02 15:12:47
	datString := getDateString(fileInfo)
	if datString == "" {
		return time.Now(), errors.New("file does not contain valid date key")
	}

	dat, err := time.Parse("2006:01:02 15:04:05", datString)
	if err != nil {
		return time.Now(), err
	}
	return dat, nil
}

func getDateString(fileInfo exiftoolz.FileMetadata) string {
	dateKeys := []string{"DateTimeOriginal", "CreateDate"}

	for _, dateKey := range dateKeys {
		datString, err := fileInfo.GetString(dateKey)
		if err != nil {
			continue
		}
		return datString
	}

	return ""
}

func parseGPS(fileInfo exiftoolz.FileMetadata) app.Coordinates {
	coordinates := app.Coordinates{}

	latVal, err := fileInfo.GetFloat("GPSLatitude")
	if err != nil {
		return coordinates
	}
	lngVal, err := fileInfo.GetFloat("GPSLongitude")
	if err != nil {
		return coordinates
	}
	coordinates.Lat = latVal
	coordinates.Lng = lngVal
	return coordinates
}

func parseExt(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("FileTypeExtension")
	if err != nil {
		return ""
	}
	return extVal
}

func parseKeywords(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("Keywords")
	if err != nil {
		return ""
	}
	return extVal
}

func parseTitle(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("Title")
	if err != nil {
		return ""
	}
	return extVal
}

func parseMimeType(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("MIMEType")
	if err != nil {
		return ""
	}
	return extVal
}

func parseCameraModel(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("Model")
	if err != nil {
		return ""
	}
	return extVal
}

func parseCameraMake(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("Make")
	if err != nil {
		return ""
	}
	return extVal
}

func parseImageWidth(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("ImageWidth")
	if err != nil {
		return ""
	}
	return extVal
}

func parseImageHeight(fileInfo exiftoolz.FileMetadata) string {
	extVal, err := fileInfo.GetString("ImageHeight")
	if err != nil {
		return ""
	}
	return extVal
}
