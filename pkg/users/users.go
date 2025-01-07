package users

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserCreateInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Inactive string `json:"inactive"`
	CreatdAt string `json:"createdAt"`
}

type usrService struct {
	db *sql.DB
}

func NewService(db *sql.DB) *usrService {
	return &usrService{
		db: db,
	}
}


func (s *usrService) CreateUser(data UserCreateInput) error {
	// TODO: create errors for users and log errors for developers
	if data.Name == "" {
		return errors.New("user name is required")
	} else if data.Email == "" {
		return errors.New("user email is required")
	} else if data.Password == "" {
		return errors.New("user password is required")
	}

	row := s.db.QueryRow("SELECT COUNT(*) > 0 FROM users WHERE email = ?;", data.Email)
	var exists bool

	if err := row.Scan(&exists); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return errors.New("already exists user with email.")
	}

	_pass, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	password := string(_pass)

	_, err = s.db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?);", data.Name, data.Email, password)

	return err
}
