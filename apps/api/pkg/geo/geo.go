package geo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

func NewTZAPILookupTimezone(baseURL string) app.LookupTimezone {
	return func(lat, lng float64, cTime time.Time) (string, error) {

		type TzAPIResponse struct {
			TZ string `json:"tz"`
		}

		apiURL := fmt.Sprintf("%s/tz/%s/%s", baseURL, strconv.FormatFloat(lat, 'f', -1, 64), strconv.FormatFloat(lng, 'f', -1, 64))

		res, err := http.Get(apiURL)
		if err != nil {
			return "", fmt.Errorf("failed to get url %s", err.Error())
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read body %s", err.Error())
		}

		resJSON := TzAPIResponse{}

		err = json.Unmarshal(body, &resJSON)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal json %s", err.Error())
		}

		return resJSON.TZ, nil
	}
}
