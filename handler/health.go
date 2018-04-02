package handler

import (
	"net/http"
)

func HealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{ \"status\": \"alive\", version: \"0.0.1\" }"))
	}
}
