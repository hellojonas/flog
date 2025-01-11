package apps

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

const APP_KEY_LEN = 32

type App struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	UserId    int64     `json:"userId"`
	Inactive  bool      `json:"inactive"`
	CreatedAt time.Time `json:"createdAt"`
}

type AppCreateInput struct {
	Name   string `json:"name"`
	UserId int64  `json:"userId"`
}

type AppMemberInput struct {
	Members []int64 `json:"members"`
}

type AppService struct {
	db *sql.DB
}

func NewService(db *sql.DB) *AppService {
	return &AppService{
		db: db,
	}
}

func (as *AppService) FindById(id int64) (*App, error) {
	row := as.db.QueryRow("SELECT name, token, inactive, user_id, created_at FROM applications WHERE id = ?", id)

	var name string
	var token string
	var inactive bool
	var userId int64
	var createdAt time.Time

	if err := row.Scan(&name, &token, &inactive, &userId, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: customise and handle this error
			return nil, errors.New("application not found")
		}
		return nil, err
	}

	return &App{
		Id:        id,
		Name:      name,
		Token:     token,
		Inactive:  inactive,
		UserId:    userId,
		CreatedAt: createdAt,
	}, nil
}

func (as *AppService) FindByName(appName string) (*App, error) {
	row := as.db.QueryRow("SELECT id, token, inactive, user_id, created_at FROM applications WHERE name = ?", appName)

	var id int64
	var token string
	var inactive bool
	var userId int64
	var createdAt time.Time

	if err := row.Scan(&id, &token, &inactive, &userId, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: customise and handle this error
			return nil, errors.New("application not found")
		}
		return nil, err
	}

	return &App{
		Id:        id,
		Name:      appName,
		Token:     token,
		Inactive:  inactive,
		UserId:    userId,
		CreatedAt: createdAt,
	}, nil
}

func (as *AppService) CreateApp(data AppCreateInput) (*App, error) {
	token, err := genKey(flog.APP_KEY_LEN)
	name := strings.ReplaceAll(data.Name, " ", "_")
	res, err := as.db.Exec("INSERT INTO applications (name, token, inactive, user_id) VALUES (?, ?, ?, ?) RETURNING id;", name, token, false, data.UserId)

	if err != nil {
		return nil, err
	}

	appId, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	if err != as.SetMembers(appId, []int64{data.UserId}) {
		return nil, err
	}

	return as.FindById(appId)
}

func (as *AppService) SetMembers(app int64, members []int64) error {
	appMember := make([]string, len(members))
	for i := range members {
		appMember[i] = "(?, ?)"
	}
	values := strings.Join(appMember, ",")

	args := make([]any, len(members)*2)
	for i := 0; i < len(members); i += 2 {
		args[i] = members[i]
		args[i+1] = app
	}

	query := "INSERT INTO user_applications (user_id, application_id) values " + values + ";"

	_, err := as.db.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (as *AppService) ListUserApps(userId int64) ([]App, error) {
	query := ` SELECT a.id, a.name, a.token, a.inactive, a.created_at
    FROM user_applications ua
    INNER JOIN users u on u.id = ua.user_id
    INNER JOIN applications a on a.id = ua.application_id
    WHERE ua.user_id = ?;
    `

	rows, err := as.db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	var id int64
	var name string
	var token string
	var inactive bool
	var createdAt time.Time

	var apps []App
	for rows.Next() {
		err := rows.Scan(&id, &name, &token, &inactive, &createdAt)
		if err != nil {
			return nil, err
		}
		apps = append(apps, App{
			Id:        id,
			Name:      name,
			Token:     token,
			Inactive:  inactive,
			CreatedAt: createdAt,
		})
	}

	return apps, nil
}

func (as *AppService) ListAppMembers(appId int64) ([]App, error) {
	query := ` SELECT u.id, u.name, u.email, u.inactive, u.created_at
    FROM user_applications ua
    INNER JOIN users u on u.id = ua.user_id
    INNER JOIN applications a on a.id = ua.application_id
    WHERE ua.application_id = ?;
    `

	rows, err := as.db.Query(query, appId)

	if err != nil {
		return nil, err
	}

	var id int64
	var name string
	var token string
	var inactive sql.NullBool
	var createdAt time.Time

	var users []App
	for rows.Next() {
		err := rows.Scan(&id, &name, &token, &inactive, &createdAt)
		if err != nil {
			return nil, err
		}
		users = append(users, App{
			Id:        id,
			Name:      name,
			Token:     token,
			Inactive:  inactive.Bool,
			CreatedAt: createdAt,
		})
	}

	return users, nil
}

func genKey(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(b)

	return key, nil
}
