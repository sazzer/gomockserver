package gomockserver

import (
	"net/http"
	"net/url"
)

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

func matchURL(matcher func(url.URL) bool) MatchRule {
	return MatchRuleFunc(func(r *http.Request) bool {
		uri, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			return false
		}

		if uri == nil {
			return false
		}

		return matcher(*uri)
	})
}

// MatchURLPath builds a `MatchRule` to check if the URL Path of the request matches the one provided.
// Note that this does a complete match, not a partial one.
func MatchURLPath(expected string) MatchRule {
	return matchURL(func(uri url.URL) bool {
		return uri.EscapedPath() == expected
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
		MatchURLPath(url),
	}
}

// MatchHeader builds a `MatchRule` to check if the given header is present and has the given value.
// If the header is repeated then only one of the repeated values needs to have the provided value.
func MatchHeader(name, value string) MatchRule {
	return MatchRuleFunc(func(r *http.Request) bool {
		values := r.Header.Values(name)

		for _, v := range values {
			if v == value {
				return true
			}
		}

		return false
	})
}

// MatchURLQuery builds a `MatchRule` to check if a query parameter is present with the given name and value.
// If the query parameter is repeated then only one of the repeated values needs to have the provided value.
func MatchURLQuery(name, value string) MatchRule {
	return matchURL(func(uri url.URL) bool {
		for n, values := range uri.Query() {
			if name == n {
				for _, v := range values {
					if value == v {
						return true
					}
				}
			}
		}

		return false
	})
}
