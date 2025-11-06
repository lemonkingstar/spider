package pgorm

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/lemonkingstar/spider/pkg/plog"
)

type DBConfig struct {
	DBType      string
	User        string
	Password    string
	Host        string
	Port        int
	DBName      string
	TablePrefix string
	LogLevel    string

	MaxIdleConn int
	MaxOpenConn int
	MaxLifetime int
}

func Default(host string, port int, user, password, db string) (*gorm.DB, error) {
	cfg := &DBConfig{
		DBType:      "mysql",
		User:        user,
		Password:    password,
		Host:        host,
		Port:        port,
		DBName:      db,
		TablePrefix: "",
		MaxIdleConn: runtime.NumCPU() * 8,
		MaxOpenConn: runtime.NumCPU() * 8,
		MaxLifetime: 30,
		LogLevel:    "DEBUG",
	}
	return cfg.build()
}

func New(c *DBConfig) (*gorm.DB, error) {
	return c.build()
}

func (c *DBConfig) build() (*gorm.DB, error) {
	var dial gorm.Dialector
	switch c.DBType {
	case "mysql":
		dial = mysql.Open(fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.DBName))
	default:
		return nil, fmt.Errorf("unsupported db type: %s", c.DBType)
	}

	plog.Infof("Build db connection[ %s ]: %s:%d/%s", c.DBType, c.Host, c.Port, c.DBName)
	logLevel := logger.Info
	if strings.ToUpper(c.LogLevel) == "WARN" {
		logLevel = logger.Warn
	} else if strings.ToUpper(c.LogLevel) == "ERROR" {
		logLevel = logger.Error
	}
	db, err := gorm.Open(dial, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	return db, nil
}
