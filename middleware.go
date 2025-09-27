package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func CheckAuthHeader(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false
	}
	auth_vals := strings.Split(auth, " ")
	if len(auth_vals) != 2 || auth_vals[0] != "Basic" {
		return false
	}
	data, err := base64.StdEncoding.DecodeString(auth_vals[1])
	if err != nil {
		return false
	}
	username, password := "Chase", "Test"
	correct := fmt.Sprintf("%s:%s", username, password)
	if string(data) != correct {
		return false
	}
	return true
}

func BasicAuth(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if ok := CheckAuthHeader(r); !ok {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"protected\"")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}
