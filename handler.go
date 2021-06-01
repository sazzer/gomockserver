package gomockserver

import (
	"net/http"
)

type handler struct {
	matches []*Match
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

	http.NotFound(w, r)
}
