package controllers

import (
	"github.com/mirfnsyh/base-service/internal/app/commons"
	services "github.com/mirfnsyh/base-service/internal/app/usecases"
)

type ControllerOption struct {
	commons.Options
	*services.Services
}
