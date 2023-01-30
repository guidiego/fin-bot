package api

import (
	"fmt"
	"net/http"

	"github.com/guidiego/fin-bot/config"
)

func IsTokenOk(w http.ResponseWriter, r *http.Request) bool {
	token, isTokenOkey := r.Header["Authorization"]

	if isTokenOkey && token[0] == config.Application.ApiToken {
		return true
	}

	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, "Not Allowed")
	return false
}
