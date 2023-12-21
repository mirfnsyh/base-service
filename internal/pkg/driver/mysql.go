package driver

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DBMysqlOption struct {
	Host                 string
	Port                 int
	User                 string
	Password             string
	DBName               string
	AdditionalParameters string
	MaxOpenCon           int
	MaxIdleCon           int
	ConnMaxLifetime      time.Duration
	Debug                bool
}

func NewMysqlDatabase(option DBMysqlOption) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		option.User, option.Password, option.Host, option.Port, option.DBName, option.AdditionalParameters))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(option.MaxOpenCon)
	db.SetMaxIdleConns(option.MaxIdleCon)
	db.SetConnMaxLifetime(option.ConnMaxLifetime)

	if option.Debug {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	gormInitialize, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), gormConfig)

	if err != nil {
		return nil, err
	}

	return gormInitialize, nil
}
