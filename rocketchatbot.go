package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/badkaktus/gorocket"
	"github.com/joho/godotenv"
)

func InitializeChat() {
	client := gorocket.NewClient("https://chat.jankari.tech")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("USER_NAME")
	password := os.Getenv("USER_PASSWORD")
	channel := os.Getenv("ROCKETCHAT_CHANNEL")

	login := gorocket.LoginPayload{
		User:     user,
		Password: password,
	}

	_, err = client.Login(&login)

	if err != nil {
		fmt.Printf("Error: %+v", err)
	}

	str := gorocket.Message{
		Channel: channel,
		Text:    "Hello World",
	}

	_, err = client.PostMessage(&str)
	if err != nil {
		fmt.Printf("Error: %+v", err)
	}

}
