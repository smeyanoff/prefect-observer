package application

import "net/http"

func isOriginAllowed(allowedOrigins []string, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
