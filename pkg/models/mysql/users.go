package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/arammikayelyan/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

// UserModel type wraps a sql.DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert a new user into the database.
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password
	hp, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) 
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Exec() inserts the user details and hashed password into the
	// users table. If it returns an error, we type assert it to a
	// *mysql.MySQLError object so we can check if the error number
	// is 1062 and, if it is, we also check whether or not it relates
	// to our users_uc_email key by checking the contents of the error
	// message. If it does, we return an ErrDuplicateEmail error. Otherwise,
	// we return the original error (or nil if everything worked).
	_, err = m.DB.Exec(stmt, name, email, string(hp))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "") {
				return models.ErrDuplicateEmail
			}
		}
	}

	fmt.Println("insert new user")

	return nil
}

// Authenticate the user.
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
