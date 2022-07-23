package database

import (
	"fmt"
	gormlog "gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	LogLevel string // silent, error, warn, info
}

func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName)
}

func GormDB(config *MySQLConfig) (*gorm.DB, error) {
	gormLogger := gormlog.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlog.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  convert2GormLogLevel(config.LogLevel),
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})
	db, err := gorm.Open(mysql.Open(config.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mysql")
	}
	return db, nil
}

func convert2GormLogLevel(level string) gormlog.LogLevel {
	switch level {
	case "silent":
		return gormlog.Silent
	case "error":
		return gormlog.Error
	case "warn":
		return gormlog.Warn
	case "info":
		return gormlog.Info
	default:
		return gormlog.Silent
	}
}
