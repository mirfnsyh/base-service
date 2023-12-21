package commons

import (
	"github.com/mirfnsyh/base-service/config"
	"github.com/mirfnsyh/base-service/internal/app/commons/appcontext"
	"gorm.io/gorm"
)

type Options struct {
	AppCtx   *appcontext.AppContext
	Config   *config.Configuration
	Database *gorm.DB
}
