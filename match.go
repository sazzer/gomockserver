package gomockserver

import "net/http"

// Match represents a matching in the mock server to potentially handle incoming requests.
type Match struct {
	rules     MatchRules
	responses ResponseBuilders
	count     int
}

// Matches will check if every rule in this `Match` passes for the incoming request.
func (m *Match) Matches(r *http.Request) bool {
	return m.rules.Matches(r)
}

// RespondsWith registers new response builders to use to build the response to an incoming request.
func (m *Match) RespondsWith(builders ...ResponseBuilder) *Match {
	m.responses = append(m.responses, ResponseBuilders(builders))

	return m
}

// Count will return the number of times this match has been used to respond to a request.
func (m *Match) Count() int {
	return m.count
}
