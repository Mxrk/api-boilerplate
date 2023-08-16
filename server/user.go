package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"api-boilerplate/models/dto"
	"github.com/go-chi/jwtauth/v5"
)

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request dto.LoginRequest

	jsonDecodeError := json.NewDecoder(r.Body).Decode(&request)
	if jsonDecodeError != nil {
		JSONError(w, "Please provide a body.", http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.Password == "" {
		JSONError(w, "Please enter a username/password.", http.StatusBadRequest)
		return
	}

	// check if username exists
	user, userExists := s.UserService.GetUserFromUsername(r.Context(), request.Username)
	if !userExists {
		JSONError(w, "Invalid login data.", http.StatusNotFound)
		return
	}

	jsonDecodeError = CheckPasswordHash(request.Password, user.Password)
	if jsonDecodeError != nil {
		JSONError(w, "Invalid login data.", http.StatusBadRequest)
		return
	}

	var claims = map[string]interface{}{"id": user.ID, "username": user.Username}

	jwtauth.SetExpiryIn(claims, 7*24*time.Hour)

	_, tokenString, _ := tokenAuth.Encode(claims)

	var serverResponse dto.UserTokenResponse
	serverResponse.User = dto.UserRequest{
		Token:    tokenString,
		Username: user.Username,
	}

	jsonDecodeError = json.NewEncoder(w).Encode(serverResponse)
	if jsonDecodeError != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("no body", err)
		JSONError(w, "Please provide a body.", http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.Password == "" {
		JSONError(w, "Please enter a username/password", http.StatusBadRequest)
		return
	}

	// check if username exists
	_, userExists := s.UserService.GetUserFromUsername(r.Context(), request.Username)
	if userExists {
		JSONError(w, "User already exists.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(request.Password)
	if err != nil {
		JSONError(w, "Could not create user.", http.StatusInternalServerError)
		return
	}

	userID, err := s.UserService.CreateUser(r.Context(), request.Username, hashedPassword)
	if err != nil {
		JSONError(w, "Could not create user.", http.StatusInternalServerError)
		return
	}

	user, err := s.UserService.GetUser(r.Context(), userID)
	if err != nil {
		JSONError(w, "Could not create user.", http.StatusInternalServerError)
		return
	}

	var claims = map[string]interface{}{"id": user.ID, "username": user.Username}

	jwtauth.SetExpiryIn(claims, 7*24*time.Hour)

	_, tokenString, _ := tokenAuth.Encode(claims)

	var serverResponse dto.UserTokenResponse
	serverResponse.User = dto.UserRequest{
		Token:    tokenString,
		Username: user.Username,
	}

	err = json.NewEncoder(w).Encode(serverResponse)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func readUserID(ctx context.Context) int {
	return int(ctx.Value("user_id").(float64))
}
