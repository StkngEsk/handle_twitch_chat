/**
 * main.go
 *
 * Copyright (c) 2017 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file)
 */

package main

import (
	"log"
	"os"
	"time"

	"github.com/StkngEsk/handle_twitch_chat/websocket_twitch_connection"
	"github.com/joho/godotenv"
)

func main() {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	myBot := websocket_twitch_connection.TwitchProps{
		Channel:    os.Getenv("CHANNEL_NAME"),
		MsgRate:    time.Duration(20/30) * time.Millisecond,
		Name:       os.Getenv("CHANNEL_NAME"),
		Port:       "6667",
		OAuthToken: os.Getenv("OAUTH_TOKEN"),
		Server:     "irc.chat.twitch.tv",
	}
	myBot.Start()
}
