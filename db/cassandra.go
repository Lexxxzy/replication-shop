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
	session *gocql.Session
}

var CassandraProxy *CassandraManager

func NewCassandraManager(configs []types.CassandraInstance) *CassandraManager {
	manager := &CassandraManager{}
	manager.connect(configs, 0)
	return manager
}

func (manager *CassandraManager) connect(configs []types.CassandraInstance, attempt int) {
	cluster := gocql.NewCluster()
	for _, config := range configs {
		cluster.Hosts = append(cluster.Hosts, fmt.Sprintf("%s:%d", config.IP, config.Port))
	}
	cluster.DisableInitialHostLookup = true
	cluster.ProtoVersion = 4
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("Failed to connect to Cassandra instances, error: %v\n", err)
		delay := time.Second * 2
		if attempt == 1 {
			delay *= 2
		}
		time.AfterFunc(delay, func() {
			manager.connect(configs, attempt+1)
		})
		return
	} else {
		log.Println("Connected to Cassandra pool")
	}
	session.SetConsistency(gocql.One)
	manager.session = session
}

func (manager *CassandraManager) GetCurrentSession() *gocql.Session {
	log.Println("INFO: Using cassandra instance")
	return manager.session
}

func InitCassandra(configPath string) error {
	config, err := util.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading Cassandra configuration: %v", err)
	}
	CassandraProxy = NewCassandraManager(config.CassandraInstances)
	return nil
}
