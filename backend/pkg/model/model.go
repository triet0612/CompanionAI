package model

import (
	"time"

	"github.com/google/uuid"
)

type APIError struct {
	Err string `json:"error_code"`
	Msg string `json:"message"`
}

type User struct {
	ID       uuid.UUID
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChatPair struct {
	PairID     uuid.UUID `json:"pairid"`
	Request    string    `json:"request"`
	Response   string    `json:"response"`
	Extension  string
	Attachment []byte
}

type Chat struct {
	ChatID       uuid.UUID  `json:"chat_id"`
	OwnerID      uuid.UUID  `json:"owner_id"`
	CreationDate time.Time  `json:"creation_date"`
	Content      []ChatPair `json:"content"`
}
