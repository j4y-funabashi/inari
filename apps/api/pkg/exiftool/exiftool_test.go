package exiftool_test

import (
	"path"
	"testing"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name                string
	backupFilename      string
	expectedMeta        app.MediaMetadata
	expectedNewFilename string
	expectError         bool
}{
	{
		name:           "photo with location",
		backupFilename: "IMG_20220103_134540.jpg",
		expectedMeta: app.MediaMetadata{
			Hash: "9b3f4e51bd961cb321ca234a0b4703f9",
			Location: app.Location{
				Coordinates: app.Coordinates{
					Lat: 53.8700189722222,
					Lng: -1.561703,
				},
			},
			Ext:         "JPG",
			MimeType:    "image/jpeg",
			Width:       "100",
			Height:      "133",
			Date:        time.Date(2022, time.January, 3, 13, 45, 40, 0, time.UTC),
			CameraMake:  "Fairphone",
			CameraModel: "FP3",
		},
	},
	{
		name:           "mov video",
		backupFilename: "P1160866.MOV",
		expectedMeta: app.MediaMetadata{
			Hash:     "1025f263450492c7a27bd44eb3a9d136",
			Location: app.Location{},
			Ext:      "MOV",
			MimeType: "video/quicktime",
			Width:    "640",
			Height:   "480",
			Date:     time.Date(2017, time.March, 20, 21, 16, 36, 0, time.UTC),
		},
	},
	{
		name:           "photo with keywords and caption",
		backupFilename: "p20140321_080118.jpg",
		expectedMeta: app.MediaMetadata{
			Hash:        "caf73e9785fa54300a051df95cfa2db9",
			Location:    app.Location{},
			Ext:         "JPG",
			MimeType:    "image/jpeg",
			Width:       "2448",
			Height:      "3264",
			CameraMake:  "Samsung",
			CameraModel: "GT-I9100",
			Keywords:    "holiday",
			Title:       "Ferry to Rotterdam",
			Date:        time.Date(2014, time.March, 21, 8, 1, 18, 0, time.UTC),
		},
	},
}

var tests2 = []struct {
	name                string
	backupFilename      string
	expectedMeta        app.MediaMetadata
	expectedNewFilename string
	expectError         bool
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
}

func TestImporter(t *testing.T) {

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange
			extractMetadata := exiftool.NewExtractor()

			// act
			fileName := path.Join("./test_data", test.backupFilename)
			t.Log(fileName)
			result, err := extractMetadata(fileName)

			// assert
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedMeta, result)
		})
	}
}
