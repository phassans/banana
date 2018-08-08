package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/phassans/banana/shared"
)

const apiKey = "AIzaSyB3U99ctNSn-x2rvk6QhDhJuqnZR9ko7z4"

// LatLong to hold lat & lon
type LatLong struct {
	Lat float64
	Lon float64
}

// GoogleResponse holds response from google
type GoogleResponse struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

// GetLatLong function is to return lat and lon of a address
func GetLatLong(address string) (LatLong, error) {
	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/geocode/json", nil)
	if err != nil {
		return LatLong{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("address", address)
	q.Add("key", apiKey)
	req.URL.RawQuery = q.Encode()

	logger := shared.GetLogger()

	logger.Info().Msgf("GetLatLong url: %s", req.URL.RawQuery)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return LatLong{}, err
	}

	googleResponse := GoogleResponse{}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return LatLong{}, err
		}
		defer resp.Body.Close()

		if err := json.Unmarshal(bodyBytes, &googleResponse); err != nil {
			panic(err)
		}

		if googleResponse.Status == "OK" {
			return LatLong{Lat: googleResponse.Results[0].Geometry.Location.Lat, Lon: googleResponse.Results[0].Geometry.Location.Lng}, nil
		}

		return LatLong{}, fmt.Errorf("error determining lat&lon for addres: %s. status:%s", address, googleResponse.Status)

	}

	return LatLong{}, fmt.Errorf("error determining lat&lon for addres: %s. status: %d", address, resp.StatusCode)
}
