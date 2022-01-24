package main

import (
	"encoding/json"
	"os"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/dynamo"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.StandardLogger()

	mediaStoreTableName := "inari-dynamodb-dev-InariDatastore-1VAD7YFUNHWKE"
	region := "eu-central-1"

	// deps
	timelineQuery := dynamo.NewTimelineQuery(mediaStoreTableName, region)
	viewTimeline := app.NewTimelineView(timelineQuery)

	timelineView, err := viewTimeline()
	if err != nil {
		logger.WithError(err).Error("failed to fetch timeline")
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(timelineView)
}
