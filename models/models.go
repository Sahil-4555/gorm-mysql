package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	Id         int       `json:"id,omitempty" gorm:"primary_key"`
	Name       string    `json:"name" binding:"required,min=2,max=49" gorm:"type:varchar(255);"`
	Age        int       `json:"age" binding:"required,gt=0"`
	Dob        string    `json:"dob" binding:"required" gorm:"type:varchar(255);"`
	Image      string    `json:"image" gorm:"type:varchar(255);"`
	Cgpa       float64   `json:"cgpa" binding:"required,gt=0"`
	Currentsem int       `json:"currentsem" binding:"required,gt=0"`
	Createdat  string    `json:"createdat" gorm:"type:varchar(255);"`
	Updatedat  time.Time `json:"updatedat" gorm:"autoUpdateTime"`
	Deleteat   *string   `json:"deletedat" gorm:"type:varchar(255);"`
}

type StudentUpadte struct {
	Id         int       `json:"id" gorm:"primary_key"`
	Name       *string   `json:"name" binding:"omitempty,min=2,max=49"`
	Age        *int      `json:"age" binding:"omitempty,gt=0"`
	Dob        *string   `json:"dob" bidning:"onitempty"`
	Image      *string   `json:"image" gorm:"type:varchar(255);"`
	Cgpa       *float64  `json:"cgpa" binding:"omitempty,gt=0"`
	Currentsem *int      `json:"currentsem" binding:"omitempty,gt=0"`
	Createdat  string    `json:"createdat" gorm:"type:varchar(255);"`
	Updatedat  time.Time `json:"updatedat" gorm:"autoUpdateTime"`
	Deleteat   *string   `json:"deletedat" gorm:"type:varchar(255);"`
}

type Pagination struct {
	Thispage   int       `json:"thispage"`
	Page       int       `json:"page"`
	Data       []Student `json:"data"`
	Totalpages int       `json:"totalpages"`
	Totalcount int       `json:"totalcount"`
}

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

type Paginationparam struct {
	Page   int `json:"page"`
	Offset int `json:"offset"`
}
