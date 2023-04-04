package middlewares

import (
	"context"
	"net/http"
)

// AuthTokenReader decodes the share session cookie and packs the session into context
func AuthTokenReader() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			cookie, _ := r.Cookie("jwt")

			if cookie == nil {
				tokenString = r.Header.Get("Authorization")
			} else {
				tokenString = cookie.Value
			}

			// put it in context and call the next with our new context
			ctx := context.WithValue(r.Context(), authTokenCtxKey, tokenString)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// OrgUIDReader decodes the share session cookie and packs the session into context
func OrgUIDReader() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orgUID := r.Header.Get("Organization")

			// put it in context and call the next with our new context
			ctx := context.WithValue(r.Context(), orgCtxKey, orgUID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
