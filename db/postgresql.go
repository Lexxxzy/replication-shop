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

type PostgresqlManager struct {
	instances []*bun.DB
	configs   []types.PgPoolInstance
	index     int
}

var PostgresqlProxy *PostgresqlManager

func NewPostgresqlManager(configs []types.PgPoolInstance) *PostgresqlManager {
	manager := &PostgresqlManager{
		configs: configs,
		index:   0,
	}
	manager.instances = make([]*bun.DB, len(configs))
	for i, config := range configs {
		manager.connect(i, config, 0)
	}
	return manager
}

func (manager *PostgresqlManager) connect(index int, config types.PgPoolInstance, attempt int) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), config.IP, config.Port, os.Getenv("POSTGRES_DB"))
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	bunDB := bun.NewDB(db, pgdialect.New())
	if err := db.Ping(); err != nil {
		log.Printf("Failed to connect to database instance at %s:%d, error: %v\n", config.IP, config.Port, err)
		delay := time.Second * 10
		if attempt == 1 {
			delay *= 2
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

func (manager *PostgresqlManager) GetCurrentDB() *bun.DB {
	for i := 0; i < len(manager.instances); i++ {
		idx := (manager.index + i) % len(manager.instances)
		if manager.instances[idx] != nil {
			log.Printf("INFO: Using database instance at %s:%d\n", manager.configs[idx].IP, manager.configs[idx].Port)
			manager.index = (idx + 1) % len(manager.instances)
			return manager.instances[idx]
		}
	}
	log.Fatal("No database connection is available")
	return nil
}

func InitPostgresql(configPath string) error {
	config, err := util.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}
	PostgresqlProxy = NewPostgresqlManager(config.PgPoolInstances)

	return nil
}
