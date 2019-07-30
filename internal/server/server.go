package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/secret/internal/config"
	"github.com/ilyakaznacheev/secret/internal/database"
	"github.com/ilyakaznacheev/secret/internal/handler"
	"github.com/ilyakaznacheev/secret/internal/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Run start the server
func Run(conf config.Config) error {
	db, err := database.NewRedisDB(conf.Redis.URL, conf.Redis.Password, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	h := handler.NewSecretHandler(db)

	router := gin.Default()

	router.GET("/", h.Redirect)

	v1 := router.Group("/v1")
	v1.POST("/secret", monitoring.MetricsMiddleware(h.PostSecret, "secret_post"))
	v1.GET("/secret/:hash", monitoring.MetricsMiddleware(h.GetSecret, "secret_get"))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run service
	return router.Run(fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port))
}
