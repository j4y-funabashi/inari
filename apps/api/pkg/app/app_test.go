package app_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/s3"
)

func TestImporter(t *testing.T) {
	tests := []struct {
		name                string
		backupFilename      string
		expectedMeta        app.MediaMetadata
		expectedNewFilename string
	}{
		{
			name:           "phone photo",
			backupFilename: "jayr/phone/Camera/IMG_20190202_151247.jpg",
			expectedMeta: app.MediaMetadata{
				Hash: "34575d530c16ca1c1e9f656d13678374",
				Location: app.Location{
					Coordinates: app.Coordinates{
						Lat: 53.8303849722222,
						Lng: -1.558461,
					},
				},
				Ext:         "JPG",
				MimeType:    "image/jpeg",
				Width:       "4000",
				Height:      "3000",
				Date:        time.Date(2020, time.April, 25, 20, 24, 12, 0, time.UTC),
				CameraMake:  "Fairphone",
				CameraModel: "FP3",
			},
		},
		{
			name:           "fairphone photo",
			backupFilename: "jayr/phone/Camera/IMG_20200425_202412.jpg",
			expectedMeta: app.MediaMetadata{
				Hash: "8ff065dc26edd00e45fdb8b3f28e3820",
				Location: app.Location{
					Coordinates: app.Coordinates{
						Lat: 53.8303849722222,
						Lng: -1.558461,
					},
				},
				Ext:         "JPG",
				MimeType:    "image/jpeg",
				Width:       "4000",
				Height:      "3000",
				Date:        time.Date(2020, time.April, 25, 20, 24, 12, 0, time.UTC),
				CameraMake:  "Fairphone",
				CameraModel: "FP3",
			},
		},
		{
			name:           "fairphone video",
			backupFilename: "jayr/phone/Camera/VID_20210315_094627.mp4",
			expectedMeta: app.MediaMetadata{
				Hash:     "276daf6038932d5e335a2539ad964dad",
				Location: app.Location{},
				Ext:      "MP4",
				MimeType: "video/mp4",
				Width:    "1920",
				Height:   "1080",
				Date:     time.Date(2021, time.March, 15, 9, 46, 28, 0, time.UTC),
			},
		},
		{
			name:           "GOPRO photo",
			backupFilename: "jayr/camera/gopro/101GOPRO/GOPR1311.JPG",
			expectedMeta: app.MediaMetadata{
				Hash:        "1ce9a7e9b0dac171e142f8a8902e64b8",
				Location:    app.Location{},
				Ext:         "JPG",
				MimeType:    "image/jpeg",
				Width:       "4000",
				Height:      "3000",
				Date:        time.Date(2018, time.September, 30, 16, 18, 25, 0, time.UTC),
				CameraMake:  "GoPro",
				CameraModel: "HERO5 Black",
			},
		},
		{
			name:           "GOPRO video",
			backupFilename: "jayr/camera/gopro/102GOPRO/GOPR2100.MP4",
			expectedMeta: app.MediaMetadata{
				Hash:     "09a9395470166b4244a058679c18206c",
				Location: app.Location{},
				Ext:      "MP4",
				MimeType: "video/mp4",
				Width:    "1920",
				Height:   "1080",
				Date:     time.Date(2021, time.January, 24, 11, 32, 15, 0, time.UTC),
			},
		},
		{
			name:           "lumix photo",
			backupFilename: "jayr/camera/lumix-dmc-lx3/116_PANA/P1160995.JPG",
			expectedMeta: app.MediaMetadata{
				Hash:        "11fc01719df33d1ae5427c2602658c95",
				Location:    app.Location{},
				Ext:         "JPG",
				MimeType:    "image/jpeg",
				Width:       "3648",
				Height:      "2736",
				Date:        time.Date(2017, time.March, 23, 12, 17, 33, 0, time.UTC),
				CameraMake:  "Panasonic",
				CameraModel: "DMC-LX3",
			},
		},
		{
			name:           "lumix video",
			backupFilename: "jayr/camera/lumix-dmc-lx3/116_PANA/P1160084.MOV",
			expectedMeta: app.MediaMetadata{
				Hash:     "3d1d7b638c9f22ba2a0dab84b3d698b6",
				Location: app.Location{},
				Ext:      "MOV",
				MimeType: "video/quicktime",
				Width:    "640",
				Height:   "480",
				Date:     time.Date(2016, time.April, 9, 15, 12, 6, 0, time.UTC),
			},
		},
	}

	for _, test := range tests {
		t.Skip("move these to exiftool pkg")
		t.Run(test.name, func(t *testing.T) {
			// arrange
			bucket := "backup.funabashi"
			region := "eu-central-1"
			mediaStoreBucket := "inari-mediastore-dev"
			mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
			downloader := s3.NewDownloader(bucket, region)
			uploader := s3.NewUploader(mediaStoreBucket, region)
			indexer := dynamo.NewIndexer(mediaStoreTableName, region)
			extractMetadata := exiftool.NewExtractor()
			importMedia := app.NewImporter(downloader, extractMetadata, uploader, indexer)

			// act
			result := importMedia(test.backupFilename)

			// assert
			if diff := cmp.Diff(test.expectedMeta, result); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
