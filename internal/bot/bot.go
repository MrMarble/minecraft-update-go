package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mrmarble/minecraft-update-go/pkg/changelog"
	"github.com/mrmarble/minecraft-update-go/pkg/manifest"
	"github.com/mrmarble/minecraft-update-go/pkg/version"
)

type Bot struct {
	ChannelID string
	LogID     string
	Token     string
}

func (b *Bot) Start() {
	latestManifest, err := manifest.GetLatest()
	if err != nil {
		log.Fatalln("Error getting manifest from server.")
	}
	latestVersion := version.FromManifest(*latestManifest)

	localVersion, err := version.Load()

	if err != nil {
		log.Println("Local version not found.")
		log.Println("Saving latest and exiting")
		latestVersion.Changelog = true
		latestVersion.Save()
		os.Exit(0)
	}

	if latestVersion.ID == localVersion.ID && localVersion.Changelog {
		log.Println("Remote version same as local. Exiting")
		os.Exit(0)
	}

	// New version
	if latestVersion.ID != localVersion.ID {
		log.Printf("New Version: %s", latestVersion.ID)
		if b.LogID != "" {
			b.sendMessage(b.LogID, fmt.Sprintf("New Minecraft version: %s\nChangelog: %s", latestVersion.ID, changelog.URL(latestVersion.ToURL())))
		}
		localVersion = &latestVersion
	}

	// Update Changelog
	if !localVersion.Changelog {
		log.Printf("Fetching changelog for %s", localVersion.ID)
		changelog, err := changelog.FromURL(localVersion.ToURL())
		if err != nil {
			log.Println("Changelog is not published")
			os.Exit(0)
		}
		b.sendMessage(b.ChannelID, changelog.String())
		latestVersion.Changelog = true
	}

	latestVersion.Save()
}

func (b *Bot) sendMessage(chatID, message string) {
	values := map[string]string{
		"chat_id":    chatID,
		"parse_mode": "HTML",
		"text":       message,
	}

	jsonData, err := json.Marshal(values)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.Token), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Fatalln(string(body))
	}
}