package models

import (
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/wesleywcr/dev-book/api/security"
)

type User struct {
	ID         uint64    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Nickname   string    `json:"nickname,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("O nome é um campo obrigatório")
	}
	if user.Nickname == "" {
		return errors.New("Nickname é um campo obrigatório")
	}
	if user.Email == "" {
		return errors.New("E-mail é um campo  obrigatório")
	}

	if error := checkmail.ValidateFormat(user.Email); error != nil {
		return errors.New("O e-mail inserido é invalido")
	}

	if step == "register" && user.Password == "" {
		return errors.New("Senha é um campo obrigatório")
	}

	return nil
}

// validate and format user
func (user *User) Prepare(step string) error {
	if error := user.validate(step); error != nil {
		return error
	}
	if error := user.formatted(step); error != nil {
		return error
	}
	return nil
}

func (user *User) formatted(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Email = strings.TrimSpace(user.Email)

	if step == "register" {
		passwordWithHash, error := security.Hash(user.Password)
		if error != nil {
			return error
		}
		user.Password = string(passwordWithHash)
	}
	return nil
}
