package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lukisxyz/nexmed-service/internal/utils/header"
	jwt "github.com/lukisxyz/nexmed-service/internal/utils/token"
)

// @Summary Register new account
// @Description Create a new user account with email and password
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email address" example:"user@example.com"
// @Param password formData string true "User password" minLength(8) example:"strongpassword123"
// @Success 201 "Account created successfully"
// @Failure 400 "invalid input data"
// @Failure 500 "Server error while creating account"
// @Router /register [post]
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
	if !isValidEmail(email) {
		writeError(w, http.StatusBadRequest, ErrEmailNotValid)
		return
	}
	password := req.FormValue("password")

	err := createAccount(ctx, email, password)
	if err != nil {
		if errors.Is(err, ErrInternalServer) {
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

// @Summary Login account
// @Description Login account with email and password
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email address" example:"user@example.com"
// @Param password formData string true "User password" minLength(8) example:"strongpassword123"
// @Success 200 {object} Tokens 
// @Failure 400 "invalid input data"
// @Failure 500 "Server error while creating account"
// @Router /login [post]
func loginAccountHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return;
	}

	ctx := req.Context()

	email := req.FormValue("email")
	if !isValidEmail(email) {
		writeError(w, http.StatusBadRequest, ErrEmailNotValid)
		return
	}
	password := req.FormValue("password")

	id, err := loginAccount(ctx, email, password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	token, err := createToken(ctx, id, email)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

// @Summary Logout account
// @Description Logout user account by invalidating refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 "Logged out successfully"
// @Failure 400 "Invalid token format"
// @Failure 401 "Invalid refresh token"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /logout [post]
func logoutAccountHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	refreshToken, err := header.GetBearerToken(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return;
	}
	ctx := req.Context()
    claims, err := jwt.VerifyToken(refreshToken)
    if err != nil {
        writeError(w,  http.StatusUnauthorized, err)
        return
    }
	rdb.Del(ctx, "refresh_token" + claims.UserID)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary Renew access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param refresh_token formData string true "Refresh token" example:"eyJhbGciOiJIUzI1NiIs..."
// @Success 200 {object} Tokens
// @Failure 400 "Invalid input data"
// @Failure 401 "Invalid refresh token"
// @Failure 500 "Server error while creating tokens"
// @Router /renew-token [post]
func renewTokenHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return;
	}

	ctx := req.Context()

	refreshToken := req.FormValue("refresh_token")

    // Verify refresh token
    claims, err := jwt.VerifyToken(refreshToken)
    if err != nil {
        writeError(w, http.StatusUnauthorized, err)
        return
    }

	// check on redis
	if !tokenIsValid(ctx, claims.UserID, refreshToken) {
        writeError(w, http.StatusUnauthorized, errors.New("token not found"))
        return
	}

    // Buat token baru
    token, err := createToken(ctx, claims.UserID, claims.Email)
    if err != nil {
        writeError(w, http.StatusInternalServerError, err)
        return
    }

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

// @Summary Change password
// @Description Change user account password
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Security BearerAuth
// @Param current_password formData string true "Current password" minLength(8) example:"oldpassword123"
// @Param new_password formData string true "New password" minLength(8) example:"newpassword123"
// @Success 200 "Password changed successfully"
// @Failure 400 "Invalid input data"
// @Failure 401 "Unauthorized"
// @Failure 500 "Server error while changing password"
// @Router /change-password [patch]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func changePasswordAccount(
	w http.ResponseWriter,
	req *http.Request,
) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return;
	}

	ctx := req.Context()

	current_password := req.FormValue("current_password")
	new_password := req.FormValue("new_password")

	err := changePassword(ctx, ctx.Value(header.UserEmail).(string), current_password, new_password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w)
}