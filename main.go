package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ilyakaznacheev/secret/internal/database"
	"github.com/ilyakaznacheev/secret/internal/handler"
	"github.com/ilyakaznacheev/secret/internal/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var conf Config
	cleanenv.ReadEnv(&conf)

	db, err := database.NewRedisDB(conf.Redis.Host+":"+conf.Redis.Port, conf.Redis.Password, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	h := handler.NewSecretHandler(db)

	router := gin.Default()

	v1 := router.Group("/v1")
	v1.POST("/secret", monitoring.MetricsMiddleware(h.PostSecret, "secret_post"))
	v1.GET("/secret/:hash", monitoring.MetricsMiddleware(h.GetSecret, "secret_get"))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run service
	if err := router.Run(fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)); err != nil {
		log.Fatal(err)
	}
}

// RedisConfig is a redis-related configuration
type RedisConfig struct {
	Port     string `env:"REDIS_PORT" env-default:"5050"`
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Password string `env:"REDIS_PASSWORD"`
}

// ServerConfig is a server-related configuration
type ServerConfig struct {
	Port string `env:"SERVER_PORT,PORT" env-default:"8080"`
	Host string `env:"SERVER_HOST"`
}

// Config is an application configuration structure
type Config struct {
	Redis  RedisConfig
	Server ServerConfig
}
