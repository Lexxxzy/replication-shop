package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/Lexxxzy/go-echo-template/internal"
	"github.com/Lexxxzy/go-echo-template/util"
)

type DBManager struct {
	instances []*bun.DB
	configs   []types.PgPoolInstance
	index     int
}

var Proxy *DBManager

func NewDBManager(configs []types.PgPoolInstance) *DBManager {
	manager := &DBManager{
		configs: configs,
		index:   0,
	}
	manager.instances = make([]*bun.DB, len(configs))
	for i, config := range configs {
		manager.connect(i, config, 0)
	}
	return manager
}

func (manager *DBManager) connect(index int, config types.PgPoolInstance, attempt int) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), config.IP, config.Port, os.Getenv("POSTGRES_DB"))
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	bunDB := bun.NewDB(db, pgdialect.New())
	if err := db.Ping(); err != nil {
		log.Printf("Failed to connect to database instance at %s:%d, error: %v\n", config.IP, config.Port, err)
		delay := time.Minute
		if attempt == 1 {
			delay = time.Hour
		}
		time.AfterFunc(delay, func() {
			manager.connect(index, config, attempt+1)
		})
		return
	} else {
		log.Printf("Connected to database instance at %s:%d\n", config.IP, config.Port)
	}

	manager.instances[index] = bunDB
}

func (manager *DBManager) GetCurrentDB() *bun.DB {
	for i := 0; i < len(manager.instances); i++ {
		idx := (manager.index + i) % len(manager.instances)
		if manager.instances[idx] != nil {
			manager.index = (idx + 1) % len(manager.instances)
			return manager.instances[idx]
		}
	}
	log.Fatal("No database connection is available")
	return nil
}

func Init(configPath string) error {
	config, err := util.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}
	Proxy = NewDBManager(config.PgPoolInstances)

	return nil
}
