package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/inconshreveable/log15"
	appconfig "github.com/j4y_funabashi/inari/apps/api/pkg/app_config"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
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

func newCollectionDetailHandler(queryCollectionDetail app.CollectionDetailQuery, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		collectionID := ps.ByName("collectionid")
		out, err := queryCollectionDetail(collectionID)
		if err != nil {
			logger.Error("failed to query collection detail",
				"err", err)
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}

}

func newDeleteMediaHandler(deleteMedia app.DeleteMedia, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		mediaID := ps.ByName("mediaid")
		err := deleteMedia(mediaID)
		if err != nil {
			logger.Error("failed to delete media",
				"err", err)
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func newUpdateMediaCaptionHandler(updateMediaCaption app.UpdateMediaTextProperty, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		mediaID := ps.ByName("mediaid")
		newCaption, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to update media caption",
				"err", err)
			panic(err)
		}
		err = updateMediaCaption(mediaID, string(newCaption))
		if err != nil {
			logger.Error("failed to update media caption",
				"err", err)
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func newUpdateMediaHashtagHandler(updateMediaHashtag app.UpdateMediaTextProperty, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		mediaID := ps.ByName("mediaid")
		newHashtag, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to update media hashtag",
				"err", err)
			panic(err)
		}
		err = updateMediaHashtag(mediaID, string(newHashtag))
		if err != nil {
			logger.Error("failed to update media hashtag",
				"err", err)
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	// conf
	baseDir := filepath.Join(os.TempDir(), "inari")

	// deps
	logger := log.New()
	listCollections := appconfig.NewListCollections(baseDir)
	collectionDetail := appconfig.NewCollectionDetail(baseDir)
	deleteMedia := appconfig.NewDeleteMedia(baseDir)
	updateMediaCaption := appconfig.NewUpdateMediaCaption(baseDir)
	updateMediaHashtag := appconfig.NewUpdateMediaHashtag(baseDir)

	// routes
	router := httprouter.New()
	router.GET("/api/timeline/months", newMonthsHandler(listCollections, logger))
	router.GET("/api/timeline/month/:collectionid", newCollectionDetailHandler(collectionDetail, logger))
	router.DELETE("/api/media/:mediaid", newDeleteMediaHandler(deleteMedia, logger))
	router.POST("/api/media/:mediaid/caption", newUpdateMediaCaptionHandler(updateMediaCaption, logger))
	router.POST("/api/media/:mediaid/hashtag", newUpdateMediaHashtagHandler(updateMediaHashtag, logger))

	http.ListenAndServe(":8090", router)
}
