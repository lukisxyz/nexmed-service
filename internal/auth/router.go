package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/lukisxyz/nexmed-service/internal/middleware"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	// public
	r.Post("/register", registerNewAccountHandler)
	r.Post("/login", loginAccountHandler)
	r.Post("/renew-token", renewTokenHandler)

	// protected
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthenticationMiddleware)
		r.Post("/logout", logoutAccountHandler)
		r.Patch("/change-password", changePasswordAccount)
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