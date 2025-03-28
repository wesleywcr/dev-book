package models

import "time"

type Publication struct {
	ID              uint64    `json:"id,omitempty"`
	Title           string    `json:"title,omitempty"`
	Content         string    `json:"content,omitempty"`
	AuthorID        uint64    `json:"authorId,omitempty"`
	AuthorNickaname string    `json:"authorNickname,omitempty"`
	Likes           uint64    `json:"likes"`
	Created_at      time.Time `json:"created_at,omitempty"`
}
