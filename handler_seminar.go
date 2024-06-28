package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Eche5/SeminarQ/internal/database"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateSeminar(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		UserID string `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing json:%v", err))
		return
	}

		userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid user_id: %v", err))
		return
	}
	seminar,err:=apiCfg.DB.CreateSeminar(r.Context(), database.CreateSeminarParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		UserID:    userID,
	})
	if err !=nil{
		respondWithError(w,400,fmt.Sprintf("something went wrong: %v",err))
		return
	}
	respondWithJson(w, 201, databaseSeminarToSeminar(seminar))

}
