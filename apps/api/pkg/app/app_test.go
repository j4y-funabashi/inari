package app_test

import (
	"testing"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/exiftool"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/j4y_funabashi/inari/apps/api/pkg/notify"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
)

func TestImport(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "it works",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			logger := app.NewNullLogger()
			downloader := storage.NewNullDownloader()
			extract := exiftool.NewNullExtractor()
			uploader := storage.NewNullUploader()
			indexer := index.NewNullIndexer()
			notifier := notify.NewNoopNotifier()

			importMedia := app.NewImporter(logger, downloader, extract, uploader, indexer, notifier)
			importMedia("hellchicken!")

		})
	}
}
