package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Lexxxzy/go-echo-template/internal"
	"github.com/Lexxxzy/go-echo-template/util"
	"github.com/go-redis/redis/v8"
)

type RedisSentinelManager struct {
	clients []*redis.Client
	configs []types.RedisSentinelInstance
	index   int
}

var (
	ctx = context.Background()
)

var RedisSentinelProxy *RedisSentinelManager

func NewRedisSentinelManager(configs []types.RedisSentinelInstance) *RedisSentinelManager {
	manager := &RedisSentinelManager{
		configs: configs,
		index:   0,
	}
	manager.clients = make([]*redis.Client, len(configs))
	for i, config := range configs {
		manager.connect(i, config, 0)
	}
	return manager
}

func (manager *RedisSentinelManager) connect(index int, config types.RedisSentinelInstance, attempt int) {
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
		Username:      os.Getenv("REDIS_USERNAME"),
		Password:      os.Getenv("REDIS_PASSWORD"),
		DialTimeout:   5 * time.Second,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Printf("Failed to connect to Redis Sentinel at %v, error: %v\n", config.SentinelAddrs, err)
		delay := time.Minute
		if attempt == 1 {
			delay = time.Hour
		}
		time.AfterFunc(delay, func() {
			manager.connect(index, config, attempt+1)
		})
		return
	} else {
		log.Printf("Connected to Redis Sentinel at %v\n", config.SentinelAddrs)
	}

	manager.clients[index] = rdb
}

func (manager *RedisSentinelManager) GetCurrentClient() *redis.Client {
	for i := 0; i < len(manager.clients); i++ {
		idx := (manager.index + i) % len(manager.clients)
		if manager.clients[idx] != nil {
			log.Printf("INFO: Using Redis Sentinel at %v\n", manager.configs[idx].SentinelAddrs)
			manager.index = (idx + 1) % len(manager.clients)
			return manager.clients[idx]
		}
	}
	log.Fatal("No Redis Sentinel connection is available")
	return nil
}

func InitRedisSentinel(configPath string) error {
	config, err := util.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}
	RedisSentinelProxy = NewRedisSentinelManager(config.RedisSentinelInstance)
	return nil
}
