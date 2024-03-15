package types

type Config struct {
	PgPoolInstances []PgPoolInstance `toml:"pg_pool_instance"`
}

type PgPoolInstance struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}
