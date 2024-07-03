package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Eche5/SeminarQ/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Define a structure to hold active WebSocket connections
type WebSocketManager struct {
	connections map[*websocket.Conn]bool
	lock        sync.Mutex
}

// Create a new WebSocket manager
var wsManager = WebSocketManager{
	connections: make(map[*websocket.Conn]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections (for demo purposes)
		return true
	},
}

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

	// Broadcast the new question to all active WebSocket connections
	broadcastQuestion(question)

	respondWithJson(w, http.StatusCreated, databaseQuestionToQuestion(question))
}

func broadcastQuestion(question database.Question) {
	wsManager.lock.Lock()
	defer wsManager.lock.Unlock()

	for conn := range wsManager.connections {
		if err := conn.WriteJSON(databaseQuestionToQuestion(question)); err != nil {
			log.Printf("WebSocket broadcast error: %v", err)
			conn.Close()
			delete(wsManager.connections, conn)
		}
	}
}

func (apiCfg apiConfig) handlerGetAllQuestions(w http.ResponseWriter, r *http.Request) {
	seminarIdStr := chi.URLParam(r, "seminarId")
	seminarId, err := uuid.Parse(seminarIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid seminar id %v", err))
		return
	}


	// Check if the request is a WebSocket upgrade request
	if websocket.IsWebSocketUpgrade(r) {
		handleWebSocket(apiCfg, w, r, seminarId)
		return
	}

	// Handle normal HTTP GET request
	questions, err := apiCfg.DB.GetAllQuestion(r.Context(), seminarId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error:%v", err))
		return
	}
	respondWithJson(w, http.StatusOK, databaseQuestionToQuestionArray(questions))
}

func handleWebSocket(apiCfg apiConfig, w http.ResponseWriter, r *http.Request, seminarId uuid.UUID) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Register the new connection
	wsManager.lock.Lock()
	wsManager.connections[conn] = true
	wsManager.lock.Unlock()

	defer func() {
		// Unregister the connection on close
		wsManager.lock.Lock()
		delete(wsManager.connections, conn)
		wsManager.lock.Unlock()
		conn.Close()
	}()

	// Send initial data to client
	questions, err := apiCfg.DB.GetAllQuestion(r.Context(),seminarId)
	if err != nil {
		log.Println(err)
		return
	}

	if err := conn.WriteJSON(databaseQuestionToQuestionArray(questions)); err != nil {
		log.Println(err)
		return
	}

	// Listen for incoming messages from client
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
	}
}
