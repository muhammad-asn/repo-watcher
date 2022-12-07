package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	token   string
	channel string
)

func checkProvider(provider string) (string, error) {
	var err error
	switch provider {
	case "telegram":
		token = os.Getenv("TELEGRAM_TOKEN")
		channel = os.Getenv("TELEGRAM_ID")
		if token == "" {
			log.Fatal("TELEGRAM_TOKEN is missing.")
		}
		if channel == "" {
			log.Fatal("TELEGRAM_ID is missing.")
		}
	default:
		err = fmt.Errorf("provider %s not supported", provider)
	}
	return provider, err
}

func sendNotification(provider string, message string) error {
	var err error
	switch provider {
	case "telegram":
		sendTelegramNotification(message)
	default:
		err = fmt.Errorf("provider %s not supported", provider)
	}
	return err
}

func sendTelegramNotification(message string) {
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	// Create the payload
	data := strings.NewReader("chat_id=" + channel + "&text=" + message)

	// Send the HTTP POST request
	_, err := http.Post(url, "application/x-www-form-urlencoded", data)
	if err != nil {
		panic(err)
	}
}
