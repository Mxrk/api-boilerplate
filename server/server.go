package server

import (
	"encoding/json"
	"net/http"

	"api-boilerplate/models/domain"
	"api-boilerplate/models/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	server *http.Server
	router *chi.Mux

	UserService domain.UserService
}

func InitServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: chi.NewRouter(),
	}

	s.router.Use(middleware.Logger)

	// CORS settings for local testing
	s.router.Use(cors.Handler(cors.Options{
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
	s.router.Group(func(r chi.Router) {
		r.Use(Verifier(tokenAuth))

		r.Use(Authenticator)

		r.Get("/protected/health", protectedHealth)
	})

	// Public routes
	s.router.Get("/health", health)
	s.router.Post("/login", s.login)
	s.router.Post("/register", s.register)

	return s
}

// StartServer validates the server options and begins listening on the bind address.
func (s *Server) StartServer(adr string) (err error) {
	return http.ListenAndServe(adr, s.router)

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
