package manifest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const manifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type Latest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type Version struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Time        time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
}

type Manifest struct {
	Latest   Latest    `json:"latest"`
	Versions []Version `json:"versions"`
}

// GetLatest fetch latest manifest from server
func GetLatest(client *http.Client) (*Manifest, error) {
	resp, err := client.Get(manifestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manifest Manifest

	err = json.Unmarshal(body, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}
