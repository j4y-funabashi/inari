package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewGeocoder(apiKey, baseURL string) app.Geocoder {
	return func(lat, lng float64) (app.Location, error) {

		// fetch reverse geocode
		geocodeURL := buildURL(lat, lng, baseURL, apiKey)
		res, err := http.Get(geocodeURL)
		if err != nil {
			return app.Location{}, err
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return app.Location{}, err
		}

		// parse results
		results := geocodeRes{}
		err = json.Unmarshal(body, &results)

		if len(results.Results) == 0 {
			return app.Location{}, fmt.Errorf("no geocode results found")
		}

		// create Location
		address := getAddress(results.Results)

		return app.Location{
			Lat:      lat,
			Lng:      lng,
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
	for _, c := range address.AddressComponents {
		if listContains(c.Types, "administrative_area_level_1") {
			return c.LongName
		}
	}
	return ""
}

func getLocality(address geocodeResItem) string {
	for _, c := range address.AddressComponents {
		if listContains(c.Types, "locality") {
			return c.LongName
		}
		if listContains(c.Types, "postal_town") {
			return c.LongName
		}
		if listContains(c.Types, "administrative_area_level_2") {
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
