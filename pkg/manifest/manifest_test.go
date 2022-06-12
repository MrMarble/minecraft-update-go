package manifest_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mrmarble/minecraft-update-go/pkg/manifest"
	"github.com/stretchr/testify/require"
)

const jsonManifest = `{
   "latest":{
      "release":"1.19",
      "snapshot":"1.19"
   },
   "versions":[
      {
         "id":"1.19",
         "type":"release",
         "url":"https://launchermeta.mojang.com/v1/packages/0ff6d277d64e547edb774346ad270e36b52cfd2a/1.19.json",
         "time":"2017-01-01T00:00:00Z",
         "releaseTime":"2017-01-01T00:00:00Z"
      }
   ]
}`

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestGetManifest(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(jsonManifest)),
		}
	})

	m, err := manifest.GetLatest(client)
	require.Nil(t, err)
	require.NotNil(t, m)
	require.Equal(t, "1.19", m.Latest.Release)
	require.Equal(t, "1.19", m.Latest.Snapshot)
	require.Equal(t, 1, len(m.Versions))
	require.Equal(t, "1.19", m.Versions[0].Id)
}

func TestManifestUnmarshal(t *testing.T) {
	var m manifest.Manifest
	err := json.Unmarshal([]byte(jsonManifest), &m)
	require.Nil(t, err)
	require.NotNil(t, m)
	require.Equal(t, "1.19", m.Latest.Release)
	require.Equal(t, "1.19", m.Latest.Snapshot)
	require.Equal(t, 1, len(m.Versions))
	require.Equal(t, "1.19", m.Versions[0].Id)
}
