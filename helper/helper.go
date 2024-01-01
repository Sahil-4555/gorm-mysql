package helper

import (
	"github.com/Sahil-4555/go-crud-api/configs"
	"github.com/Sahil-4555/go-crud-api/models"
	"gorm.io/gorm"
)

var DB *gorm.DB = configs.InitDB()

func GetCount() (int64, error) {
	var counter int64
	err := DB.Model(&models.Student{}).Where("deleteat IS NULL").Count(&counter).Error
	if err != nil {
		return counter, err
	}
	return counter, nil
}

func GetData(key, limit int) ([]models.Student, error) {
	var result []models.Student
	err := DB.Model(models.Student{}).Where("deleteat IS NULL").Limit(limit).Offset((key - 1) * int(limit)).Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}
