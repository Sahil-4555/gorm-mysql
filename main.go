package main

import (
	"fmt"
	"log"

	"github.com/Sahil-4555/go-crud-api/configs"
	"github.com/Sahil-4555/go-crud-api/models"
	"github.com/Sahil-4555/go-crud-api/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func migrateTable(v interface{}) error {
	err := DB.AutoMigrate(&v)
	if err != nil {
		log.Println(err)
	}
	return err
}

func main() {
	DB = configs.InitDB()
	fmt.Printf("Connected To Database Sucessfully on Port %s\n", configs.Port())
	err := migrateTable(models.Student{})
	if err != nil {
		panic(err.Error())
	}
	err = migrateTable(models.User{})
	if err != nil {
		panic(err.Error())
	}

	port := "5000"
	router := gin.New()
	router.Use(gin.Logger())
	routes.Routes(router)
	router.Run(":" + port)
}
