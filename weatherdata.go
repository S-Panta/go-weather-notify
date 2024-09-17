package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func GetWeatherData() {
	type WeatherData struct {
		Condition   string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	type Temperature struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
	}

	type WeatherResponse struct {
		Weather []WeatherData `json:"weather"`
		Main    Temperature   `json:"main"`
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OWM_API_KEY")

	city := "Pokhara"

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
	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON response into the struct
	var response WeatherResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("Weather Condition:", response.Weather[0].Condition)
	fmt.Println("Weather Description:", response.Weather[0].Description)
	fmt.Println("Weather Icon:", response.Weather[0].Icon)

	fmt.Println("Temperature:", response.Main.Temp)
	fmt.Println("Feels Like:", response.Main.FeelsLike)
	fmt.Println("Temperature Min:", response.Main)
	fmt.Println("Temperature Max:", response.Main.TempMax)
}
