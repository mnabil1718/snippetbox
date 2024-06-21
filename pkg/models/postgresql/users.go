package postgresql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mnabil1718/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert(name, email, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) // cost below 12 is not recommended
	if err != nil {
		return err
	}

	SQL := "INSERT INTO users (name, email, password, created_at) VALUES ($1, $2, $3, NOW())"
	_, err = model.DB.Exec(SQL, name, email, hashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_email_key") {
				return models.ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (model *UserModel) Authenticate(email, password string) (int, error) {
	user := &models.User{}
	SQL := "SELECT id,email,password,active FROM users WHERE email=$1 AND active=TRUE"
	err := model.DB.QueryRow(SQL, email).Scan(&user.ID, &user.Email, &user.Password, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return user.ID, err
	}

	// remember, if struct is pointer the fields are not
	return user.ID, nil
}

func (model *UserModel) Get(id int) (*models.User, error) {
	user := &models.User{}

	SQL := "SELECT id, name, email, active, created_at FROM users WHERE id=$1 AND active=TRUE"
	err := model.DB.QueryRow(SQL, id).Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return user, nil
}
