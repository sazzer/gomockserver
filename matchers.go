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

// MatchRules is a type representing a slice of `MatchRule`.
// This allows for multiple rules to be treated as a single rule.
type MatchRules []MatchRule

func (m MatchRules) Matches(r *http.Request) bool {
	for _, match := range m {
		if !match.Matches(r) {
			return false
		}
	}

	return true
}

// MatchURL builds a `MatchRule` to check if the URL of the request matches the one provided.
// Note that this does a complete match, not a partial one.
func MatchURL(url string) MatchRule {
	return MatchRuleFunc(func(r *http.Request) bool {
		return r.RequestURI == url
	})
}

// MatchMethod builds a `MatchRule` to check if the HTTP Method of the request matches the one provided.
func MatchMethod(method string) MatchRule {
	return MatchRuleFunc(func(r *http.Request) bool {
		return r.Method == method
	})
}

// MatchRequest is a helper that matches both the HTTP Method and URL of the request.
func MatchRequest(method, url string) MatchRule {
	return MatchRules{
		MatchMethod(method),
		MatchURL(url),
	}
}
