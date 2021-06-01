package gomockserver_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/matryer/is"
	"github.com/sazzer/gomockserver"
)

func TestStartStopServer(t *testing.T) {
	t.Parallel()

	server := gomockserver.New(t)
	defer server.Close()
}

func TestServerURL(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	url := server.URL()

	match, err := regexp.MatchString(`^http://127\.0\.0\.1:\d+$`, url)
	is.NoErr(err)
	is.True(match)
}

func TestMultipleServers(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server1 := gomockserver.New(t)
	defer server1.Close()

	server2 := gomockserver.New(t)
	defer server2.Close()

	is.True(server1.URL() != server2.URL())
}

func TestNoRoutes(t *testing.T) {
	t.Parallel()

	tests := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}

	for _, tt := range tests { //nolint:paralleltest
		tt := tt

		t.Run(tt, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			server := gomockserver.New(t)
			defer server.Close()

			resp := makeRequest(t, tt, server.URL())
			defer resp.Body.Close()

			is.Equal(resp.StatusCode, http.StatusNotFound)
		})
	}
}

func TestMatchRequestMatches(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchRequest("GET", "/testing/abc"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing/abc", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
}

func TestMatchRequestDoesntMatchURL(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchRequest("GET", "/testing/abc"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/wrong", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusNotFound)
}

func TestMatchRequestDoesntMatchMethod(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchRequest("POST", "/testing/abc"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/wrong", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusNotFound)
}

func makeRequest(t *testing.T, method, url string) *http.Response {
	t.Helper()
	is := is.New(t)

	req, err := http.NewRequestWithContext(context.Background(), method, url, nil)
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	return resp
}
