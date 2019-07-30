package config

// RedisConfig is a redis-related configuration
type RedisConfig struct {
	URL      string `env:"REDIS_URL" env-default:":5050" env-description:"URL of Redis server"`
	Password string `env:"REDIS_PASSWORD" env-description:"Redis password"`
}

// ServerConfig is a server-related configuration
type ServerConfig struct {
	Port string `env:"SERVER_PORT,PORT" env-default:"8080" env-description:"Server port"`
	Host string `env:"SERVER_HOST" env-description:"Server host"`
}

// Config is an application configuration structure
type Config struct {
	Redis  RedisConfig
	Server ServerConfig
}
