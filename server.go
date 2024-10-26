package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
)

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
	Weather  []WeatherData `json:"weather"`
	Main     Temperature   `json:"main"`
	CityName string        `json:"name"`
}

type Config struct {
	ApiKey    string
	ServerUrl string
	Username  string
	Password  string
	Channel   string
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	serverURL := os.Getenv("ROCKETCHAT_SERVER_URL")
	username := os.Getenv("USER_NAME")
	password := os.Getenv("USER_PASSWORD")
	channel := os.Getenv("ROCKETCHAT_CHANNEL")

	return Config{
		ApiKey:    apiKey,
		ServerUrl: serverURL,
		Username:  username,
		Password:  password,
		Channel:   channel,
	}

}
func GetWeatherData(apiKey string, city string) (WeatherResponse, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

func main() {
	config := GetConfig()

	u, err := url.Parse(config.ServerUrl)
	if err != nil {
		panic(err)
	}

	client, err := realtime.NewClient(u, true)
	if err != nil {
		log.Fatalf("Failed to create real-time client: %v", err)
	}

	if _, err := client.Login(&models.UserCredentials{Email: config.Username, Password: config.Password}); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Println("Logged in successfully!")
	roomId, err := client.GetChannelId(config.Channel)
	if err != nil {
		fmt.Println("error getting room id", err)
	}
	messageChannel := make(chan models.Message)

	if err := client.SubscribeToMessageStream(&models.Channel{Name: config.Channel, ID: roomId}, messageChannel); err != nil {
		log.Fatalf("Failed to subscribe to room: %v", err)
	}
	fmt.Println("Subscribed to message stream for room:", config.Channel, roomId)

	go func() {
		for msg := range messageChannel {
			fmt.Printf("Received message %s\n", msg.Msg)

			if strings.HasPrefix((msg.Msg), "weather") {
				words := strings.Fields(msg.Msg)
				cityName := words[1]
				response, err := GetWeatherData(config.ApiKey, cityName)
				if err != nil {
					fmt.Println("Function call error")
				}
				condition := response.Weather[0].Condition
				description := response.Weather[0].Description
				feelsLike := response.Main.FeelsLike
				maxTemp := response.Main.TempMax
				icon := response.Weather[0].Icon
				city := response.CityName

				iconURL := fmt.Sprintf("https://openweathermap.org/img/wn/%s.png", icon)
				msg := fmt.Sprintf("Weather Update:\nCondition: %s\nDescription: %s\nFeels Like: %.2f°C\nMax Temperature: %.2f°C \n City: %s \n Status: ![Weather Icon](%s)",
					condition, description, feelsLike, maxTemp, city, iconURL)

				reply := &models.Message{
					RoomID: roomId,
					Msg:    msg,
				}

				go func(reply *models.Message) {
					if _, err := client.SendMessage(reply); err != nil {
						log.Printf("Failed to send reply message: %v", err)
					}
				}(reply)
			}

		}

	}()
	select {}
}