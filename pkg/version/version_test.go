package version_test

import (
	"testing"

	"github.com/mrmarble/minecraft-update-go/pkg/manifest"
	"github.com/mrmarble/minecraft-update-go/pkg/version"
)

var urlTests = []struct {
	name string
	in   version.Version
	out  string
}{
	{"Release", version.Version{Type: version.Release, ID: "1.16.5"}, "java-edition-1-16-5"},
	{"Snapshot", version.Version{Type: version.Snapshot, ID: "21w11a"}, "snapshot-21w11a"},
	{"Release Candidate", version.Version{Type: version.ReleaseCandidate, ID: "1.17-rc2"}, "1-17-release-candidate-2"},
	{"Pre Release", version.Version{Type: version.PreRelease, ID: "1.17-pre1"}, "1-17-pre-release-1"},
}

func TestURL(t *testing.T) {
	for _, tt := range urlTests {
		t.Run(tt.name, func(t *testing.T) {
			url := tt.in.ToURL()
			if url != tt.out {
				t.Errorf("got %q, want %q", url, tt.out)
			}
		})
	}
}

var manifestTest = []struct {
	name string
	in   manifest.Manifest
	out  version.Version
}{
	{"Release", manifest.Manifest{Versions: []manifest.Version{{Type: "release", Id: "1.16.5"}}}, version.Version{Type: version.Release, ID: "1.16.5"}},
	{"Snapshot", manifest.Manifest{Versions: []manifest.Version{{Type: "snapshot", Id: "21w11a"}}}, version.Version{Type: version.Snapshot, ID: "21w11a"}},
	{"Snapshot", manifest.Manifest{Versions: []manifest.Version{{Type: "snapshot", Id: "1.17-rc2"}}}, version.Version{Type: version.ReleaseCandidate, ID: "1.17-rc2"}},
	{"Snapshot", manifest.Manifest{Versions: []manifest.Version{{Type: "snapshot", Id: "1.17-pre1"}}}, version.Version{Type: version.PreRelease, ID: "1.17-pre1"}},
}

func TestFromManifest(t *testing.T) {
	for _, tt := range manifestTest {
		t.Run(tt.name, func(t *testing.T) {
			version := version.FromManifest(tt.in)
			if version != tt.out {
				t.Errorf("got %v, want %v", version, tt.out)
			}
		})
	}
}
