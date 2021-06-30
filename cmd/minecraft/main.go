package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mrmarble/minecraft-update-go/internal/bot"
)

var (
	channel    = os.Getenv("MINE_CHANNEL")
	logChannel = os.Getenv("MINE_LOG_CHANNEL")
	token      = os.Getenv("MINE_TOKEN")
	output     = os.Getenv("MINE_OUTPUT")
	flagAlias  = map[string]string{
		"channel": "c",
		"output":  "o",
		"log":     "l",
		"token":   "t",
	}
)

func init() {

	flag.StringVar(&channel, "channel", channel, "Telegram notifications channel ID")
	flag.StringVar(&logChannel, "log", logChannel, "Telegram log channel ID")
	flag.StringVar(&token, "token", token, "Telegram bot token")
	flag.StringVar(&output, "output", output, "Output directory. Defaults to CWD")

	for from, to := range flagAlias {
		flagSet := flag.Lookup(from)
		flag.Var(flagSet.Value, to, fmt.Sprintf("alias to %s", flagSet.Name))

	}

	flag.Parse()
}

func main() {
	if channel == "" {
		log.Println("Telegram channel ID required.")
		os.Exit(1)
	}
	if token == "" {
		log.Println("Telegram bot token required.")
		os.Exit(1)
	}

	bot := bot.Bot{
		ChannelID: channel,
		LogID:     logChannel,
		Token:     token,
	}

	// Log to file
	f, err := os.OpenFile(path.Join(output, "output.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	log.SetOutput(f)

	bot.Start(output)
}
