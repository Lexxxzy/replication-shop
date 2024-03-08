package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"os"
)

var Bun *bun.DB

func Setup() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		dbname   = os.Getenv("POSTGRES_DB")
		dbuser   = os.Getenv("POSTGRES_USER")
		dbpasswd = os.Getenv("POSTGRES_PASSWORD")
		dbhost   = os.Getenv("DB_HOST")
		dbport   = os.Getenv("DB_PORT")
		dsn      = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbuser, dbpasswd, dbhost, dbport, dbname)
	)

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Init initializes the database connection.
//
// It sets up the database connection using the Setup function and
// adds a query hook using the bundebug package. The function returns
// an error if there was an issue connecting to the database.
//
// Parameters:
//
//	None
//
// Return:
//
//	error: An error if there was an issue connecting to the database.
func Init() error {
	db, err := Setup()
	if err != nil {
		log.Fatal("Error connecting to database. " + err.Error())
		return err
	}

	Bun = bun.NewDB(db, pgdialect.New())
	Bun.AddQueryHook(bundebug.NewQueryHook())

	return nil
}
