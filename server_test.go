package gomockserver_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
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

			is.Equal(server.UnmatchedCount(), 1)
		})
	}
}

func TestMatchRequestMatches(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	match := server.Matches(gomockserver.MatchRequest("GET", "/testing/abc"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing/abc", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)

	is.Equal(server.UnmatchedCount(), 0)
	is.Equal(match.Count(), 1)
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

func TestMatchHeader(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchHeader("X-Test", "2"))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(), nil)
	is.NoErr(err)

	req.Header.Add("X-Test", "1")
	req.Header.Add("X-Test", "2")
	req.Header.Add("X-Test", "3")

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
}

func TestMatchJSONBodyFullIdentical(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchJSONFull(`{"a": 1, "b": {"c": 2, "d": "e"}}`))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(),
		bytes.NewReader([]byte(`{"a": 1, "b": {"c": 2, "d": "e"}}`)))
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
}

func TestMatchJSONBodyFullReordered(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchJSONFull(`{"b": {"d": "e", "c": 2}, "a": 1}`))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(),
		bytes.NewReader([]byte(`{"a": 1, "b": {"c": 2, "d": "e"}}`)))
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
}

func TestMatchJSONBodyFullMissingKey(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchJSONFull(`{"b": {"d": "e", "c": 2}}`))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(),
		bytes.NewReader([]byte(`{"a": 1, "b": {"c": 2, "d": "e"}}`)))
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusNotFound)
}

func TestMatchJSONBodyCompatibleIdentical(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchJSONCompatible(`{"a": 1, "b": {"c": 2, "d": "e"}}`))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(),
		bytes.NewReader([]byte(`{"a": 1, "b": {"c": 2, "d": "e"}}`)))
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
}

func TestMatchJSONBodyCompatibleReordered(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchJSONCompatible(`{"b": {"d": "e", "c": 2}, "a": 1}`))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(),
		bytes.NewReader([]byte(`{"a": 1, "b": {"c": 2, "d": "e"}}`)))
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
}

func TestMatchJSONBodyCompatibleMissingKey(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchJSONCompatible(`{"b": {"d": "e", "c": 2}}`))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(),
		bytes.NewReader([]byte(`{"a": 1, "b": {"c": 2, "d": "e"}}`)))
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
}

func TestMatchJSONBodyCompatibleExtraKey(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchJSONCompatible(`{"a": 1, "b": {"d": "e", "c": 2}}`))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL(),
		bytes.NewReader([]byte(`{"b": {"c": 2, "d": "e"}}`)))
	is.NoErr(err)

	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusNotFound)
}

func TestCustomResponseStatus(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchRequest("GET", "/testing/abc")).
		RespondsWith(gomockserver.ResponseStatus(http.StatusAccepted))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing/abc", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusAccepted)
}

func TestCustomResponseHeaders(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchRequest("GET", "/testing/abc")).
		RespondsWith(gomockserver.ResponseSetHeader("content-type", "application/json"),
			gomockserver.ResponseAppendHeader("X-Test", "1"),
			gomockserver.ResponseAppendHeader("X-Test", "2"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing/abc", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
	is.Equal(resp.Header.Get("content-type"), "application/json")
	is.Equal(resp.Header.Values("X-Test"), []string{"1", "2"})
}

func TestCustomResponseBody(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	server.Matches(gomockserver.MatchRequest("GET", "/testing/abc")).
		RespondsWith(gomockserver.ResponseBody([]byte("Hello")))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing/abc", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)

	body, err := ioutil.ReadAll(resp.Body)
	is.NoErr(err)
	is.Equal(body, []byte("Hello"))
}

func TestCustomResponseBodyJSON(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	match := server.Matches(gomockserver.MatchRequest("GET", "/testing/abc")).
		RespondsWith(gomockserver.ResponseJSON("Hello"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing/abc", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
	is.Equal(resp.Header.Get("content-type"), "application/json")

	body, err := ioutil.ReadAll(resp.Body)
	is.NoErr(err)
	is.Equal(body, []byte("\"Hello\""))

	is.Equal(match.Count(), 1)
}

func TestMatchQuery(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	match := server.Matches(gomockserver.MatchRequest("GET", "/testing"), gomockserver.MatchURLQuery("answer", "42"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing?answer=42", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)

	is.Equal(server.UnmatchedCount(), 0)
	is.Equal(match.Count(), 1)
}

func TestMatchQueryWrongValue(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	match := server.Matches(gomockserver.MatchRequest("GET", "/testing"), gomockserver.MatchURLQuery("answer", "42"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing?answer=41", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusNotFound)

	is.Equal(server.UnmatchedCount(), 1)
	is.Equal(match.Count(), 0)
}

func TestMatchQueryRepeated(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	server := gomockserver.New(t)
	defer server.Close()

	match := server.Matches(gomockserver.MatchRequest("GET", "/testing"), gomockserver.MatchURLQuery("answer", "42"))

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/testing?answer=123&answer=42", server.URL()))
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)

	is.Equal(server.UnmatchedCount(), 0)
	is.Equal(match.Count(), 1)
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
