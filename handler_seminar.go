package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Eche5/SeminarQ/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateSeminar(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name   string `json:"name"`
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
	seminar, err := apiCfg.DB.CreateSeminar(r.Context(), database.CreateSeminarParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		UserID:    userID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("something went wrong: %v", err))
		return
	}
	respondWithJson(w, 201, databaseSeminarToSeminar(seminar))

}

func (apiCfg apiConfig) handlerGetAllSeminars(w http.ResponseWriter, r *http.Request) {

	userIdStr := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid user id %v", err))
		return
	}
	seminars, err := apiCfg.DB.GetAllSeminars(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error:%v", err))
		return
	}
	respondWithJson(w, http.StatusOK, databaseSeminarToSeminarArray(seminars))
}

func (apiCfg apiConfig) handlerGetSeminarByAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apiKey")

	seminar, err := apiCfg.DB.GetAllSeminarsByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Something went wrong: %v", err))
		return
	}
	var result []Seminar
	for _, seminar := range seminar {
		result = append(result, databaseSeminarToSeminar(seminar))
	}

	respondWithJson(w, http.StatusOK, result)
}

func (apiCfg apiConfig) handlerDeleteSeminar(w http.ResponseWriter, r *http.Request) {
	seminadIdStr := chi.URLParam(r, "seminarId")
	seminarId, err := uuid.Parse(seminadIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid user id %v", err))
		return
	}

	err = apiCfg.DB.DeleteSeminar(r.Context(), seminarId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error:%v", err))
	}
	respondWithJson(w, http.StatusNoContent, struct{}{})
}

func (apiCfg apiConfig) handlerUpdateSeminarName(w http.ResponseWriter, r *http.Request) {
	seminarIdStr := chi.URLParam(r, "seminarId")
	seminarId, err := uuid.Parse(seminarIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid seminar id %v", err))
		return
	}
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing json:%v", err))
		return
	}

	updatedSeminar, err := apiCfg.DB.EditSeminarName(r.Context(), database.EditSeminarNameParams{
		ID:        seminarId,
		Name:      params.Name,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Something went wrong: %v", err))
	}

	respondWithJson(w, http.StatusOK, databaseSeminarToSeminar(updatedSeminar))
}

func (apiCfg apiConfig) handlerGetSeminarByName(w http.ResponseWriter, r *http.Request) {

	seminarNameStr := chi.URLParam(r, "seminarName")
	userIdStr := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid user id %v", err))
		return
	}
	seminar, err := apiCfg.DB.GetSeminarByName(r.Context(), database.GetSeminarByNameParams{
		Name:   seminarNameStr,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Something went wrong: %v", err))
	}

	var result []Seminar
	for _, seminar := range seminar {
		result = append(result, databaseSeminarToSeminar(seminar))
	}

	respondWithJson(w, http.StatusOK, result)
}
