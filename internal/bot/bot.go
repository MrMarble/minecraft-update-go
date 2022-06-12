package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mrmarble/minecraft-update-go/pkg/changelog"
	"github.com/mrmarble/minecraft-update-go/pkg/manifest"
	"github.com/mrmarble/minecraft-update-go/pkg/version"
	"github.com/rs/zerolog/log"
)

// Bot holds telegram information
type Bot struct {
	ChannelID string
	LogID     string
	Token     string
}

// Start runs the job once
func (b *Bot) Start(workingDir string) {
	latestManifest, err := manifest.GetLatest()
	if err != nil {
		log.Fatal().AnErr("Err", err).Msg("Error getting manifest from server.")
	}

	latestVersion := version.FromManifest(*latestManifest)

	localVersion, err := version.Load(workingDir)
	if err != nil {
		log.Info().Msg("Local version not found.")
		log.Info().Interface("Remote Version", latestVersion).Msg("Saving remote and exiting.")
		latestVersion.Changelog = true
		latestVersion.Save(workingDir)
		os.Exit(0)
	}

	if latestVersion.ID == localVersion.ID && localVersion.Changelog {
		log.Info().Interface("Version", localVersion).Msg("Remote version same as local. Exiting.")
		os.Exit(0)
	}

	// New version
	if latestVersion.ID != localVersion.ID {
		log.Info().Str("Version", latestVersion.ID).Msg("New version.")

		if b.LogID != "" {
			b.sendMessage(b.LogID, fmt.Sprintf("New Minecraft version: %s\nChangelog: %s", latestVersion.ID, changelog.URL(latestVersion.ToURL())))
		}

		localVersion = &latestVersion
	}

	// Update Changelog
	if !localVersion.Changelog {
		log.Info().Str("Version", latestVersion.ID).Msg("Fetching changelog.")

		changelog, err := changelog.FromURL(localVersion.ToURL())
		if err != nil {
			log.Info().AnErr("Err", err).Msg("Changelog is not published. Exiting")
			os.Exit(0)
		} else {
			log.Info().Str("Title", changelog.Title).Msg("Changelog found.")
			b.sendMessage(b.ChannelID, changelog.String())
			latestVersion.Changelog = true
		}
	}

	latestVersion.Save(workingDir)
}

// Parse runs the bot against a provided url
func (b *Bot) Parse(url string) {
	log.Info().Str("url", url).Msg("Fetching changelog from url")
	changelog, err := changelog.FromURL(url)

	if err != nil {
		log.Info().AnErr("Err", err).Msg("Changelog could not be reached. Exiting")
		os.Exit(0)
	} else {
		log.Info().Str("Title", changelog.Title).Msg("Changelog found.")
		changelog.URL = url
		b.sendMessage(b.ChannelID, changelog.String())
	}
}

func (b *Bot) sendMessage(chatID, message string) {
	values := map[string]string{
		"chat_id":    chatID,
		"parse_mode": "HTML",
		"text":       message,
	}

	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Fatal().AnErr("Err", err).Msg("Error marshalling message")
	}

	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.Token), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal().AnErr("Err", err).Msg("Could not reach telegram.")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal().AnErr("Err", err).Msg("Error reading HTTP response from Telegram")
		}

		log.Fatal().Str("response", string(body)).Msg("Telegram error.")
	}
}
