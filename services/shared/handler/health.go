package handler

import "net/http"

// HealthHandler answers with an empty 200. Platforms and Consul only look at
// the status code, so a body would just be one more thing to get wrong.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
