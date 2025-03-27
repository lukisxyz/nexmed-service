package jwt

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Claims struct {
    UserID    string    `json:"userId"`
    Email     string    `json:"email"`
    ExpiresAt int64     `json:"exp"`
    Nonce string `json:"nonce"`
}

var secret string

func SetJwtConfig(_secret string) {
    secret = _secret
}

func CreateToken(ctx context.Context, claims Claims) (string, error) {
    // Generate a unique nonce for each token
    nonceBytes := make([]byte, 16)
    _, err := rand.Read(nonceBytes)
    if err != nil {
        return "", fmt.Errorf("failed to generate nonce: %v", err)
    }
    claims.Nonce = base64.RawURLEncoding.EncodeToString(nonceBytes)

    // Set a precise expiration time
    claims.ExpiresAt = time.Now().Add(time.Minute * 30).Unix() // 30 minutes from now

    // Header
    header := map[string]string{
        "alg": "HS256",
        "typ": "JWT",
    }
    headerBytes, _ := json.Marshal(header)
    headerBase64 := base64.RawURLEncoding.EncodeToString(headerBytes)

    // Payload
    payloadBytes, _ := json.Marshal(claims)
    payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

    // Signature
    signatureInput := headerBase64 + "." + payloadBase64
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(signatureInput))
    signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

    return headerBase64 + "." + payloadBase64 + "." + signature, nil
}

func VerifyToken(tokenString string) (Claims, error) {
    parts := strings.Split(tokenString, ".")
    if len(parts) != 3 {
        return Claims{}, fmt.Errorf("invalid token format")
    }

    // Verify signature
    signatureInput := parts[0] + "." + parts[1]
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(signatureInput))
    expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

    if parts[2] != expectedSignature {
        return Claims{}, fmt.Errorf("invalid signature")
    }

    // Decode payload
    payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return Claims{}, err
    }

    var claims Claims
    if err := json.Unmarshal(payloadBytes, &claims); err != nil {
        return Claims{}, err
    }

    // Check expiration
    if claims.ExpiresAt < time.Now().Unix() {
        return Claims{}, fmt.Errorf("token expired")
    }

    return claims, nil
}
