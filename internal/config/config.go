package config

// RedisConfig is a redis-related configuration
type RedisConfig struct {
	URL string `env:"REDIS_URL" env-default:":5050" env-description:"URL of Redis server including options"`
}

// ServerConfig is a server-related configuration
type ServerConfig struct {
	Port string `env:"SERVER_PORT,PORT" env-default:"8080" env-description:"Server port"`
	Host string `env:"SERVER_HOST" env-description:"Server host"`
}

// RedirectConfig contains redirection settings
type RedirectConfig struct {
	Root string `env:"REDIRECT_ROOT" env-default:"http://bit.ly/2MrGGf8" env-description:"Redirect from root page"`
	API  string `env:"REDIRECT_API" env-default:"http://bit.ly/2MrGGf8" env-description:"Redirect to API documentation"`
}

// Config is an application configuration structure
type Config struct {
	Redis    RedisConfig
	Server   ServerConfig
	Redirect RedirectConfig
}
