package gomockserver

// MockServer represents the actual server that will be used in the tests.
type MockServer interface {
	// Close will shut the mock server down. This must always be called, preferably via `defer`.
	Close()
	// URL will generate a URL representing the mock server. This includes the scheme, host and post of the server.
	URL() string
	// Matches will record a new match against the server that will potentially process any incoming requests.
	Matches(...MatchRule) *Match
	// UnmatchedCount will return the number of times a request has been handmed and not matched.
	UnmatchedCount() int
}
