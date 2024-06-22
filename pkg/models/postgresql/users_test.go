package postgresql

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mnabil1718/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {

	if testing.Short() {
		t.Skip("postgresql: skipping integration test...")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:        1,
				Name:      "Alice",
				Email:     "alice@email.com",
				CreatedAt: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:    true,
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Non-existent ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			SQL := "INSERT INTO users (name, email, password, active, created_at) VALUES ($1, $2, $3, $4, $5)"

			db, cleanUpFunc := createAndInsertTable(t, "users", SQL, "Alice", "alice@email.com", "okiedokie", true, time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC))
			defer cleanUpFunc()

			userModel := &UserModel{DB: db}

			user, err := userModel.Get(test.userID)

			if !errors.Is(err, test.wantError) {
				t.Errorf("want %v; got %s", test.wantError, err)
			}
			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("want %v; got %v", test.wantUser, user)
			}
		})
	}

}
