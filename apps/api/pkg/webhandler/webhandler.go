package webhandler

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	appconfig "github.com/j4y_funabashi/inari/apps/api/pkg/app_config"
	"github.com/j4y_funabashi/inari/apps/api/pkg/storage"
	"github.com/julienschmidt/httprouter"
)

const (
	CollectionsPath = "/api/timeline/months"
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

func NewListCollectionsHandler(listCollections app.CollectionLister, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		collectionType := parseCollectionType(r)

		out, err := listCollections(collectionType)
		if err != nil {
			logger.Error("failed to list collections",
				"err", err)
			panic(err)
		}

		w.Header().Set(ContentType, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}
}

func parseCollectionType(r *http.Request) app.CollectionType {
	collectionType := app.CollectionTypeHashTag

	collectionTypeFilter := r.URL.Query().Get("type")
	if collectionTypeFilter != "" {
		collectionType = app.CollectionType(collectionTypeFilter)
	}

	return collectionType
}

type CreateCollectionRequest struct {
	Title string             `json:"title,omitempty"`
	Type  app.CollectionType `json:"type,omitempty"`
}

type CreateCollectionResponse struct {
	ID    string             `json:"id,omitempty"`
	Title string             `json:"title,omitempty"`
	Type  app.CollectionType `json:"type,omitempty"`
}

func NewCreateCollectionHandler(createCollection app.CreateCollection) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// parse request
		createRequest := CreateCollectionRequest{}
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(requestBody, &createRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newCollection, err := createCollection(createRequest.Title, createRequest.Type)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// response
		createResponse := CreateCollectionResponse{
			ID:    newCollection.ID,
			Title: newCollection.Title,
			Type:  newCollection.Type,
		}
		w.Header().Set(ContentType, ContentTypeJSON)
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func NewCollectionDetailHandler(queryCollectionDetail app.CollectionDetailQuery, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		collectionID := ps.ByName("collectionid")
		out, err := queryCollectionDetail(collectionID)
		if err != nil {
			logger.Error("failed to query collection detail",
				"err", err)
			panic(err)
		}

		w.Header().Set(ContentType, ContentTypeJSON)
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

		w.Header().Set(ContentType, ContentTypeJSON)
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

		w.Header().Set(ContentType, ContentTypeJSON)
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

		w.Header().Set(ContentType, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
	}
}

func newExportMediaHandler(export app.Exporter, logger app.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		mediaID := ps.ByName("mediaid")

		err := export(mediaID)
		if err != nil {
			logger.Error("failed to export media",
				"err", err)
			panic(err)
		}

		w.Header().Set(ContentType, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
	}
}

func NewWebHandler() http.Handler {
	// conf
	baseDir := filepath.Join(os.TempDir(), "inari")

	inariApp := appconfig.New()

	// deps
	logger := slog.Default()
	collectionDetail := appconfig.NewCollectionDetail(baseDir)
	deleteMedia := appconfig.NewDeleteMedia(baseDir)
	exportMedia := appconfig.NewExportMedia(baseDir)
	updateMediaCaption := appconfig.NewUpdateMediaCaption(baseDir)
	updateMediaHashtag := appconfig.NewUpdateMediaHashtag(baseDir)
	queryMediaDetail := appconfig.NewMediaDetail(baseDir)

	// uploader
	micropubBucket := "micropub.funabashi.co.uk"
	mediaBucket := "media.funabashi.co.uk"
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	s3Client := s3.NewFromConfig(cfg)
	s3Uploader := manager.NewUploader(s3Client)

	micropubUploader := storage.NewUploader(micropubBucket, s3Uploader, s3Client)
	mediaUploader := storage.NewUploader(mediaBucket, s3Uploader, s3Client)

	exporter := appconfig.NewExporter(logger, queryMediaDetail, mediaUploader, micropubUploader, baseDir, exportMedia)

	// routes
	router := httprouter.New()

	// collections
	router.POST(CollectionsPath, NewCreateCollectionHandler(inariApp.CreateCollection))
	router.GET(CollectionsPath, NewListCollectionsHandler(inariApp.ListCollections, logger))
	router.GET("/api/timeline/month/:collectionid", NewCollectionDetailHandler(collectionDetail, logger))

	// media
	router.DELETE("/api/media/:mediaid", newDeleteMediaHandler(deleteMedia, logger))
	router.POST("/api/media/:mediaid/caption", newUpdateMediaCaptionHandler(updateMediaCaption, logger))
	router.POST("/api/media/:mediaid/hashtag", newUpdateMediaHashtagHandler(updateMediaHashtag, logger))
	router.POST("/api/media/:mediaid/export", newExportMediaHandler(exporter, logger))

	return router
}
