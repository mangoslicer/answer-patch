package models

import (
	"time"
)

type User struct {
	ID             string    `json:"userID"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashedPassword"`
	CreatedAt      time.Time `json:"createdAt"`
}
