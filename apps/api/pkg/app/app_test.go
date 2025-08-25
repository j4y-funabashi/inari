package app_test

import (
	"context"
	log "github.com/inconshreveable/log15"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/google/uuid"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	appconfig "github.com/j4y_funabashi/inari/apps/api/pkg/app_config"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
	"gotest.tools/v3/assert"
)

func TestImport(t *testing.T) {
	testCases := []struct {
		desc     string
		filePath string
		expected app.Media
	}{
		{
			desc:     "it imports photo with keywords and caption",
			filePath: "p20140321_080118.jpg",
			expected: app.Media{
				ID: "caf73e9785fa54300a051df95cfa2db9",
				MediaMetadata: app.MediaMetadata{
					Hash:        "caf73e9785fa54300a051df95cfa2db9",
					Coordinates: app.Coordinates{},
					Ext:         "jpg",
					MimeType:    "image/jpeg",
					Width:       "2448",
					Height:      "3264",
					CameraMake:  "Samsung",
					CameraModel: "GT-I9100",
					Keywords:    "holiday",
					Title:       "Ferry to Rotterdam",
					Date:        time.Date(2014, time.March, 21, 8, 1, 18, 0, time.UTC),
				},
				FilePath: "2014/20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
				Thumbnails: app.MediaSrc{
					Small:  "sqsm_20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
					Medium: "sqmd_20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
					Large:  "lg_20140321_080118_caf73e9785fa54300a051df95cfa2db9.jpg",
				},
				Collections: []app.Collection{
					{
						ID:    "inbox__2014-03",
						Title: "inbox Mar 2014",
						Type:  app.CollectionTypeInbox,
					},
					{
						ID:    "camera__samsung-gt-i9100",
						Title: "Samsung GT-I9100",
						Type:  app.CollectionTypeCamera,
					},
					{
						ID:    "timeline_month__2014-03",
						Title: "2014 March",
						Type:  app.CollectionTypeTimelineMonth,
					},
					{
						ID:    "timeline_day__2014-03-21",
						Title: "Fri, 21 Mar 2014",
						Type:  app.CollectionTypeTimelineDay,
					},
				},
				Caption: "Ferry to Rotterdam",
			},
		},
		{
			desc:     "photo with location",
			filePath: "IMG_20220103_134540.jpg",
			expected: app.Media{
				ID: "9b3f4e51bd961cb321ca234a0b4703f9",
				MediaMetadata: app.MediaMetadata{
					Hash: "9b3f4e51bd961cb321ca234a0b4703f9",
					Coordinates: app.Coordinates{
						Lat: 53.8700189722222,
						Lng: -1.561703,
					},
					Ext:         "jpg",
					MimeType:    "image/jpeg",
					Width:       "100",
					Height:      "133",
					Date:        time.Date(2022, time.January, 3, 13, 45, 40, 0, time.UTC),
					CameraMake:  "Fairphone",
					CameraModel: "FP3",
				},
				FilePath: "2022/20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
				Thumbnails: app.MediaSrc{
					Small:  "sqsm_20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
					Medium: "sqmd_20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
					Large:  "lg_20220103_134540_9b3f4e51bd961cb321ca234a0b4703f9.jpg",
				},
				Collections: []app.Collection{
					{
						ID:    "inbox__2022-01",
						Title: "inbox Jan 2022",
						Type:  app.CollectionTypeInbox,
					},
					{
						ID:    "camera__fairphone-fp3",
						Title: "Fairphone FP3",
						Type:  app.CollectionTypeCamera,
					},
					{
						ID:    "timeline_month__2022-01",
						Title: "2022 January",
						Type:  app.CollectionTypeTimelineMonth,
					},
					{
						ID:    "timeline_day__2022-01-03",
						Title: "Mon, 03 Jan 2022",
						Type:  app.CollectionTypeTimelineDay,
					},
					{
						ID:    "places_country__united-kingdom",
						Title: "United Kingdom",
						Type:  app.CollectionTypePlacesCountry,
					},
					{
						ID:    "places_region__west-yorkshire-united-kingdom",
						Title: "West Yorkshire, United Kingdom",
						Type:  app.CollectionTypePlacesRegion,
					},
				},
				Location: app.Location{
					Country: app.Country{
						Long:  "United Kingdom",
						Short: "GB",
					},
					Region:   "West Yorkshire",
					Locality: "Leeds",
					Coordinates: app.Coordinates{
						Lat: 53.8700189722222,
						Lng: -1.561703,
					},
					Timezone: "Europe/London",
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			// create uniq test dir to store db + output files
			testID := "inari-test-" + uuid.New().String()
			testDir := filepath.Join(os.TempDir(), testID)

			// arrange
			importMedia := appconfig.NewMediaImporter(
				testDir,
				appconfig.WithNullLogger(),
				appconfig.WithNullGeocoder(),
			)
			queryMediaDetail := appconfig.NewMediaDetail(testDir)

			// act
			iMedia, err := importMedia(path.Join("./test_data", tC.filePath))
			assert.NilError(t, err)

			actual, err := queryMediaDetail(iMedia.ID)
			assert.NilError(t, err)

			// assert
			assert.DeepEqual(t, tC.expected, actual)

		})
	}
}

func TestExport(t *testing.T) {

	testCases := []struct {
		desc     string
		filePath string
		expected app.Media
	}{
		{
			desc:     "it imports photo with keywords and caption",
			filePath: "p20140321_080118.jpg",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			// create uniq test dir to store db + output files
			testID := "inari-test-" + uuid.New().String()
			testDir := filepath.Join(os.TempDir(), testID)
			logger := log.New()

			// arrange
			importMedia := appconfig.NewMediaImporter(
				testDir,
				appconfig.WithNullLogger(),
				appconfig.WithNullGeocoder(),
			)
			queryMediaDetail := appconfig.NewMediaDetail(testDir)
			exportMedia := appconfig.NewExportMedia(testDir)
			updateHashtag := appconfig.NewUpdateMediaHashtag(testDir)

			// uploader
			micropubBucket := "micropub.funabashi.co.uk"
			mediaBucket := "media.funabashi.co.uk"
			cfg, _ := config.LoadDefaultConfig(context.TODO())
			s3Client := s3.NewFromConfig(cfg)
			s3Uploader := manager.NewUploader(s3Client)

			micropubUploader := storage.NewUploader(micropubBucket, s3Uploader, s3Client)
			mediaUploader := storage.NewUploader(mediaBucket, s3Uploader, s3Client)

			export := appconfig.NewExporter(logger, queryMediaDetail, mediaUploader, micropubUploader, testDir, exportMedia)

			// act
			iMedia, err := importMedia(path.Join("./test_data", tC.filePath))
			assert.NilError(t, err)

			err = updateHashtag(iMedia.ID, "new-test-tag")
			assert.NilError(t, err)

			err = export(iMedia.ID)
			assert.NilError(t, err)

			t.FailNow()
		})
	}

}
