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
		"update users set name = ?, nickname = ?, email = ? where id = ?",
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

func (repository Users) Delete(ID uint64) error {
	statement, error := repository.db.Prepare(
		"delete from users where id = ?",
	)
	if error != nil {
		return error
	}
	defer statement.Close()

	if _, error = statement.Exec(ID); error != nil {
		return error
	}

	return nil
}

func (repository Users) SearchEmail(email string) (models.User, error) {

	row, error := repository.db.Query(
		"select id, password from users where email = ?", email)
	if error != nil {
		return models.User{}, error
	}
	defer row.Close()

	var user models.User

	if row.Next() {
		if error = row.Scan(&user.ID, &user.Password); error != nil {
			return models.User{}, error
		}
	}
	return user, error
}
func (repository Users) Follow(userId, followerId uint64) error {
	statement, error := repository.db.Prepare(
		"insert ignore into followers (user_id, follower_id) values (?, ?)",
	)
	if error != nil {
		return error
	}
	defer statement.Close()

	if _, error = statement.Exec(userId, followerId); error != nil {
		return error
	}
	return nil
}
func (repository Users) UnFollow(userId, followerId uint64) error {
	statement, error := repository.db.Prepare(
		"delete from followers where user_id = ? and follower_id = ?",
	)
	if error != nil {
		return error
	}
	defer statement.Close()

	if _, error = statement.Exec(userId, followerId); error != nil {
		return error
	}
	return nil
}

func (repository Users) SearchFollowers(userId uint64) ([]models.User, error) {
	rows, error := repository.db.Query(`
	select u.id, u.name, u.nickname, u.email, u.created_at
	from users u inner join followers s on u.id = s.follower_id where s.user_id = ?
`, userId)
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
func (repository Users) SearchFollowing(userId uint64) ([]models.User, error) {
	rows, error := repository.db.Query(`
	select u.id, u.name, u.nickname, u.email, u.created_at
	from users u inner join followers s on u.id = s.user_id where s.follower_id = ?
`, userId)

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

func (repository Users) GetPassword(userId uint64) (string, error) {
	row, error := repository.db.Query(`select password from users where id = ?`, userId)
	if error != nil {
		return "", error
	}
	defer row.Close()

	var user models.User

	if row.Next() {
		if error = row.Scan(&user.Password); error != nil {
			return "", error
		}
	}
	return user.Password, nil
}
func (repository Users) UpdatePassword(userId uint64, password string) error {
	statement, error := repository.db.Prepare("update users set password = ? where id = ?")
	if error != nil {
		return error
	}
	defer statement.Close()

	if _, error := statement.Exec(password, userId); error != nil {
		return error
	}
	return nil
}
