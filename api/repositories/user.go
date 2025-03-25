package repositories

import (
	"database/sql"

	"github.com/wesleywcr/dev-book/api/models"
)

type Users struct {
	db *sql.DB
}

func NewRepositoryOfUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (repository Users) Create(user models.User) (uint64, error) {
	statement, error := repository.db.Prepare(
		"insert into users (name, nickname, email, password) values(?, ?, ?, ?)",
	)
	if error != nil {
		return 0, nil
	}

	defer statement.Close()

	result, error := statement.Exec(user.Name, user.Nickname, user.Email, user.Password)
	if error != nil {
		return 0, error
	}

	lastInsertId, error := result.LastInsertId()
	if error != nil {
		return 0, error
	}
	return uint64(lastInsertId), nil
}
