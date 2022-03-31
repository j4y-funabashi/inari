package main

import (
	"os"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/j4y_funabashi/inari/apps/api/pkg/google"
	"go.uber.org/zap"
)

func main() {

	lat := 53.8700189722222
	lng := -1.561703
	lat = 40.416775
	lng = -3.703790
	lat = 51.173
	lng = -1.776
	lat = 51.179
	lng = -1.826

	zlogger, _ := zap.NewDevelopment()
	logger := zlogger.Sugar()
	defer logger.Sync()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"

	geocoder := google.NewGeocoder(apiKey, baseURL)
	reverseGeocode := app.NewGeocoder(logger, geocoder)
	location, err := reverseGeocode(lat, lng)
	if err != nil {
		logger.Fatal("OH NOES")
	}

	logger.Infow(
		"geocode completed",
		"location", location,
	)
}
