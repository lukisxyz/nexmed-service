package auth

import (
	"context"
	"time"

	jwt "github.com/lukisxyz/nexmed-service/internal/utils/token"
)

type Tokens struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

func tokenIsValid(ctx context.Context, userID, tokenFromRequest string) bool {
    res, err := rdb.Get(ctx, "refresh_token:" + userID).Result()
    if err != nil {
        return false
    }
    return tokenFromRequest == res
}

func createToken(ctx context.Context, userID string, email string) (Tokens, error) {
    // Access token - expires in 5 minute
    accessClaims := jwt.Claims{
        UserID: userID,
		Email: email,
        ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
    }
    accessToken, err := jwt.CreateToken(ctx, accessClaims)
    if err != nil {
        return Tokens{}, err
    }

    // Refresh token - expires in 7 days
    refreshClaims := jwt.Claims{
        UserID: userID,
		Email: email,
        ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
    }
    refreshToken, err := jwt.CreateToken(ctx, refreshClaims)
    if err != nil {
        return Tokens{}, err
    }

    refreshKey := "refresh_token:" + userID
    rdb.Del(ctx, refreshKey)
    err = rdb.Set(ctx, refreshKey, refreshToken, time.Hour*24*7).Err()
    if err != nil {
        return Tokens{}, err
    }

    return Tokens{
        AccessToken: accessToken,
        RefreshToken: refreshToken,
    }, nil
}