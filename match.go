package gomockserver

import "net/http"

// Match represents a matching in the mock server to potentially handle incoming requests.
type Match struct {
	rules     MatchRules
	responses ResponseBuilders
}

// Matches will check if every rule in this `Match` passes for the incoming request.
func (m *Match) Matches(r *http.Request) bool {
	return m.rules.Matches(r)
}

// RespondsWith registers new response builders to use to build the response to an incoming request.
func (m *Match) RespondsWith(builders ...ResponseBuilder) {
	m.responses = append(m.responses, ResponseBuilders(builders))
}
