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
