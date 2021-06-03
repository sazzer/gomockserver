package gomockserver

import (
	"fmt"
	"net/http"
	"testing"
)

type handler struct {
	t              *testing.T
	matches        []*Match
	unmatchedCount int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, match := range h.matches {
		if match.Matches(r) {
			response := Response{
				Status:  http.StatusOK,
				Headers: http.Header{},
			}

			match.responses.PopulateResponse(&response, r)

			response.Write(w)

			return
		}
	}

	requestOutput := fmt.Sprintf("%s %s", r.Method, r.RequestURI)

	for name, values := range r.Header {
		for _, value := range values {
			requestOutput = fmt.Sprintf("%s\n%s: %s", requestOutput, name, value)
		}
	}

	h.t.Logf("Unmatched request: %s", requestOutput)

	h.unmatchedCount++

	http.NotFound(w, r)
}
