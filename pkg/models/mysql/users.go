package mysql

import (
	"database/sql"

	"github.com/arammikayelyan/snippetbox/pkg/models"
)

// UserModel type wraps a sql.DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// This will insert a new user into the database.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Insert a new user into the database.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get returns a specific user based on id.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

// This will return the 10 most recently created snippets.
// func (m *UserModel) Latest() ([]*models.User, error) {
// 	return nil, nil
// }
