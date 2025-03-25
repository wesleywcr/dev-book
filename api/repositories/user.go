package repositories

import (
	"database/sql"
	"fmt"

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

func (repository Users) Search(nameOrNickname string) ([]models.User, error) {
	nameOrNickname = fmt.Sprintf("%%%s%%", nameOrNickname) // %nameOrNickname%

	rows, error := repository.db.Query(
		"select id, name, nickname, email, created_at from users where name LIKE ? or nickname LIKE ?",
		nameOrNickname, nameOrNickname,
	)
	if error != nil {
		return nil, error
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if error = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.Created_at,
		); error != nil {
			return nil, error
		}
		users = append(users, user)
	}
	return users, nil

}

func (repository Users) SearchPerId(ID uint64) (models.User, error) {
	rows, error := repository.db.Query(
		"select id, name, nickname, email, created_at from users where id = ?", ID,
	)

	if error != nil {
		return models.User{}, error
	}

	defer rows.Close()

	var user models.User

	if rows.Next() {
		if error = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.Created_at,
		); error != nil {
			return models.User{}, error
		}
	}
	return user, nil
}

func (repository Users) Update(ID uint64, user models.User) error {
	statement, error := repository.db.Prepare(
		"update users set name = ?, nickname = ?, email = ?, where id = ?",
	)
	if error != nil {
		return error
	}
	defer statement.Close()
	if _, error := statement.Exec(user.Name, user.Nickname, user.Email, ID); error != nil {
		return error
	}
	return nil
}
