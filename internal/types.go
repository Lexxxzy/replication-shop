package types

type Config struct {
	PgPoolInstances       []PgPoolInstance        `toml:"pgpool_instance"`
	RedisSentinelInstance []RedisSentinelInstance `toml:"redis_sentinel_instance"`
	CassandraInstances    []CassandraInstance     `toml:"cassandra_instance"`
}

type PgPoolInstance struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}

type RedisSentinelInstance struct {
	MasterName    string   `toml:"master_name"`
	SentinelAddrs []string `toml:"sentinel_addrs"`
}

type CassandraInstance struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}
