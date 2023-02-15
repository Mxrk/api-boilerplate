package server

import (
	"encoding/json"
	"log"
	"net/http"

	"api-boilerplate/config"
	"api-boilerplate/models/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func InitServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// CORS settings for local testing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           24, // Maximum value not ignored by any of major browsers
	}))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(Verifier(tokenAuth))

		r.Use(Authenticator)

		r.Get("/protected/health", protectedHealth)
	})

	// Public routes
	r.Get("/health", health)
	r.Post("/login", login)
	r.Post("/register", register)

	log.Println("Server started on port " + config.Config.Server.Port)
	err := http.ListenAndServe(":"+config.Config.Server.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}

// JSONError is a helper function which is creating an error as json.
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	var errorResponse dto.ErrorResponse
	errorResponse.ErrorType = err.(string)
	errorResponse.StatusCode = code

	encodeError := json.NewEncoder(w).Encode(errorResponse)
	if encodeError != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
