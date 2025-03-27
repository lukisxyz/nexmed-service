package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	// public
	r.Post("/register", registerNewAccountHandler)
	r.Post("/login", loginAccount)
	r.Post("/refreh-token", loginAccount)

	// protected
	r.Group(func(r chi.Router) {
		// TODO: apply middleware here
		r.Post("/logout", loginAccount)
	})

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