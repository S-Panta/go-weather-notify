package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	owm "github.com/briandowns/openweathermap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
        log.Fatal("Error loading .env file")
    }
	var apiKey = os.Getenv("OWM_API_KEY")

	w, err := owm.NewCurrent("K", "EN", apiKey)
	if err != nil {
		log.Fatalf("Failed to create OWM client: %v", err)
	}

	err = w.CurrentByName("Phoenix,AZ")
	if err != nil {
		log.Fatalf("Failed to fetch weather data: %v", err)
	}

	fmt.Println(w)
}
