package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrmarble/minecraft-update-go/internal/bot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	channel    = os.Getenv("MINE_CHANNEL")
	logChannel = os.Getenv("MINE_LOG_CHANNEL")
	token      = os.Getenv("MINE_TOKEN")
	output     = os.Getenv("MINE_OUTPUT")
	url        = ""
	flagAlias  = map[string]string{
		"channel": "c",
		"output":  "o",
		"log":     "l",
		"token":   "t",
		"url":     "u",
	}
)

func init() {

	flag.StringVar(&channel, "channel", channel, "Telegram notifications channel ID")
	flag.StringVar(&logChannel, "log", logChannel, "Telegram log channel ID")
	flag.StringVar(&token, "token", token, "Telegram bot token")
	flag.StringVar(&output, "output", output, "Output directory. Defaults to CWD")
	flag.StringVar(&url, "url", url, "Changelog URL to parse. Optional")

	for from, to := range flagAlias {
		flagSet := flag.Lookup(from)
		flag.Var(flagSet.Value, to, fmt.Sprintf("alias to %s", flagSet.Name))
	}

	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	if channel == "" {
		log.Fatal().Msg("Telegram channel ID is required.")
	}
	if token == "" {
		log.Fatal().Msg("Telegram BOT Token is required.")
	}

	bot := bot.Bot{
		ChannelID: channel,
		LogID:     logChannel,
		Token:     token,
	}

	if url == "" {
		bot.Start(output)
	} else {
		bot.Parse(url)
	}

}
