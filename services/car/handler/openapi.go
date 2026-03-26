package handler

import (
	"log"
	"net/http"
)

func OpenapiHandler(spec []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(spec)
	}
}
