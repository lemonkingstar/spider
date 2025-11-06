package pgorm

import "gorm.io/gorm"

func Migrate(db *gorm.DB, dst ...interface{}) error {
	return db.Migrator().AutoMigrate(dst...)
}
