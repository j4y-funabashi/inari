package gpx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/tkrajina/gpxgo/gpx"
)

func NewGpxImporter(addLocationToGPXPoints addLocationToGPXPoints, logger app.Logger) app.Importer {
	return func(inputFilename string) (app.Media, error) {
		m := app.Media{}
		startTime := time.Now()

		ext := strings.ToLower(filepath.Ext(inputFilename))
		if ext != ".gpx" {
			return m, nil
		}

		gpxBytes, err := os.ReadFile(inputFilename)
		if err != nil {
			logger.Error("failed to open file",
				"err", err,
				"filename", inputFilename)
			return m, err
		}
		gpxFile, err := gpx.ParseBytes(gpxBytes)
		if err != nil {
			logger.Error("failed to parse file",
				"err", err,
				"filename", inputFilename)
			return m, err
		}

		allPoints := []app.GPXPoint{}
		for _, track := range gpxFile.Tracks {
			for _, segment := range track.Segments {
				for _, p := range segment.Points {
					allPoints = append(allPoints, app.GPXPoint{
						Timestamp: p.Timestamp,
						Location: app.Location{
							Coordinates: app.Coordinates{
								Lat: p.Latitude,
								Lng: p.Longitude,
							}},
					})
				}
			}
		}

		_, err = addLocationToGPXPoints(allPoints)
		if err != nil {
			logger.Error("failed to add location to gpx points",
				"err", err,
				"filename", inputFilename)
			return m, err
		}
		logger.Info(
			"imported file",
			"time", time.Since(startTime),
			"points", len(allPoints),
			"filename", inputFilename)

		return m, err
	}
}

type addLocationToGPXPoints func(points []app.GPXPoint) ([]app.GPXPoint, error)

func NewAddLocationToGPXPoints(fetchLocation app.LookupTimezone) addLocationToGPXPoints {
	return func(points []app.GPXPoint) ([]app.GPXPoint, error) {
		out := []app.GPXPoint{}

		first := points[0]
		fTimezone, err := fetchLocation(first.Lat, first.Lng, first.Timestamp)
		if err != nil {
			return out, fmt.Errorf("failed to fetch location for first point: %w", err)
		}
		fmt.Printf("first location:: %+v\n", fTimezone)

		last := points[len(points)-1]
		lTimezone, err := fetchLocation(last.Lat, last.Lng, last.Timestamp)
		if err != nil {
			return out, fmt.Errorf("failed to fetch location for last point: %w", err)
		}
		fmt.Printf("last location:: %+v\n", lTimezone)

		// timezones match
		if fTimezone == lTimezone {
			fmt.Println("timezones match")
			for _, point := range points {
				point, err := applyTimezoneToGPXPoint(point, fTimezone)
				if err != nil {
					return out, fmt.Errorf("failed to apply location to gpx point: %w", err)
				}

				out = append(out, point)
			}
			return out, nil
		}

		// timezones differ
		seenHours := map[string]string{}
		fmt.Println("timezones differ")
		for _, point := range points {

			key := point.Timestamp.Format("2006-01-02T15")

			if _, exists := seenHours[key]; !exists {
				pTimezone, err := fetchLocation(point.Lat, point.Lng, point.Timestamp)
				if err != nil {
					return out, fmt.Errorf("failed to fetch location for point: %w", err)
				}
				fmt.Printf("%s :: %s\n", key, pTimezone)
				seenHours[key] = pTimezone
			}

			point, err := applyTimezoneToGPXPoint(point, seenHours[key])
			if err != nil {
				return out, fmt.Errorf("failed to apply location to gpx point: %w", err)
			}

			out = append(out, point)
		}
		return out, nil
	}
}

func applyTimezoneToGPXPoint(point app.GPXPoint, timezone string) (app.GPXPoint, error) {
	// change time to new timezone
	nLocation, err := time.LoadLocation(timezone)
	if err != nil {
		return point, fmt.Errorf("failed to load location: %w", err)
	}
	nt := point.Timestamp.In(nLocation)
	ntUTC := time.Date(nt.Year(), nt.Month(), nt.Day(), nt.Hour(), nt.Minute(), nt.Second(), nt.Nanosecond(), time.UTC)

	point.Timestamp = ntUTC
	point.Location.Timezone = timezone

	return point, nil

}
