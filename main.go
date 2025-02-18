package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Eche5/SeminarQ/internal/database"

	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dBURL := os.Getenv("DB_URL")
	if dBURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	conn, err := sql.Open("postgres", dBURL)
	if err != nil {
		log.Fatal("can't connect to database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
    apiCfg.startCronJobs()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "https://queryhub.adaptable.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Get("/refresh-token", apiCfg.handlerRefreshToken)
	v1Router.Get("/seminar/{apiKey}", apiCfg.handlerGetSeminarByAPIKey)

	v1Router.Post("/users", apiCfg.handlerCreateUsers)
	v1Router.Post("/seminar", apiCfg.handlerCreateSeminar)
	v1Router.Get("/seminars/{userId}", apiCfg.handlerGetAllSeminars)
	v1Router.Put("/seminar/{seminarId}", apiCfg.handlerUpdateSeminarName)
	v1Router.Get("/seminar/{userId}/{seminarName}", apiCfg.handlerGetSeminarByName)

	v1Router.Delete("/seminar/{seminarId}", apiCfg.handlerDeleteSeminar)

	v1Router.Post("/question/{seminarId}", apiCfg.handlerCreateQuestion)
	v1Router.Get("/question/{seminarId}", apiCfg.handlerGetAllQuestions)

	v1Router.Post("/login", apiCfg.handlerLoginUser)
	v1Router.Post("/logout", apiCfg.haandleLogoutUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Printf("Server running on PORT %v\n", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
