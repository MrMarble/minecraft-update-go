package version

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/mrmarble/minecraft-update-go/pkg/manifest"
)

const outputFile = "latest_version.json"

const (
	// Snapshot is a weekly release
	Snapshot = iota
	// Release is a full version release
	Release
	// PreRelease is a pre release
	PreRelease
	// ReleaseCandidate is a release candidate
	ReleaseCandidate
)

var (
	preRelease       = regexp.MustCompile(`^(?:\d+|\.)+-pre\d+$`)
	releaseCandidate = regexp.MustCompile(`^(?:\d+|\.)+-rc\d+$`)
)

var versionTypeFromString = map[string]Type{
	"snapshot": Snapshot,
	"release":  Release,
}

// Type holds iota for version type
type Type uint

// Version holds minecraft version info
type Version struct {
	Type      Type
	ID        string
	Changelog bool
}

// ToURL converts Version ID to URL for changelog
func (v *Version) ToURL() string {
	result := ""

	switch v.Type {
	case Snapshot:
		result = fmt.Sprintf("snapshot-%s", v.ID)
	case Release:
		result = fmt.Sprintf("java-edition-%s", strings.ReplaceAll(v.ID, ".", "-"))
	case PreRelease:
		re := regexp.MustCompile(`^((?:\d+|\.)+)-pre(\d+)$`)
		str := re.ReplaceAllString(v.ID, "$1-pre-release-$2")
		result = strings.ReplaceAll(str, ".", "-")
	case ReleaseCandidate:
		re := regexp.MustCompile(`^((?:\d+|\.)+)-rc(\d+)$`)
		str := re.ReplaceAllString(v.ID, "$1-release-candidate-$2")
		result = strings.ReplaceAll(str, ".", "-")
	}

	return result
}

// FromManifest instanciates a Version from the latest published version on manifest
func FromManifest(m manifest.Manifest) Version {
	return FromManifestVersion(m.Versions[0])
}

func FromManifestVersion(v manifest.Version) Version {
	t := versionTypeFromString[v.Type]
	if t == Snapshot {
		switch {
		case preRelease.MatchString(v.Id):
			t = PreRelease
		case releaseCandidate.MatchString(v.Id):
			t = ReleaseCandidate
		}
	}

	return Version{
		Type:      t,
		ID:        v.Id,
		Changelog: false,
	}
}

// Save writes the Version to a json file
func (v *Version) Save(basePath string) {
	marshaled, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(path.Join(basePath, outputFile), marshaled, fs.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}

// Load a Version from a file
func Load(basePath string) (*Version, error) {
	file, err := ioutil.ReadFile(path.Join(basePath, outputFile))
	if err != nil {
		return nil, err
	}

	var version Version

	err = json.Unmarshal(file, &version)
	if err != nil {
		return nil, err
	}

	return &version, nil
}
