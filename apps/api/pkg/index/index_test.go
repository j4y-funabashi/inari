package index_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/matryer/is"
)

func TestIndex(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "it saves media to db",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			is := is.New(t)

			// arrange
			dbFilepath := filepath.Join(os.TempDir(), fmt.Sprintf("inari-test-db-%s", uuid.New().String()))
			db, err := sql.Open("sqlite3", dbFilepath)
			if err != nil {
				t.Fatalf("failed to open sqlite db: %s", err)
			}

			indexMedia := index.NewSqliteIndexer(db)
			getMedia := index.NewQueryMediaDetail(db)

			// act
			media := app.MediaMetadata{
				Hash: "test-hash",
				Date: time.Now(),
			}

			err = index.CreateIndex(db)
			if err != nil {
				t.Fatalf("%s", err)
			}

			err = indexMedia(media)
			if err != nil {
				t.Fatalf("%s", err)
			}

			actualMedia, err := getMedia(media.ID())
			if err != nil {
				t.Fatalf("%s", err)
			}

			// assert
			is.Equal(media.ID(), actualMedia.ID())
		})
	}
}
