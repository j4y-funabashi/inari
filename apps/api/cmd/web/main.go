package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/inconshreveable/log15"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/index"
	"github.com/julienschmidt/httprouter"
)

func newMonthsHandler(listCollections app.CollectionLister, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		out, err := listCollections(app.CollectionTypeInbox)
		if err != nil {
			logger.Error("failed to list collections",
				"err", err)
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}

}

func main() {
	// conf
	baseDir := filepath.Join(os.TempDir(), "inari")
	dbFilepath := filepath.Join(baseDir, "inari-media-db.db")

	// deps
	logger := log.New()
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		logger.Error("failed to open db",
			"err", err)
		panic(err)
	}
	err = index.CreateIndex(db)
	if err != nil {
		logger.Error("failed to create index",
			"err", err)
		panic(err)
	}

	listCollections := index.NewSqliteCollectionLister(db)

	router := httprouter.New()

	router.GET("/api/timeline/months", newMonthsHandler(listCollections, logger))

	http.ListenAndServe(":8090", router)
}
