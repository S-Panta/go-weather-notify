package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func GetWeatherData(apiKey string, city string) (WeatherResponse, error) {
	var response WeatherResponse
	baseUrl := "https://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&units=metric"

	url := fmt.Sprintf(baseUrl, apiKey, city)
	res, err := http.Get(url)

	if err != nil {
		return response, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		var httpError HTTPErrorResponse
		err := json.NewDecoder(res.Body).Decode(&httpError)
		if err != nil {
			return response, fmt.Errorf("failed to decode HTTP error response: %w", err)
		}
		return response, fmt.Errorf("error fetching data for %s, %s: %s", city, httpError.Code, httpError.Message)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return response, nil
}
