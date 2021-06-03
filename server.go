package gomockserver

import (
	"net/http/httptest"
	"testing"
)

type server struct {
	t       *testing.T
	handler *handler
	server  *httptest.Server
}

// New will create a new mock server ready for use in tests.
func New(t *testing.T) MockServer {
	t.Helper()

	handler := handler{
		t:              t,
		unmatchedCount: 0,
	}

	return &server{
		t:       t,
		handler: &handler,
		server:  httptest.NewServer(&handler),
	}
}

func (s *server) Close() {
	if s.server != nil {
		s.server.Close()
		s.server = nil
	}
}

func (s *server) URL() string {
	if s.server == nil {
		s.t.Error("Server has been closed")
	}

	return s.server.URL
}

func (s *server) Matches(rules ...MatchRule) *Match {
	match := &Match{
		rules: rules,
	}

	s.handler.matches = append(s.handler.matches, match)

	return match
}

func (s *server) UnmatchedCount() int {
	return s.handler.unmatchedCount
}
