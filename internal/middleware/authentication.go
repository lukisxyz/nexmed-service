package mw

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/lukisxyz/nexmed-service/internal/utils/header"
	jwt "github.com/lukisxyz/nexmed-service/internal/utils/token"
)

var (
	ErrExpired = errors.New("auth: token is expired")
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := header.GetBearerToken(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return;
		}
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return;
		}
		if (time.Now().Unix() > claims.ExpiresAt) {
			writeError(w, http.StatusUnauthorized, ErrExpired)
			return;
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, header.UserIDKey, claims.UserID)
        ctx = context.WithValue(ctx, header.UserEmail, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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