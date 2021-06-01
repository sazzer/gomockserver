package gomockserver

import "net/http"

// Match represents a matching in the mock server to potentially handle incoming requests.
type Match struct {
	rules MatchRules
}

// Matches will check if every rule in this `Match` passes for the incoming request.
func (m *Match) Matches(r *http.Request) bool {
	return m.rules.Matches(r)
}
