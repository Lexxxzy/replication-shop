package data

import (
	"context"
	"net/http"
	"time"

	"github.com/Lexxxzy/go-echo-template/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id"`
	Name      string    `bun:"type:char(64),notnull" json:"name"`
	Email     string    `bun:"type:char(64),unique,notnull" json:"email"`
	Password  string    `bun:"type:text,notnull" json:"password"`
	CreatedAt time.Time `bun:"type:timestamptz,default:current_timestamp,notnull" json:"created_at"`
}

func CreateUser(user *User, c echo.Context) error {
	query := `
		INSERT INTO users (name, email, password) VALUES (?, ?, ?)
		RETURNING id, created_at
	`
	_, err := db.Proxy.GetCurrentDB().NewRaw(query, user.Name, user.Email, user.Password).Exec(c.Request().Context())
	if err != nil {
		log.Error("Error creating user. ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return nil
}

func GetUser(user *User) error {
	query := "SELECT * FROM users WHERE id = ?"
	return db.Proxy.GetCurrentDB().NewRaw(query, user.ID).Scan(context.Background(), user)
}

func GetUserById[T string | uuid.UUID](id T) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = ?"
	err := db.Proxy.GetCurrentDB().NewRaw(query, id).Scan(context.Background(), &user)
	return user, err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE email = ?"
	err := db.Proxy.GetCurrentDB().NewRaw(query, email).Scan(context.Background(), &user)
	return user, err
}

func IsUserExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	err := db.Proxy.GetCurrentDB().NewRaw(query, email).Scan(context.Background(), &exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
