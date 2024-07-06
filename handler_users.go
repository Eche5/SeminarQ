package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Eche5/SeminarQ/internal/database"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing json:%v", err))
		return
	}

	if len(params.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "password must be at least 8 characters long")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password: %v", err))
		return
	}

	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FullName:  params.FullName,
		Email:     params.Email,
		Password:  string(hashedPassword),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("cannot create user:%v", err))
		return
	}

	accessToken, err := generateJWT(newUser.ID.String(), 15*time.Minute, "JWT_SECRET")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error generating token: %v", err))
		return
	}
	refreshToken, err := generateJWT(newUser.ID.String(), 7*24*time.Hour, "JWT_REFRESH_SECRET")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error generating token: %v", err))
		return
	}
	response := map[string]interface{}{
		"user":  databaseUserToUser(newUser),
		"token": accessToken,
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	respondWithJson(w, 200, response)
}

// Login User
func (apiCfg apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing json:%v", err))
		return
	}
	email, err := apiCfg.DB.GetUserEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("email does not exist"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(email.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	accessToken, err := generateJWT(params.Email, 15*time.Minute, "JWT_SECRET")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error generating token: %v", err))
		return
	}
	refreshToken, err := generateJWT(params.Email, 7*24*time.Hour, "JWT_REFRESH_SECRET")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error generating token: %v", err))
		return
	}
	response := map[string]interface{}{
		"user":        databaseUserToUser(email),
		"accessToken": accessToken,
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	respondWithJson(w, 200, response)
}

func generateJWT(userID string, expirationTime time.Duration, secret string) (string, error) {
	var jwtKey = []byte(os.Getenv(secret))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expirationTime).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (apiCfg apiConfig) haandleLogoutUser(w http.ResponseWriter, r *http.Request) {
_, err := r.Cookie("jwt")
	if err != nil {
		respondWithError(w, http.StatusNoContent, "")
		return
	}
	
	// Invalidate the cookie by setting a past expiration date
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // Ensure SameSite is consistent with login cookie settings
		Secure:   true,                  // Ensure Secure is consistent with login cookie settings
	})

	respondWithJson(w, http.StatusNoContent, "")
}
