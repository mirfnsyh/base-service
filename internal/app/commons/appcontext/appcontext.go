package appcontext

import (
	"errors"

	"github.com/mirfnsyh/base-service/config"
	"github.com/mirfnsyh/base-service/internal/pkg/driver"
	"gorm.io/gorm"
)

const (
	DBDialectMysql = "mysql"
)

type AppContext struct {
	config *config.Configuration
}

func NewAppContext(config *config.Configuration) *AppContext {
	return &AppContext{
		config: config,
	}
}

func (ctx *AppContext) GetAppOption() AppOption {
	return AppOption{
		Host: ctx.config.App.Host,
		Port: ctx.config.App.Port,
		Env:  ctx.config.App.Env,
	}
}

func (ctx *AppContext) GetDBInstance(dbType string) (*gorm.DB, error) {
	var gorm *gorm.DB
	var err error

	switch dbType {
	case DBDialectMysql:
		dbOption := ctx.GetMysqlOption()
		gorm, err = driver.NewMysqlDatabase(dbOption)
	default:
		err = errors.New("error get db instance, unknown db type")
	}

	return gorm, err
}

func (ctx *AppContext) GetMysqlOption() driver.DBMysqlOption {
	return driver.DBMysqlOption{
		Host:                 ctx.config.Database.Host,
		Port:                 ctx.config.Database.Port,
		User:                 ctx.config.Database.User,
		Password:             ctx.config.Database.Password,
		DBName:               ctx.config.Database.Name,
		AdditionalParameters: ctx.config.Database.AdditionalParameters,
		MaxOpenCon:           ctx.config.Database.MaxOpenConns,
		MaxIdleCon:           ctx.config.Database.MaxIdleConns,
		ConnMaxLifetime:      ctx.config.Database.ConnMaxLifetime,
		Debug:                ctx.config.Database.Debug,
	}
}
