package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/mrmarble/minecraft-update-go/internal/try"
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

type Context struct {
	manifest      *manifest.Manifest
	latestVersion version.Version
	localVersion  version.Version
}

// Start runs the job once
func (b *Bot) Start(workingDir string) {
	ctx := &Context{}
	fetchManifest(ctx)
	loadLocalVersion(ctx, workingDir)
	compareVersion(ctx, b)
	updateChangelog(ctx, b)

	ctx.latestVersion.Save(workingDir)
}

func fetchManifest(ctx *Context) {
	// Fetch latest manifest from minecraft api
	var latestManifest *manifest.Manifest

	err := try.Do(func(attempt int) (bool, error) {
		var err error
		latestManifest, err = manifest.GetLatest(http.DefaultClient)
		time.Sleep(30 * time.Second)
		return attempt < 3, err
	})
	if err != nil {
		log.Fatal().AnErr("Err", err).Msg("Error getting manifest from server.")
	}

	ctx.manifest = latestManifest
}

func loadLocalVersion(ctx *Context, workingDir string) {
	ctx.latestVersion = version.FromManifest(*ctx.manifest)

	localVersion, err := version.Load(workingDir)
	if err != nil {
		log.Info().Msg("Local version not found.")
		log.Info().Interface("Remote Version", ctx.latestVersion).Msg("Saving remote and exiting.")
		ctx.latestVersion.Changelog = true
		ctx.latestVersion.Save(workingDir)
		os.Exit(0)
	}

	if ctx.latestVersion.ID == localVersion.ID && localVersion.Changelog {
		log.Info().Interface("Version", localVersion).Msg("Remote version same as local. Exiting.")
		os.Exit(0)
	}

	ctx.localVersion = *localVersion
}

func compareVersion(ctx *Context, b *Bot) {
	// New version
	if ctx.latestVersion.ID != ctx.localVersion.ID {
		log.Info().Str("Version", ctx.latestVersion.ID).Msg("New version.")

		if b.LogID != "" {
			b.sendMessage(b.LogID, fmt.Sprintf("New Minecraft version: %s\nChangelog: %s", ctx.latestVersion.ID, changelog.URL(ctx.latestVersion.ToURL())))
		}

		ctx.localVersion = ctx.latestVersion
	}
}

func updateChangelog(ctx *Context, b *Bot) {
	// Update Changelog
	if !ctx.localVersion.Changelog {
		log.Info().Str("Version", ctx.latestVersion.ID).Msg("Fetching changelog.")

		var chlog *changelog.Changelog

		err := try.Do(func(attempt int) (bool, error) {
			var err error
			chlog, err = changelog.FromURL(ctx.latestVersion.ToURL())
			time.Sleep(5 * time.Minute)
			return attempt < 3, err
		})

		if err != nil {
			log.Info().AnErr("Err", err).Msg("Changelog is not published. Exiting")
			os.Exit(0)
		} else {
			log.Info().Str("Title", chlog.Title).Msg("Changelog found.")
			b.sendMessage(b.ChannelID, chlog.String())
			ctx.latestVersion.Changelog = true
		}
	}
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
			log.Warn().AnErr("Err", err).Msg("Error reading HTTP response from Telegram")
			return
		}

		log.Warn().Str("response", string(body)).Msg("Telegram error.")
	}
}
