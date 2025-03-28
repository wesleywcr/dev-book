package models

import (
	"errors"
	"strings"
	"time"
)

type Publication struct {
	ID              uint64    `json:"id,omitempty"`
	Title           string    `json:"title,omitempty"`
	Content         string    `json:"content,omitempty"`
	AuthorID        uint64    `json:"authorId,omitempty"`
	AuthorNickaname string    `json:"authorNickname,omitempty"`
	Likes           uint64    `json:"likes"`
	Created_at      time.Time `json:"created_at,omitempty"`
}

func (publication *Publication) Prepare() error {
	if error := publication.validate(); error != nil {
		return error
	}
	publication.format()
	return nil
}

func (publication *Publication) validate() error {
	if publication.Title == "" {
		return errors.New("O título é obrigatório e não pode estar em branco")
	}
	if publication.Content == "" {
		return errors.New("O conteúdo é obrigatório e não pode estar em branco")
	}

	return nil
}

func (publication *Publication) format() {
	publication.Title = strings.TrimSpace(publication.Title)
	publication.Content = strings.TrimSpace(publication.Content)
}
