package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mirfnsyh/base-service/config"
	httpHelper "github.com/mirfnsyh/base-service/internal/app/commons/helper/http"
	"github.com/mirfnsyh/base-service/internal/app/commons/helper/structs"
	"github.com/mirfnsyh/base-service/internal/app/controllers"
	"github.com/mirfnsyh/base-service/internal/app/server/http/middlewares"
)

func Router(opt controllers.ControllerOption) *gin.Engine {
	var r *gin.Engine

	cfg := config.GetConfig()
	basePath := "base/"

	handlerCtx := httpHelper.NewContextHandler(structs.Meta{
		APIEnv:  cfg.App.Env,
		Version: cfg.App.Version,
	})

	r = gin.Default()

	r.Use(middlewares.API())
	r.Use(middlewares.CORS())

	healthCheckController := controllers.NewHealthCheckController(opt, &handlerCtx)
	r.GET(fmt.Sprintf("%s/ping", basePath), healthCheckController.Ping)

	return r
}
