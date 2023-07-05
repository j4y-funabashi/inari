package index

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	_ "github.com/mattn/go-sqlite3"
)

func CreateIndex(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS media (
			id TEXT NOT NULL PRIMARY KEY,
			date_created DATETIME NOT NULL,
			date_deleted DATETIME,
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

	return nil
}

func NewQueryMediaDetail(db *sql.DB) app.QueryMediaDetail {
	return func(mediaID string) (app.Media, error) {
		return fetchMediaByID(db, mediaID)
	}
}

func fetchMediaByID(db *sql.DB, mediaID string) (app.Media, error) {
	out := app.Media{}

	q := `SELECT
			media_data
			FROM media
			WHERE id = ?;
			`
	row := db.QueryRow(q, mediaID)

	jsonStr := ""
	err := row.Scan(&jsonStr)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal([]byte(jsonStr), &out)
	return out, err

}

func NewDeleteMedia(db *sql.DB) app.DeleteMedia {
	return func(mediaID string) error {

		now := time.Now().Format(time.RFC3339Nano)
		q := `UPDATE media SET date_deleted = ? WHERE id = ?;`
		pq, err := db.Prepare(q)
		if err != nil {
			return err
		}
		_, err = pq.Exec(now, mediaID)
		return err

	}
}

func NewUpdateMediaCaption(db *sql.DB) app.UpdateMediaCaption {
	return func(mediaID, newCaption string) error {
		media, err := fetchMediaByID(db, mediaID)
		if err != nil {
			return err
		}

		media.Caption = newCaption

		return updateMediaDataByID(db, media)
	}
}

func updateMediaDataByID(db *sql.DB, media app.Media) error {
	mediaData, err := json.Marshal(media)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`UPDATE media SET media_data = ? WHERE id = ?;`,
		string(mediaData),
		media.ID)

	return err
}

func NewSqliteIndexer(db *sql.DB) app.Indexer {
	return func(media app.Media) (app.Media, error) {

		media.ID = media.Hash

		// inbox
		media, err := addMediaToCollection(
			db,
			media.Date.Format("2006-01"),
			app.CollectionTypeInbox,
			fmt.Sprintf("inbox %s", media.Date.Format("Jan 2006")),
			media,
		)
		if err != nil {
			return app.Media{}, err
		}

		// camera
		media, err = addMediaToCollection(
			db,
			fmt.Sprintf("%s %s", media.CameraMake, media.CameraModel),
			app.CollectionTypeCamera,
			fmt.Sprintf("%s %s", media.CameraMake, media.CameraModel),
			media,
		)
		if err != nil {
			return app.Media{}, err
		}

		// month
		media, err = addMediaToCollection(
			db,
			media.Date.Format("2006-01"),
			app.CollectionTypeTimelineMonth,
			media.Date.Format("2006 January"),
			media,
		)
		if err != nil {
			return app.Media{}, err
		}

		// day
		media, err = addMediaToCollection(
			db,
			media.Date.Format("2006-01-02"),
			app.CollectionTypeTimelineDay,
			media.Date.Format("Mon, 02 Jan 2006"),
			media,
		)
		if err != nil {
			return app.Media{}, err
		}

		if media.Location.Country.Long != "" {
			// country
			media, err = addMediaToCollection(
				db,
				media.Location.Country.Long,
				app.CollectionTypePlacesCountry,
				media.Location.Country.Long,
				media,
			)
			if err != nil {
				return app.Media{}, err
			}
		}
		if media.Location.Region != "" && media.Location.Country.Long != "" {
			// region
			media, err = addMediaToCollection(
				db,
				fmt.Sprintf("%s, %s", media.Location.Region, media.Location.Country.Long),
				app.CollectionTypePlacesRegion,
				fmt.Sprintf("%s, %s", media.Location.Region, media.Location.Country.Long),
				media,
			)
			if err != nil {
				return app.Media{}, err
			}
		}

		mediaData, err := json.Marshal(media)
		if err != nil {
			return app.Media{}, err
		}
		_, err = db.Exec(
			`INSERT OR IGNORE INTO
			media (id, date_created, media_data)
			VALUES (?,?,?);
			`,
			media.ID,
			media.Date.Format(time.RFC3339),
			string(mediaData))
		if err != nil {
			return app.Media{}, err
		}

		return media, nil
	}
}

