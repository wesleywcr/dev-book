package models

import "time"

type User struct {
	ID         uint64    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Nickname   string    `json:"nickname,omitempty"`
	Email      string    `json:"email,omitempty"`
	Passoword  string    `json:"password,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
}
