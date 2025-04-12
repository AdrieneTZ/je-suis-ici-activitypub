package middlewares

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// AuthJWT verify JWT token validity and user identity
func AuthJWT(tokenAuth *jwtauth.JWTAuth) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// step1. verify JWT token validity
		verifier := jwtauth.Verifier(tokenAuth)

		// step2. verify user identity
		authenticator := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token, _, err := jwtauth.FromContext(r.Context())

				if err != nil {
					http.Error(w, fmt.Sprintf("fail to extract JWT token from context: %v", err), http.StatusUnauthorized)
					return
				}

				if token == nil {
					http.Error(w, "invalid or missing token", http.StatusUnauthorized)
					return
				}

				// Token is valid, proceed
				next.ServeHTTP(w, r)
			})
		}

		// step3. combine verify
		return verifier(authenticator(next))
	}
}
