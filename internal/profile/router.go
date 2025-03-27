package profile

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/lukisxyz/nexmed-service/internal/middleware"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Use(mw.AuthenticationMiddleware)
	r.Get("/", getProfileHandler)
	r.Post("/", createProfileHandler)
	r.Put("/", updateProfileHandler)

	return r
}

func writeMessage(
	w http.ResponseWriter,
	status int,
	msg string,
) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(j)
}

func writeError(
	w http.ResponseWriter,
	status int,
	err error,
) {
	writeMessage(w, status, err.Error())
}