package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zedisdog/armor/config"
)

var DB *gorm.DB

func Init() {
	var err error
	dbConfig := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.String("db.mysql.user"),
		config.Conf.String("db.mysql.pass"),
		config.Conf.String("db.mysql.host"),
		config.Conf.Int("db.mysql.port"),
		config.Conf.String("db.mysql.name"),
	)
	DB, err = gorm.Open("mysql", dbConfig)
	if err != nil {
		panic(err)
	}
}

func Instance() *gorm.DB {
	if DB == nil {
		Init()
	}
	return DB
}

type AutoMigrate func(db *gorm.DB)
