package model

import (
	"time"
)

type APIError struct {
	Err string `json:"error_code"`
	Msg string `json:"message"`
}

type UserAccount struct {
	UserID       string
	Email        string `json:"email"`
	Password     string `json:"password"`
	CreationDate time.Time
}

type QA struct {
	QAID      string `json:"qa_id"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Extension string `json:"extension"`
}

type Story struct {
	StoryID      string    `json:"story_id"`
	UserID       string    `json:"user_id"`
	CreationDate time.Time `json:"creation_date"`
	StoryContext []int     `json:"story_context"`
	Content      []QA      `json:"content"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Context  []int  `json:"context"`
}
