package mysql

import (
	"database/sql"
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

	return nil
}

// Authenticate the user.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Retreive the id and hashed password associated with the given email
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)

	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// Check whether the hashed password and plain-text password match
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// Password is correct. Return the user ID
	return id, nil
}

// Get returns a specific user based on id.
func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}

	return s, nil
}

// This will return the 10 most recently created snippets.
// func (m *UserModel) Latest() ([]*models.User, error) {
// 	return nil, nil
// }
