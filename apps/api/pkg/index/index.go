package index

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	_ "github.com/mattn/go-sqlite3"
)

const (
	collectionTypeInbox         = "inbox"
	collectionTypeTimelineMonth = "timeline_month"
	collectionTypeTimelineDay   = "timeline_day"
	collectionTypePlacesCountry = "places_country"
	collectionTypePlacesRegion  = "places_region"
)

func CreateIndex(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS media (
  id TEXT NOT NULL PRIMARY KEY,
  date_created DATETIME NOT NULL,
  media_data TEXT
  );`
	if _, err := db.Exec(q); err != nil {
		return err
	}

	q = `CREATE TABLE IF NOT EXISTS collection (
			id TEXT NOT NULL PRIMARY KEY,
			collection_type TEXT,
			title TEXT
		);
		CREATE INDEX IF NOT EXISTS
		idx_collection_type ON collection (collection_type);
  `
	if _, err := db.Exec(q); err != nil {
		return err
	}

	q = `CREATE TABLE IF NOT EXISTS
		media_collection (
			media_id TEXT NOT NULL,
			collection_id TEXT NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS
		idx_media_collections ON media_collection (media_id, collection_id);
  `
	if _, err := db.Exec(q); err != nil {
		return err
	}

	fmt.Println("index created")
	return nil
}

func NewSqliteIndexer(db *sql.DB) app.Indexer {
	return func(mediaMeta app.MediaMetadata) error {

		mediaData, err := json.Marshal(mediaMeta)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			`INSERT OR IGNORE INTO
			media (id, date_created, media_data)
			VALUES (?,?,?);
			`,
			mediaMeta.ID(),
			mediaMeta.Date.Format(time.RFC3339),
			string(mediaData))
		if err != nil {
			return err
		}

		// inbox
		err = addMediaToCollection(
			db,
			fmt.Sprintf("inbox_%s", mediaMeta.Date.Format("2006-01")),
			collectionTypeInbox,
			fmt.Sprintf("inbox_%s", mediaMeta.Date.Format("2006-01")),
			mediaMeta.ID(),
		)
		if err != nil {
			return err
		}

		// month
		err = addMediaToCollection(
			db,
			mediaMeta.Date.Format("2006-01"),
			collectionTypeTimelineMonth,
			mediaMeta.Date.Format("2006 January"),
			mediaMeta.ID(),
		)
		if err != nil {
			return err
		}

		// day
		err = addMediaToCollection(
			db,
			mediaMeta.Date.Format("2006-01-02"),
			collectionTypeTimelineDay,
			mediaMeta.Date.Format("Mon, 02 Jan 2006"),
			mediaMeta.ID(),
		)
		if err != nil {
			return err
		}

		return err
	}
}

func addMediaToCollection(db *sql.DB, collectionID, collectionType, collectionTitle, mediaID string) error {

	_, err := db.Exec(
		`INSERT OR IGNORE INTO
		collection (id, collection_type, title)
		VALUES (?,?,?);
		`,
		collectionID,
		collectionType,
		collectionTitle)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`INSERT OR IGNORE INTO
		media_collection (media_id, collection_id)
		VALUES (?,?);
		`,
		mediaID,
		collectionID)

	return err
}
