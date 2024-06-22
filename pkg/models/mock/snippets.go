package mock

import (
	"time"

	"github.com/mnabil1718/snippetbox/pkg/models"
)

var mockSnippet *models.Snippet = &models.Snippet{
	ID:        1,
	Title:     "Test Snippet",
	Content:   "This is a test snippet.",
	CreatedAt: time.Now(),
	ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
}

type SnippetModel struct{}

func (model *SnippetModel) Insert(title, content, expiresAt string) (int, error) {
	return 2, nil
}

func (model *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (model *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
