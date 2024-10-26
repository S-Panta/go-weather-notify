This project fetch upon openweather api to get weather data and show that in your rocketchat room

### Prerequisites

- Go 1.23
- Hosted RocketChat Server. 

### Steps to Install

1. Clone the repository:
   ```bash
   git clone git@github.com:S-Panta/go-weather-notify.
2. Go to project directory
    ```bash
    cd go-weather-notify
3. Install dependencies
    ```bash
    go mod tidy
4. Get all environment variable
    ```bash
    USER_NAME='rocket-chat-username'
    ROCKETCHAT_SERVER_URL = 'rocket-chat-server-url'
    USER_PASSWORD='rocket-chat-password'
    ROCKETCHAT_CHANNEL='rocket-chat-channel'
    OPEN_WEATHER_API_KEY='open-weather-api-key' (signup to https://openweathermap.org/ to get api key)
    

5. Build the project
    ```bash
    go build
 6. Run the executable file
