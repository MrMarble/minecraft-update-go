package main

import (
	"flag"
	"log"
	"os"

	"github.com/mrmarble/minecraft-update-go/internal/bot"
)

func main() {

	channel := os.Getenv("MINE_CHANNEL")
	logChannel := os.Getenv("MINE_LOG_CHANNEL")
	token := os.Getenv("MINE_TOKEN")

	flag.StringVar(&channel, "channel", channel, "Telegram notifications channel ID")
	flag.StringVar(&logChannel, "log", logChannel, "Telegram log channel ID")
	flag.StringVar(&token, "token", token, "Telegram bot token")

	flag.Parse()

	if channel == "" {
		log.Println("Telegram channel ID required")
		os.Exit(1)
	}
	if token == "" {
		log.Println("Telegram bot token required")
		os.Exit(1)
	}

	bot := bot.Bot{
		ChannelID: channel,
		LogID:     logChannel,
		Token:     token,
	}

	// Log to file
	f, err := os.OpenFile("output.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	log.SetOutput(f)

	bot.Start()
}
