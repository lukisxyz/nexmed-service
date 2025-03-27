package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

func registerNewAccountHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return;
	}

	ctx := req.Context()

	email := req.FormValue("email")
	password := req.FormValue("password")

	err := createAccount(ctx, email, password)
	if err != nil {
		if errors.Is(err, ErrFailedCreateAccount) {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w)
}