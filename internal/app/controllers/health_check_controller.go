package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mirfnsyh/base-service/internal/app/commons"
	httpHelper "github.com/mirfnsyh/base-service/internal/app/commons/helper/http"
)

type HealthCheckController struct {
	ControllerOption
	httpHelper.HttpHandlerContext
}

func NewHealthCheckController(opt ControllerOption, handlerCtx *httpHelper.HttpHandlerContext) *HealthCheckController {
	return &HealthCheckController{
		ControllerOption:   opt,
		HttpHandlerContext: *handlerCtx,
	}
}

func (ctl *HealthCheckController) Ping(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			commons.PrintError(err, c, &ctl.HttpHandlerContext)
			return
		}
	}()

	result := ctl.Services.HealthCheckService.Ping()

	ctl.WriteData(c, result)
}
