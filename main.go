package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ilyakaznacheev/secret/database"
	"github.com/ilyakaznacheev/secret/handler"
)

func main() {

	var conf RedisConfig
	cleanenv.ReadEnv(&conf)

	db, err := database.NewRedisDB(conf.Host+":"+conf.Port, conf.Password, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	h := handler.NewSecretHandler(db)

	router := gin.Default()
	router.POST("/secret", h.PostSecret)
	router.GET("/secret/:hash", h.GetSecret)

	// Run service
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

type RedisConfig struct {
	Port     string `env:"REDIS_PORT" env-default:"5050"`
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Password string `env:"REDIS_PASSWORD"`
}

type Config struct {
	Redis RedisConfig
}
