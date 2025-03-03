package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
	"github.com/joho/godotenv"
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

func GetConfig() (Config, error) {
	err := godotenv.Load()
	// os error will be self explanatory
	// if env not found,exiting the code
	if err != nil {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}
	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	serverURL := os.Getenv("ROCKETCHAT_SERVER_URL")
	username := os.Getenv("USER_EMAIL")
	password := os.Getenv("USER_PASSWORD")
	channel := os.Getenv("ROCKETCHAT_CHANNEL")

	if apiKey == "" || serverURL == "" || username == "" || password == "" || channel == "" {
		return Config{}, fmt.Errorf("missing required environment variables")
	}

	return Config{
		ApiKey:    apiKey,
		ServerUrl: serverURL,
		Username:  username,
		Password:  password,
		Channel:   channel,
	}, nil
}

func main() {
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	u, _ := url.Parse(config.ServerUrl)

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
		log.Fatalf("error getting room id.", err.Error())
	}
	messageChannel := make(chan models.Message)

	if err := client.SubscribeToMessageStream(&models.Channel{Name: config.Channel, ID: roomId}, messageChannel); err != nil {
		log.Fatalf("Failed to subscribe to room: %v", err)
	}
	fmt.Println("Subscribed to message stream for room:", config.Channel, roomId)

	go func() {
		for msg := range messageChannel {
			fmt.Printf("Received message %s\n", msg.Msg)

			if strings.HasPrefix(strings.ToLower(msg.Msg), "weather") {
				words := strings.Fields(msg.Msg)
				cityName := words[1]
				response, err := GetWeatherData(config.ApiKey, cityName)
				if err != nil {
					log.Printf("error fetching weather data for %s: %v\n. Try again", cityName, err)
					continue
				}

				condition := response.Weather[0].Condition
				description := response.Weather[0].Description
				feelsLike := response.Main.FeelsLike
				maxTemp := response.Main.TempMax
				icon := response.Weather[0].Icon
				city := response.CityName

				iconURL := fmt.Sprintf("https://openweathermap.org/img/wn/%s.png", icon)
				msg := fmt.Sprintf(
					"Current Weather Report:\n"+
						"----------------------------\n"+
						"Condition        : %s\n"+
						"Description      : %s\n"+
						"Feels Like       : %.2f°C\n"+
						"Max Temperature  : %.2f°C\n"+
						"City             : %s\n"+
						"Status           : ![Weather Icon](%s)\n"+
						"----------------------------",
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
