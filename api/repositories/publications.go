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

func (repository Publications) SearchById(publicationId uint64) (models.Publication, error) {
	row, error := repository.db.Query(`
	select p.*, p.nickname from 
	publications p inner join users u
	on u.id = p.author_id where p.id = ?
	`, publicationId)
	if error != nil {
		return models.Publication{}, error
	}

	defer row.Close()

	var publication models.Publication

	if row.Next() {
		if error = row.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.Created_at,
			&publication.AuthorNickaname,
		); error != nil {
			return models.Publication{}, error
		}
	}
	return publication, nil
}
func (repository Publications) Search(userId uint64) ([]models.Publication, error) {
	rows, error := repository.db.Query(`
	select distinct p.*, u.nickname from publications p
	inner join users u on u.id = p.author_id
	inner join followers s on p.author_id = s.user_id
	where u.id = ? or s.follower_id = ?
	order by 1 desc
	`, userId, userId)
	if error != nil {
		return nil, error
	}

	var publications []models.Publication

	if rows.Next() {
		var publication models.Publication

		if error = rows.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.Created_at,
			&publication.AuthorNickaname,
		); error != nil {
			return nil, error
		}
		publications = append(publications, publication)
	}
	return publications, nil
}
