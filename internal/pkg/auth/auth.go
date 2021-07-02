package auth

import (
	"net/http"
	"strings"
)

// get token from Header
func ExtractToken(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return ""
	}
	splitted := strings.Split(tokenHeader, " ")

	if len(splitted) != 2 {
		return ""
	}
	tokenpath := splitted[1]

	return tokenpath
}
