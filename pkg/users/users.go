package users

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserCreateInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Inactive bool      `json:"inactive"`
	CreatdAt time.Time `json:"createdAt"`
}

type userService struct {
	db *sql.DB
}

func NewService(db *sql.DB) *userService {
	return &userService{
		db: db,
	}
}

func (us *userService) CreateUser(data UserCreateInput) error {
	// TODO: create errors for users and log internal errors
	if data.Name == "" {
		return errors.New("user name is required")
	} else if data.Email == "" {
		return errors.New("user email is required")
	} else if data.Password == "" {
		return errors.New("user password is required")
	}

	row := us.db.QueryRow("SELECT COUNT(*) > 0 FROM users WHERE email = ?;", data.Email)
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

	_, err = us.db.Exec("INSERT INTO users (name, email, password, inactive) VALUES (?, ?, ?, ?);", data.Name, data.Email, password, false)

	return err
}

func (us *userService) FindById(id int64) (*User, error) {
	row := us.db.QueryRow("SELECT name, email, password, inactive, created_at FROM users where id = ?;", id)
	var name string
	var email string
	var password string
	var inactive sql.NullBool
	var createdAt time.Time
	if err := row.Scan(&name, &email, &password, &inactive, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: customise and handle this error
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &User{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Inactive: inactive.Bool,
		CreatdAt: createdAt,
	}, nil
}

func (us *userService) FindByEmail(email string) (*User, error) {
	row := us.db.QueryRow("SELECT name, email, password, inactive, created_at FROM users where email = ?;", email)
	var id int64
	var name string
	var password string
	var inactive sql.NullBool
	var createdAt time.Time
	if err := row.Scan(&name, &email, &password, &inactive, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: customise and handle this error
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &User{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Inactive: inactive.Bool,
		CreatdAt: createdAt,
	}, nil
}

func (us *userService) Exists(ids []int64) (bool, error) {
	_ids := make([]any, len(ids))
	for i, id := range ids {
		_ids[i] = id
	}

	args := strings.Repeat("?, ", len(_ids))
	args = args[:len(args)-2]

	query := "SELECT id FROM users where id in (" + args + ");"
	fmt.Println(query)
	rows, err := us.db.Query(query, _ids...)

	if err != nil {
		return false, err
	}

	found := make([]int64, 0)
	for rows.Next() {
		var uid int64
		if err := rows.Scan(&uid); err != nil {
			return false, err
		}
		found = append(found, uid)
	}

	if len(ids) != len(found) {
		return false, nil
	}

	return true, nil
}
