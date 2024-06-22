package mock

import (
	"time"

	"github.com/mnabil1718/snippetbox/pkg/models"
)

type UserModel struct{}

var mockUser = &models.User{
	ID:        1,
	Email:     "alice@gmail.com",
	Active:    true,
	Name:      "Alice",
	CreatedAt: time.Now(),
}

func (model *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@email.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (model *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "alice@gmail.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (model *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
