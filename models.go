package main

import (
	"time"

	"github.com/Eche5/SeminarQ/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		FullName:  dbUser.FullName,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
	}
}

type Seminar struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseSeminarToSeminar(dbBase database.Seminar) Seminar {
	return Seminar{
		ID:        dbBase.ID,
		CreatedAt: dbBase.CreatedAt,
		UpdatedAt: dbBase.UpdatedAt,
		Name:      dbBase.Name,
		ApiKey:    dbBase.ApiKey,
		UserID:    dbBase.UserID,
	}
}

type Question struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	SeminarID uuid.UUID `json:"seminar_id"`
	Question  string    `json:"question"`
}

func databaseQuestionToQuestion(dbBase database.Question) Question {
	return Question{
		ID:        dbBase.ID,
		CreatedAt: dbBase.CreatedAt,
		UpdatedAt: dbBase.UpdatedAt,
		UserID:    dbBase.UserID,
		SeminarID: dbBase.SeminarID,
		Question:  dbBase.Question,
	}
}


func databaseQuestionToQuestionArray(dbBases []database.Question)[]Question{
	questions := []Question{}
	for _, dbBase :=range dbBases{
		questions = append(questions, databaseQuestionToQuestion(dbBase))
	}
	return questions
}

func databaseSeminarToSeminarArray(dbBases []database.Seminar)[]Seminar{
	seminars := []Seminar{}
	for _, dbBase :=range dbBases{
		seminars = append(seminars, databaseSeminarToSeminar(dbBase))
	}
	return seminars
}