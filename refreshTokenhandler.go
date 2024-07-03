package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Refresh Token Handler
func (apiCfg apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("jwt")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token not found")
		return
	}
	refreshToken := refreshTokenCookie.Value
	claims := &jwt.MapClaims{}

	jwtKey := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	userID := (*claims)["user_id"].(string)

	email, err := apiCfg.DB.GetUserEmail(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("email does not exist:%v", err))
		return
	}
	accessToken, err := generateJWT(userID, 15*time.Minute, "JWT_SECRET")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error generating access token: %v", err))
		return
	}

	response := map[string]interface{}{
		"accessToken": accessToken,
		"user":        databaseUserToUser(email),
	}

	respondWithJson(w, http.StatusOK, response)
}

func verifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Printf("authHeader:%v", authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// Store user information in context
		username := (*claims)["user_id"].(string)
		ctx := context.WithValue(r.Context(), "user", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
