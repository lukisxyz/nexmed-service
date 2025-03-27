package header

import (
	"fmt"
	"net/http"
	"strings"
)

type contextKey string
const (
    UserIDKey   contextKey = "id"
    UserEmail   contextKey = "email"
)

func GetBearerToken(r *http.Request) (string, error) {
    // Get the Authorization header
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return "", fmt.Errorf("no authorization header")
    }

    // Check if it's a Bearer token
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        return "", fmt.Errorf("invalid authorization header format")
    }

    // Return the token part
    return parts[1], nil
}