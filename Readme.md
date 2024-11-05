This is under development !!!!!
This project fetch upon openweather api to get weather data and show that in your rocketchat room

### Prerequisites

- Go 1.23
- Hosted RocketChat Server. 

### Steps to Install

1. Clone the repository:
   ```bash
   git clone git@github.com:S-Panta/go-weather-notify
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
    OPEN_WEATHER_API_KEY='open-weather-api-key' (signup to https://openweathermap.org/ to get API key)
    

5. Build the project
    ```bash
    go build
6. Run the executable file

### To run using docker

1. rename .env.example and set your own credentials
2. run `docker-compose up` and Enjoy

### To run local Kubernetes setup using miniKube
kubectl is command line tool for using the Kubernetes API and minikube is local Kubernete.
1. Install Kubernetes and minikube locally. Go to https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/ for kubectl installation
   ```
   curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
   sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64
   ```
2. Start your cluster
   ```
   minikube start --driver=docker
   ```
   For more installation queries, go to `https://minikube.sigs.k8s.io/docs/start`

3. Build the docker image of app
    ```
    make docker-build
    ```
4. Apply YAML files
   In deployment.yaml, set your credentials inside `value` of `env`
   ```
   kubectl apply -f deployment.yaml
   kubectl apply -f service.yaml
   ```
5. Access the application
   ```
   minikube service go-weather-notify-service --url
   ```