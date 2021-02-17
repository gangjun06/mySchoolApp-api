package tools

import (
	"context"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		ctx := r.Context()
		if header != "" {
			ctx = context.WithValue(ctx, "authHeader", header)
		}

		ctx = context.WithValue(ctx, "ip", getRealIP(r))

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-IP")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarder-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
