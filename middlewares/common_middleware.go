package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"viraj_golang/utils"
)

type Middleware func(http.Handler) http.Handler

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("\n Method: %v , Request Origin : %v , time_taken : %v , path : %v \n", r.Method, r.RemoteAddr, time.Since(start), r.URL.Path)

	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Connection", "close")
				errr := fmt.Errorf("%s", r)
				http.Error(w, errr.Error(), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := utils.ExtractJwtFromAuthHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNetworkAuthenticationRequired)
			return
		}
		claims, err := utils.ValidateJwtAndExtractClaims(jwtToken)
		if err != nil {

			http.Error(w, err.Error(), http.StatusNetworkAuthenticationRequired)
			return
		}
		type UserIdClaim string
		var userId UserIdClaim
		ctx := r.Context()
		ctx = context.WithValue(ctx, userId, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func CreateMiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}
		return next
	}
}
