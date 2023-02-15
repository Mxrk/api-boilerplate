package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"api-boilerplate/database"
	"api-boilerplate/models/dto"
	"github.com/go-chi/jwtauth/v5"
)

func login(w http.ResponseWriter, r *http.Request) {
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
	databaseUser, userExists := database.GetUserFromUsername(request.Username)
	if !userExists {
		JSONError(w, "Invalid login data.", http.StatusNotFound)
		return
	}

	jsonDecodeError = CheckPasswordHash(request.Password, databaseUser.Password)
	if jsonDecodeError != nil {
		JSONError(w, "Invalid login data.", http.StatusBadRequest)
		return
	}

	var claims = map[string]interface{}{"id": databaseUser.ID, "username": databaseUser.Username}

	jwtauth.SetExpiryIn(claims, 7*24*time.Hour)

	_, tokenString, _ := tokenAuth.Encode(claims)

	var serverResponse dto.UserTokenResponse
	serverResponse.User = dto.UserRequest{
		Token:    tokenString,
		Username: databaseUser.Username,
	}

	jsonDecodeError = json.NewEncoder(w).Encode(serverResponse)
	if jsonDecodeError != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func register(w http.ResponseWriter, r *http.Request) {
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
	_, userExists := database.GetUserFromUsername(request.Username)
	if userExists {
		JSONError(w, "User already exists.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(request.Password)
	if err != nil {
		JSONError(w, "Could not create user.", http.StatusInternalServerError)
		return
	}

	databaseUser, err := database.CreateUser(request.Username, hashedPassword)
	if err != nil {
		JSONError(w, "Could not create user.", http.StatusInternalServerError)
		return
	}

	var claims = map[string]interface{}{"id": databaseUser.ID, "username": databaseUser.Username}

	jwtauth.SetExpiryIn(claims, 7*24*time.Hour)

	_, tokenString, _ := tokenAuth.Encode(claims)

	var serverResponse dto.UserTokenResponse
	serverResponse.User = dto.UserRequest{
		Token:    tokenString,
		Username: databaseUser.Username,
	}

	err = json.NewEncoder(w).Encode(serverResponse)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userID := readUserID(r.Context())

	user, err := database.GetUser(userID)
	if err != nil {
		JSONError(w, "could not get user.", http.StatusInternalServerError)
		log.Println("could not get user", err)
		return
	}

	userResp := dto.DatabaseUserToUserResponse(user)
	json.NewEncoder(w).Encode(userResp)
}

func readUserID(ctx context.Context) int {
	return int(ctx.Value("user_id").(float64))
}
