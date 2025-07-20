package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
	PasswordUpdate(id int, crrentPassword, newPassword string) error
}

type User struct {
	Created        time.Time
	Email          string
	HashedPassowrd []byte
	ID             int
	Name           string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Get(id int) (*User, error) {
	var user User
	stmt := `SELECT id, name, email, created FROM users WHERE id = ?;`
	err := m.DB.QueryRow(stmt, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Created,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassowrd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `
		INSERT INTO users (name, email, hashed_password, created)
		VALUES (?, ?, ?, UTC_TIMESTAMP());
	`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassowrd))
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if mysqlError.Number == 1062 &&
				strings.Contains(mysqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassowrd []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?;`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassowrd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassowrd), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?);`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var currentHashedPasswrod []byte
	stmt := `SELECT hashed_password FROM users WHERE id = ?;`
	err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPasswrod)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(currentHashedPasswrod, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}
	hashedPassowrd, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	stmt = `UPDATE users SET hashed_password = ? WHERE id = ?;`
	_, err = m.DB.Exec(stmt, string(hashedPassowrd), id)
	return err
}
