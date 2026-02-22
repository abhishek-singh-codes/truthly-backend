package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type ExtractLocation struct {
	City    *string
	State   *string
	Country *string
	logger  *slog.Logger
}

func GetNewExtractLocation(logger *slog.Logger) *ExtractLocation {
	return &ExtractLocation{
		logger: logger,
	}
}

/*
Reverse geocoding response structure
*/
type reverseGeoResponse struct {
	Address struct {
		City    *string `json:"city"`
		Town    *string `json:"town"`
		Village *string `json:"village"`
		Suburb  *string `json:"suburb"`
		County  *string `json:"county"`
		State   *string `json:"state"`
		Country *string `json:"country"`
	} `json:"address"`
}

/*
Latitude + Longitude → City, State, Country
*/
func (e *ExtractLocation) FromLatLong(longitude, latitude float32) (*ExtractLocation, error) {
	e.logger.Info(
		"Extracting location from lat/long",
		"longitude", longitude,
		"latitude", latitude,
	)

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?lon=%f&lat=%f&format=json&addressdetails=1",
		longitude,
		latitude,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// REQUIRED by Nominatim policy
	req.Header.Set("User-Agent", "truthly-backend/1.0")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nominatim returned status %d", resp.StatusCode)
	}

	var geoResp reverseGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return nil, err
	}

	addr := geoResp.Address

	// 🔑 City resolution priority
	switch {
	case addr.City != nil:
		e.City = addr.City
	case addr.Town != nil:
		e.City = addr.Town
	case addr.Village != nil:
		e.City = addr.Village
	case addr.Suburb != nil:
		e.City = addr.Suburb
	case addr.County != nil:
		e.City = addr.County
	default:
		e.City = nil
	}

	e.State = addr.State
	e.Country = addr.Country

	e.logger.Info(
		"Location extracted successfully",
		"city", e.City,
		"state", e.State,
		"country", e.Country,
	)

	return e, nil
}
