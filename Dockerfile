FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go /app/

RUN go build -o /go-weather-notify

EXPOSE 8080

CMD ["/go-weather-notify"]