func NewSqliteCollectionLister(db *sql.DB) app.CollectionLister {
	return func(collectionType string) ([]app.Collection, error) {
		out := []app.Collection{}

		q := `SELECT
			c.id, c.collection_type, c.title, count(*) as media_count
			FROM collection AS c
			INNER JOIN media_collection ON media_collection.collection_id = c.id
			INNER JOIN media ON media_collection.media_id = media.id
			WHERE c.collection_type = ? AND media.date_deleted IS NULL
			GROUP BY c.id
			ORDER BY c.id DESC;
			`
		rows, err := db.Query(q, collectionType)
		if err != nil {
			return out, err
		}

		for rows.Next() {
			c := app.Collection{}
			err = rows.Scan(&c.ID, &c.Type, &c.Title, &c.MediaCount)
			if err != nil {
				return out, err
			}
			out = append(out, c)
		}

		return out, nil
	}
}

func NewSqliteCollectionDetail(db *sql.DB) app.CollectionDetailQuery {
	return func(collectionID string) (app.CollectionDetail, error) {
		out := app.CollectionDetail{}

		c, err := fetchCollectionByID(db, collectionID)
		if err != nil {
			return out, err
		}
		out.CollectionMeta = c

		media, err := fetchMediaByCollectionID(db, collectionID)
		if err != nil {
			return out, err
		}
		out.Media = media

		return out, nil
	}
}

func fetchMediaByCollectionID(db *sql.DB, collectionID string) ([]app.Media, error) {
	out := []app.Media{}

	q := `SELECT
			media_data
			FROM collection AS c
			INNER JOIN media_collection ON media_collection.collection_id = c.id
			INNER JOIN media ON media_collection.media_id = media.id
			WHERE c.id = ? AND media.date_deleted IS NULL
			ORDER BY media.id DESC;
			`
	rows, err := db.Query(q, collectionID)
	if err != nil {
		return out, err
	}

	for rows.Next() {
		m := app.Media{}
		jsonStr := ""
		err = rows.Scan(&jsonStr)
		if err != nil {
			return out, err
		}
		err = json.Unmarshal([]byte(jsonStr), &m)
		if err != nil {
			return out, err
		}

		m.FormattedDate = m.MediaMetadata.Date.Format(time.RFC3339Nano)

		out = append(out, m)
	}

	return out, nil
}

func fetchCollectionByID(db *sql.DB, collectionID string) (app.Collection, error) {
	q := `SELECT
			c.id, c.collection_type, c.title, count(*) as media_count
			FROM collection AS c
			INNER JOIN media_collection ON media_collection.collection_id = c.id
			INNER JOIN media ON media_collection.media_id = media.id
			WHERE c.id = ? AND media.date_deleted IS NULL
			GROUP BY c.id
			ORDER BY c.id DESC;
			`

	c := app.Collection{}
	err := db.QueryRow(q, collectionID).Scan(&c.ID, &c.Type, &c.Title, &c.MediaCount)
	if err != nil {
		return c, err
	}

	return c, nil
}

func addMediaToCollection(db *sql.DB, collectionID, collectionType, collectionTitle string, media app.Media) (app.Media, error) {

	collectionID = slug.Make(fmt.Sprintf("%s__%s", collectionType, collectionID))

	_, err := db.Exec(
		`INSERT OR IGNORE INTO
		collection (id, collection_type, title)
		VALUES (?,?,?);
		`,
		collectionID,
		collectionType,
		collectionTitle)
	if err != nil {
		return media, err
	}

	_, err = db.Exec(
		`INSERT OR IGNORE INTO
		media_collection (media_id, collection_id)
		VALUES (?,?);
		`,
		media.ID,
		collectionID)

	media.Collections = append(
		media.Collections,
		app.Collection{
			ID:    collectionID,
			Title: collectionTitle,
			Type:  collectionType,
		},
	)

	return media, err
}
