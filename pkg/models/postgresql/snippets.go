package postgresql

import (
	"database/sql"
	"errors"

	"github.com/mnabil1718/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content, expiresAt string) (int, error) {
	var id int
	SQL := "INSERT INTO snippets (title, content, created_at, expires_at) VALUES ($1, $2, NOW(), NOW() + $3::INTERVAL) RETURNING id" // need to cast $3 into interval type
	err := model.DB.QueryRow(SQL, title, content, expiresAt+" days").Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (model *SnippetModel) Get(id int) (*models.Snippet, error) {

	snippet := &models.Snippet{}
	SQL := "SELECT id, title, content, created_at, expires_at FROM snippets WHERE id=$1 AND expires_at > NOW()"
	err := model.DB.QueryRow(SQL, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.CreatedAt, &snippet.ExpiresAt) // has to pass each field's memory address
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return snippet, nil
}

func (model *SnippetModel) Latest() ([]*models.Snippet, error) {
	snippets := []*models.Snippet{}
	SQL := "SELECT id, title, content, created_at, expires_at FROM snippets WHERE expires_at > NOW() ORDER BY created_at LIMIT 10"
	rows, err := model.DB.Query(SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		snippet := &models.Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.CreatedAt, &snippet.ExpiresAt)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
