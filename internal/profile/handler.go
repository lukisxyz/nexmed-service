package profile

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/lukisxyz/nexmed-service/internal/utils/header"
	"gopkg.in/guregu/null.v4"
)

// @Summary Create a new user profile
// @Description Create a profile for the authenticated user
// @Tags profile
// @Accept x-www-form-urlencoded
// @Produce json
// @Security BearerAuth
// @Param full_name formData string true "Full Name"
// @Param bio formData string false "Short biography"
// @Param phone_number formData string true "Phone Number"
// @Param address formData string true "Address"
// @Param birth_date formData string true "Birth Date" format(date-time) example:"2000-01-01T00:00:00Z"
// @Success 201 {object} map[string][]string "Profile created successfully"
// @Failure 400 "Invalid input data"
// @Failure 422 "Unable to process the request"
// @Router /profile [post]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func createProfileHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return;
	}

	ctx := req.Context()

	id := ctx.Value(header.UserIDKey).(string)
	fullName := req.FormValue("full_name")
	bio := req.FormValue("bio")
	phoneNumber := req.FormValue("phone_number")
	address := req.FormValue("address")
	birthDateRaw := req.FormValue("birth_date")
	birthDateTime, _ := time.Parse("2006-01-02", birthDateRaw)
	birthDate := null.TimeFrom(birthDateTime)

	profile, err := createProfile(
		ctx,
		id,
		fullName,
		address,
		phoneNumber,
		birthDate,
		bio,
	)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return;
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags profile
// @Accept x-www-form-urlencoded
// @Produce json
// @Security BearerAuth
// @Param full_name formData string false "Full Name"
// @Param bio formData string false "Short biography"
// @Param phone_number formData string false "Phone Number"
// @Param address formData string false "Address"
// @Param birth_date formData string false "Birth Date" format(date-time) example:"2000-01-01T00:00:00Z"
// @Success 200 {object} map[string][]string "Profile updated successfully"
// @Failure 400 "Invalid input data"
// @Failure 422 "Unable to process the request"
// @Router /profile [put]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func updateProfileHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	if err := req.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return;
	}

	ctx := req.Context()

	id := ctx.Value(header.UserIDKey).(string)
	fullName := req.FormValue("full_name")
	bio := req.FormValue("bio")
	phoneNumber := req.FormValue("phone_number")
	address := req.FormValue("address")
	birthDateRaw := req.FormValue("birth_date")
	birthDateTime, _ := time.Parse("2006-01-02", birthDateRaw)
	birthDate := null.TimeFrom(birthDateTime)

	profile, err := updateProfile(
		ctx,
		id,
		fullName,
		address,
		phoneNumber,
		birthDate,
		bio,
	)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return;
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// @Summary Get user profile
// @Description Retrieve the profile details of the authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]string "Profile retrieved successfully"
// @Failure 404 "Profile not found"
// @Failure 400 "Invalid request"
// @Router /profile [get]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func getProfileHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	ctx := req.Context()
	id := ctx.Value(header.UserIDKey).(string)
	
	profile, err := getProfile(ctx, id)
	if err != nil {
		if (errors.Is(err, ErrProfileNotFound)) {
			writeError(w, http.StatusNotFound, err)
			return;
		}
		writeError(w, http.StatusBadRequest, err)
		return;
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}