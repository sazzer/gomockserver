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

func (s *server) Matches(rules ...MatchRule) Match {
	match := Match{
		rules: rules,
	}

	s.handler.matches = append(s.handler.matches, match)

	return match
}