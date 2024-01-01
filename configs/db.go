package configs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	User := Username()
	Password := Password()
	Host := Host()
	Port := Port()
	DbName := DbName()
	createDBDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", User, Password, Host, Port)
	database, _ := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})
	_ = database.Exec("CREATE DATABASE IF NOT EXISTS " + DbName + ";")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", User, Password, Host, Port, DbName)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return DB
}
