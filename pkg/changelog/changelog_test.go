package changelog_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrmarble/minecraft-update-go/pkg/changelog"
	"github.com/stretchr/testify/require"
)

func TestFromUrl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, req.Method, "GET")
		require.Equal(t, req.Header.Get("User-Agent"), "FeedFetcher-Google")
		require.Equal(t, req.Header.Get("Cache-Control"), "no-cache")

		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

	changelog.FromURL(server.URL)
}
