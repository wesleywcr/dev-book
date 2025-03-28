package repositories

import (
	"database/sql"

	"github.com/wesleywcr/dev-book/api/models"
)

type Publications struct {
	db *sql.DB
}

func NewRepositoryOfPublications(db *sql.DB) *Publications {
	return &Publications{db}
}

func (repository Publications) Create(publications models.Publication) (uint64, error) {
	statement, error := repository.db.Prepare(
		"insert into publications (title, content, author_id) values (?, ?, ?)",
	)
	if error != nil {
		return 0, error
	}

	result, error := statement.Exec(publications.Title, publications.Content, publications.AuthorID)
	if error != nil {
		return 0, error
	}
	lastIdInsert, error := result.LastInsertId()

	if error != nil {
		return 0, error
	}

	return uint64(lastIdInsert), nil
}
