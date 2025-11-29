package pgorm

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/lemonkingstar/spider/pkg/plog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Config struct {
	DBType      string
	User        string
	Password    string
	Host        string
	Port        int
	DBName      string
	TablePrefix string
	LogLevel    string
	Dsn         string // dsn

	MaxIdleConn int
	MaxOpenConn int
	MaxLifetime int
}

func CreateDefault(host string, port int, user, pw, db string) (*gorm.DB, error) {
	cfg := &Config{
		DBType:      "mysql",
		User:        user,
		Password:    pw,
		Host:        host,
		Port:        port,
		DBName:      db,
		TablePrefix: "",
		MaxIdleConn: runtime.NumCPU() * 8,
		MaxOpenConn: runtime.NumCPU() * 8,
		MaxLifetime: 30,
	}
	return cfg.build()
}

func AcquireDefault(dsn string) (*gorm.DB, error) {
	cfg := &Config{
		DBType:      "mysql",
		Dsn:         dsn,
		TablePrefix: "",
		MaxIdleConn: runtime.NumCPU() * 8,
		MaxOpenConn: runtime.NumCPU() * 8,
		MaxLifetime: 30,
	}
	return cfg.build()
}

func New(c *Config) (*gorm.DB, error) {
	return c.build()
}

func (c *Config) build() (*gorm.DB, error) {
	var dial gorm.Dialector
	dsn := c.Dsn
	switch c.DBType {
	case "mysql":
		if dsn == "" {
			dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
				c.User, c.Password, c.Host, c.Port, c.DBName)
		}
		dial = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", c.DBType)
	}

	plog.Infof("Build db connection[ %s ]: %s", c.DBType, dsn)
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
