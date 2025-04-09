package middlewares

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// AuthJWT verify JWT token validity and user identity
func AuthJWT(tokenAuth *jwtauth.JWTAuth) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// step1. verify JWT token validity
		verifier := jwtauth.Verifier(tokenAuth)

		// step2. verify user identity
		authenticator := jwtauth.Authenticator(tokenAuth)

		// step3. combine verify
		return authenticator(verifier(next))
	}
}
