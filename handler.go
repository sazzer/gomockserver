package gomockserver

import "net/http"

type handler struct {
	matches []Match
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, match := range h.matches {
		if match.Matches(r) {
			return
		}
	}

	http.NotFound(w, r)
}
