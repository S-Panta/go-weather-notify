package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetWeatherData(apiKey string, city string) (WeatherResponse, error) {

	baseUrl := "https://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&units=metric"

	url := fmt.Sprintf(baseUrl, apiKey, city)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Expected status code be 200 but received:", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var response WeatherResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)

	}
	return response, nil
}
