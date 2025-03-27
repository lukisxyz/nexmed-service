package custom_middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
func SetConnectionRedis(newConn *redis.Client) error {
	if newConn == nil {
		return errors.New("cannot assign nil connection")
	}
	redisClient = newConn
	return nil
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: auth middleware here

		next.ServeHTTP(w, r)
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