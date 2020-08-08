package model

import (
	"fmt"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

type AutoMigrate func(db *gorm.DB)

func New(v *viper.Viper) (db *gorm.DB, err error) {
	config := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		v.GetString("db.mysql.user"),
		v.GetString("db.mysql.pass"),
		v.GetString("db.mysql.host"),
		v.GetInt("db.mysql.port"),
		v.GetString("db.mysql.name"),
	)
	db, err = gorm.Open("mysql", config)
	return db, err
}

var ProviderSet = wire.NewSet(New)
