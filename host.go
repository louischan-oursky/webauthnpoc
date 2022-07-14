package main

import (
	"net/http"
)

func GetHost(r *http.Request) string {
	if host := r.Header.Get("X-Forwarded-Host"); host != "" {
		return host
	}

	if host := r.Header.Get("X-Original-Host"); host != "" {
		return host
	}

	return r.Host
}
