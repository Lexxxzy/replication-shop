package db

import (
	"fmt"
	"log"
	"time"

	"github.com/Lexxxzy/go-echo-template/internal"
	"github.com/Lexxxzy/go-echo-template/util"
	"github.com/gocql/gocql"
)

type CassandraManager struct {
	instances []*gocql.Session
	configs   []types.CassandraInstance
	index     int
}

var CassandraProxy *CassandraManager

func NewCassandraManager(configs []types.CassandraInstance) *CassandraManager {
	manager := &CassandraManager{
		configs: configs,
		index:   0,
	}
	manager.instances = make([]*gocql.Session, len(configs))
	for i, config := range configs {
		manager.connect(i, config, 0)
	}
	return manager
}

func (manager *CassandraManager) connect(index int, config types.CassandraInstance, attempt int) {
	cluster := gocql.NewCluster(config.IP)
	cluster.Port = config.Port
	cluster.ProtoVersion = 4

	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("Failed to connect to Cassandra instance at %s:%d, error: %v\n", config.IP, config.Port, err)
		delay := time.Minute
		if attempt == 1 {
			delay = time.Hour
		}
		time.AfterFunc(delay, func() {
			manager.connect(index, config, attempt+1)
		})
		return
	} else {
		log.Printf("Connected to Cassandra instance at %s:%d\n", config.IP, config.Port)
	}

	manager.instances[index] = session
}

func (manager *CassandraManager) GetCurrentSession() *gocql.Session {
	for i := 0; i < len(manager.instances); i++ {
		idx := (manager.index + i) % len(manager.instances)
		if manager.instances[idx] != nil {
			log.Printf("INFO: Using Cassandra instance at %s:%d\n", manager.configs[idx].IP, manager.configs[idx].Port)
			manager.index = (idx + 1) % len(manager.instances)
			return manager.instances[idx]
		}
	}
	log.Fatal("No Cassandra connection is available")
	return nil
}

func InitCassandra(configPath string) error {
	config, err := util.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading Cassandra configuration: %v", err)
	}
	CassandraProxy = NewCassandraManager(config.CassandraInstances)

	return nil
}
