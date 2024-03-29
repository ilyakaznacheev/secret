package secret

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/secret/internal/config"
	"github.com/ilyakaznacheev/secret/internal/database"
	"github.com/ilyakaznacheev/secret/internal/handler"
	"github.com/ilyakaznacheev/secret/internal/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Run start the server
func Run(conf config.Config) error {
	var (
		db  *database.RedisDB
		err error
	)
	if conf.Redis.URL != "" {
		db, err = database.NewRedisDBWithOpts(conf.Redis.URL)
	} else {
		db, err = database.NewRedisDB(conf.Redis.Host + ":" + conf.Redis.Port)
	}
	if err != nil {
		return err
	}

	h := handler.NewSecretHandler(db)

	router := gin.Default()

	router.GET("/", handler.RedirectTo(conf.Redirect.Root))

	v1 := router.Group("/v1")
	v1.POST("/secret", monitoring.MetricsMiddleware(h.PostSecret, "secret_post"))
	v1.GET("/secret/:hash", monitoring.MetricsMiddleware(h.GetSecret, "secret_get"))
	v1.GET("/", handler.RedirectTo(conf.Redirect.API))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run service
	return router.Run(fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port))
}
