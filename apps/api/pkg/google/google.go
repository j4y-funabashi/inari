package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewNullGeocoder() app.Geocoder {
	return func(lat, lng float64, cTime time.Time) (app.Location, error) {

		if lat == 0 && lng == 0 {
			return app.Location{}, nil
		}

		return app.Location{
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
		}, nil
	}
}

func NewGeocoder(queryNearestGPX app.QueryNearestGPX, logger app.Logger, apiKey, baseURL string) app.Geocoder {
	return func(lat, lng float64, cTime time.Time) (app.Location, error) {

		if lat == 0 && lng == 0 {
			nearestGPX, err := queryNearestGPX(cTime)
			if err != nil {
				return app.Location{}, err
			}
			lat = nearestGPX.Lat
			lng = nearestGPX.Lng
			logger.Info("nearest gpx point found",
				"cTime", cTime,
				"nearestGPXTime", nearestGPX.Timestamp,
				"time diff", cTime.Sub(nearestGPX.Timestamp))
		}

		if lat == 0 && lng == 0 {
			return app.Location{}, nil
		}

		// fetch reverse geocode
		geocodeURL := buildURL(lat, lng, baseURL, apiKey)
		res, err := http.Get(geocodeURL)
		if err != nil {
			return app.Location{}, err
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return app.Location{}, err
		}

		// parse results
		results := geocodeRes{}
		err = json.Unmarshal(body, &results)
		if err != nil {
			return app.Location{}, err
		}

		if len(results.Results) == 0 {
			return app.Location{}, fmt.Errorf("no geocode results found: %s", body)
		}

		// create Location
		address := getAddress(results.Results)
		// fmt.Printf("\n\n%+v\n\n", address)

		return app.Location{
			Coordinates: app.Coordinates{
				Lat: lat,
				Lng: lng,
			},
			Country:  getCountry(address),
			Region:   getRegion(address),
			Locality: getLocality(address),
		}, nil
	}
}

type geocodeRes struct {
	Results []geocodeResItem `json:"results"`
}

type geocodeResItem struct {
	Types             []string    `json:"types"`
	AddressComponents []component `json:"address_components"`
}

type component struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

func buildURL(lat, lng float64, baseURL, apiKey string) string {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	latlng := fmt.Sprintf("%f,%f", lat, lng)
	q.Add("latlng", latlng)
	q.Add("key", apiKey)
	u.RawQuery = q.Encode()
	return u.String()
}

func getAddress(results []geocodeResItem) geocodeResItem {
	for _, item := range results {
		if listContains(item.Types, "street_address") {
			return item
		}
	}

	return results[0]
}

func getCountry(address geocodeResItem) app.Country {
	for _, c := range address.AddressComponents {
		if listContains(c.Types, "country") {
			return app.Country{
				Short: c.ShortName,
				Long:  c.LongName,
			}
		}
	}
	return app.Country{}
}

func getRegion(address geocodeResItem) string {
	componentTypes := []string{
		"administrative_area_level_2",
		"administrative_area_level_1",
	}
	for _, componentType := range componentTypes {
		r := pickAddressComponent(address, componentType)
		if r != "" {
			return r
		}
	}
	return ""
}

func getLocality(address geocodeResItem) string {
	componentTypes := []string{
		"sublocality",
		"locality",
		"postal_town",
		"administrative_area_level_2",
	}
	for _, componentType := range componentTypes {
		r := pickAddressComponent(address, componentType)
		if r != "" {
			return r
		}
	}
	return ""
}

func pickAddressComponent(address geocodeResItem, componentType string) string {
	for _, c := range address.AddressComponents {
		if listContains(c.Types, componentType) {
			return c.LongName
		}
	}
	return ""
}

func listContains(l []string, t string) bool {
	for _, i := range l {
		if i == t {
			return true
		}
	}
	return false
}
