package gomockserver

import "net/http"

// MatchRule represents a rule to match against to see if a request should be processed.
type MatchRule interface {
	// Matches will check to see if the provided HTTP request matches this rule.
	Matches(r *http.Request) bool
}

// MatchRuleFunc is a function type that implements the `MatchRule` interface.
// This allows for simple functions to be used in place of the interface.
type MatchRuleFunc func(*http.Request) bool

func (m MatchRuleFunc) Matches(r *http.Request) bool {
	return m(r)
}

// MatchURL builds a `MatchRule` to check if the URL of the request matches the one provided.
// Note that this does a complete match, not a partial one.
func MatchURL(url string) MatchRule {
	return MatchRuleFunc(func(r *http.Request) bool {
		return r.RequestURI == url
	})
}
