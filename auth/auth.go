package auth

import (
	"net/http"
)

func loginAccount(
	w http.ResponseWriter,
	req *http.Request,
) {}

func logoutAccount(
	w http.ResponseWriter,
	req *http.Request,
) {}

func renewToken(
	w http.ResponseWriter,
	req *http.Request,
) {}