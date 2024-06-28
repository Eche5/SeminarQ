package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	respondWithJson(w, 200, databaseUserToUser(newUser))
}
