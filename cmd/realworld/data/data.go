package data

import (
	"github.com/lemonkingstar/spider/pkg/predis"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb predis.Client
)

func Startup() error {
	return nil
}

type Data struct{}

func (d *Data) db() *gorm.DB       { return db }
func (d *Data) rdb() predis.Client { return rdb }
