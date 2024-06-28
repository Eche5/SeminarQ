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

func (apiCfg apiConfig) handlerCreateQuestion(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Question string `json:"question"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing json:%v", err))
		return
	}
	seminarIdStr := chi.URLParam(r, "seminarId")
	seminarId, err := uuid.Parse(seminarIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid seminar id %v", err))
		return
	}
	userIdStr := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid seminar id %v", err))
		return
	}

	question, err := apiCfg.DB.CreateQuestion(r.Context(), database.CreateQuestionParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userId,
		SeminarID: seminarId,
		Question:  params.Question,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error:%v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseQuestionToQuestion(question))
}

func (apiCfg apiConfig) handlerGetAllQuestions(w http.ResponseWriter, r *http.Request) {
	seminarIdStr := chi.URLParam(r, "seminarId")
	seminarId, err := uuid.Parse(seminarIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid seminar id %v", err))
		return
	}
	userIdStr := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid seminar id %v", err))
		return
	}
	questions, err := apiCfg.DB.GetAllQuestion(r.Context(), database.GetAllQuestionParams{
		SeminarID: seminarId,
		UserID:    userId,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error:%v", err))
		return
	}
	respondWithJson(w, http.StatusOK, databaseQuestionToQuestionArray(questions))
}
