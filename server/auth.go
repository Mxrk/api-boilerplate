package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

var tokenAuth *jwtauth.JWTAuth
var signKey = []byte("secret")

func init() {
	tokenAuth = jwtauth.New("HS256", signKey, nil)
	randomUserID := 123
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": randomUserID})
	log.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

// Verifier is the JWTAuth middleware.
func Verifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verify(ja, TokenFromQuery, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)(next)
	}
}

// TokenFromQuery exports the key "token" from a given query.
func TokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("token")
}

// Authenticator authenticates a user.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, claims, err := jwtauth.FromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "user_id", claims["id"])
		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HashPassword hashes a password and returns the hashed password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares the password with a given hash.
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
