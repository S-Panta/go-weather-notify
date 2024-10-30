IMAGE_NAME = go-weather-notify
TAG = alpha

COMPOSE_FILE = docker-compose.yml

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) .