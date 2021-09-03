package gomockserver

// MockServer represents the actual server that will be used in the tests.
type MockServer interface {
	// Close will shut the mock server down. This must always be called, preferably via `defer`.
	Close()
	// URL will generate a URL representing the mock server. This includes the scheme, host and post of the server.
	URL() string
	// Matches will record a new match against the server that will potentially process any incoming requests.
	Matches(...MatchRule) *Match
	// Mount will create a new
	Mount(Mock) *Match
	// UnmatchedCount will return the number of times a request has been handmed and not matched.
	UnmatchedCount() int
}

// Mock represents a lightweight representation of a mock to add to the server.
type Mock struct {
	Matches  []MatchRule
	Response []ResponseBuilder
}

// NewMock will create a new Mock instance from an assortment of MatchRule and ResponseBuilder types.
func NewMock(params ...interface{}) Mock {
	matches := []MatchRule{}
	response := []ResponseBuilder{}

	for _, param := range params {
		if m, ok := param.(MatchRule); ok {
			matches = append(matches, m)
		} else if r, ok := param.(ResponseBuilder); ok {
			response = append(response, r)
		}
	}

	return Mock{
		Matches:  matches,
		Response: response,
	}
}
