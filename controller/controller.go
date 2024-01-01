package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/Sahil-4555/go-crud-api/configs"
	"github.com/Sahil-4555/go-crud-api/helper"
	"github.com/Sahil-4555/go-crud-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IError struct {
	Field string
	Tag   string
	Value string
}

var DB *gorm.DB = configs.InitDB()

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "max":
		return "Should be less than " + fe.Param()
	case "min":
		return "Should be greater than " + fe.Param()
	case "lte":
		return "Should be less than " + fe.Param()
	case "gt":
		return "Should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func Validatedate(date string) (string, string) {
	re := regexp.MustCompile("^(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-((19|20)\\d\\d)")
	if !re.MatchString(date) {
		return "Dob", "INVALID DATE"
	}
	return "", ""
}

func Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var data models.Student
	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	data.Createdat = now.Format("02-01-2006 15:04:05")

	defer cancel()
	var out []ErrorMsg
	if err := c.ShouldBind(&data); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				out = append(out, ErrorMsg{fe.Field(), getErrorMsg(fe)})
			}
		}
	}
	field, err := Validatedate(data.Dob)
	if err != "" {
		out = append(out, ErrorMsg{field, err})
	}
	if len(out) != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errors": out,
		})
		return
	}
	DB.Model(models.Student{}).WithContext(ctx).Create(&data)

	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

const folder = "images"

func UploadImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	reqParamsId := c.Param("id")
	idtodo, _ := strconv.Atoi(reqParamsId)
	defer cancel()
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var datai models.Student
	datatoUpdate := DB.Model(models.Student{}).WithContext(ctx).Where("id = ? AND deleteat IS NULL", idtodo).First(&datai)
	if datatoUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not Found",
		})
		return
	}
	ext1 := filepath.Ext(handler.Filename)
	datai.Image = (datai.Createdat + ext1)
	result := DB.WithContext(ctx).Save(&datai)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to Update the data",
		})
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fullFolderPath := filepath.Join(wd, folder)
	if _, err := os.Stat(fullFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullFolderPath, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	newFilePath := filepath.Join(fullFolderPath, (datai.Createdat)+ext1)

	newFile, err := os.Create(newFilePath)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	defer newFile.Close()
	io.Copy(newFile, file)
	c.JSON(http.StatusCreated, gin.H{
		"data": result.RowsAffected,
	})
}

func UpdateImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	reqParamsId := c.Param("id")
	idtodo, _ := strconv.Atoi(reqParamsId)
	defer cancel()
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var datai models.Student
	datatoUpdate := DB.Model(models.Student{}).WithContext(ctx).Where("id = ? AND deleteat IS NULL", idtodo).First(&datai)
	if datatoUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not Found",
		})
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fullFolderPath := filepath.Join(wd, folder)
	if _, err := os.Stat(fullFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullFolderPath, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	os.Remove(filepath.Join(fullFolderPath, datai.Image))
	ext1 := filepath.Ext(handler.Filename)
	datai.Image = (datai.Createdat + ext1)
	result := DB.WithContext(ctx).Save(&datai)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to Update the data",
		})
		return
	}

	newFilePath := filepath.Join(fullFolderPath, (datai.Createdat)+ext1)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	defer newFile.Close()
	io.Copy(newFile, file)
	c.JSON(http.StatusCreated, gin.H{
		"data": result.RowsAffected,
	})
}

func GetImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	reqParamsId := c.Param("id")
	idtodo, _ := strconv.Atoi(reqParamsId)
	defer cancel()
	var datai models.Student
	datatoUpdate := DB.Model(models.Student{}).WithContext(ctx).Where("id = ? AND deleteat IS NULL", idtodo).First(&datai)
	if datatoUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not Found",
		})
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fullFolderPath := filepath.Join(wd, folder)
	if _, err := os.Stat(fullFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullFolderPath, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	f, err := os.Open(filepath.Join(fullFolderPath, datai.Image))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to Open the File",
		})
		return
	}

	ext := filepath.Ext(datai.Image)
	var extArr = []string{".png", ".jpeg", ".jpg", ".gif"}
	for i := 0; i <= len(extArr); i++ {
		if ext == extArr[i] {
			c.Set("Content-Type", "image/"+ext[1:])
			img, _, err := image.Decode(f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Unable to Decode the %s file", ext),
				})
				return
			}
			buffer := new(bytes.Buffer)
			if ext == ".png" {
				err = png.Encode(buffer, img)
			} else if ext == ".jpg" || ext == ".jpeg" {
				err = jpeg.Encode(buffer, img, nil)
			} else if ext == ".gif" {
				err = gif.Encode(buffer, img, nil)
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Unable to Encode the %s file", ext),
				})
				return
			}
			c.Data(http.StatusOK, "image/"+ext[1:], buffer.Bytes())
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Unsupported Image Format",
	})
}

func Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var data models.StudentUpadte
	reqParamsId := c.Param("id")
	idtodo, _ := strconv.Atoi(reqParamsId)
	defer cancel()
	var out []ErrorMsg
	if err := c.ShouldBindJSON(&data); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				out = append(out, ErrorMsg{fe.Field(), getErrorMsg(fe)})
			}
		}
	}
	field, err := Validatedate(*data.Dob)
	if err != "" {
		out = append(out, ErrorMsg{field, err})
	}
	if len(out) != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errors": out,
		})
		return
	}

	var datai models.Student
	datatoUpdate := DB.WithContext(ctx).Where("id = ? AND deleteat IS NULL", idtodo).First(&datai)
	if datatoUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not Found",
		})
		return
	}

	if data.Age != nil {
		datai.Age = *data.Age
	}
	if data.Name != nil {
		datai.Name = *data.Name
	}
	if data.Dob != nil {
		datai.Dob = *data.Dob
	}
	if data.Cgpa != nil {
		datai.Cgpa = *data.Cgpa
	}
	if data.Currentsem != nil {
		datai.Currentsem = *data.Currentsem
	}
	result := DB.WithContext(ctx).Save(&datai)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to Update the data",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": result.RowsAffected,
	})
}

func Getall(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var students []models.Student

	err := DB.Model(models.Student{}).WithContext(ctx).Where("deleteat IS NULL").Find(&students)
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error Getting Data",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": students,
	})
}

func Getbyid(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var data models.Student
	reqParamsId := c.Param("id")
	idtodo, _ := strconv.Atoi(reqParamsId)
	defer cancel()

	datatoUpdate := DB.Model(models.Student{}).WithContext(ctx).Where("id = ? AND deleteat IS NULL", idtodo).First(&data)
	if datatoUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

func Deletebyid(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	reqParamId := c.Param("id")
	idtodo, _ := strconv.Atoi(reqParamId)
	defer cancel()

	var datai models.Student
	datatoUpdate := DB.Model(models.Student{}).WithContext(ctx).Where("id = ? AND deleteat IS NULL", idtodo).First(&datai)
	if datatoUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not Found!!",
		})
		return
	}

	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	v := now.Format("02-01-2006 15:04:05")
	datai.Deleteat = &v
	result := DB.WithContext(ctx).Save(&datai)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to Delete the data",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": result.RowsAffected,
	})

}

func Pagination(c *gin.Context) {
	var pagi models.Paginationparam
	param := c.Request.URL.Query()
	pagi.Page, _ = strconv.Atoi(param.Get("page"))
	pagi.Offset, _ = strconv.Atoi(param.Get("offset"))

	if pagi.Page == 0 || pagi.Offset == 0 {
		pagi.Page = 5
		pagi.Offset = 10
	}
	var pagination models.Pagination

	data, err := helper.GetData(pagi.Page, pagi.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Not Found.",
		})
		return
	}
	pagination.Thispage = len(data)
	pagination.Page = pagi.Page
	pagination.Data = make([]models.Student, 0)
	if len(data) > 0 {
		pagination.Data = data
	}
	counter, err := helper.GetCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if int(counter)%pagi.Offset == 0 {
		pagination.Totalpages = int(counter) / pagi.Offset
	} else {
		pagination.Totalpages = (int(counter) / pagi.Offset) + 1
	}
	pagination.Totalcount = int(counter)
	c.JSON(http.StatusCreated, gin.H{
		"data": pagination,
	})
}

func SearchHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := c.Query("q")
	fields := c.Query("fields")

	if query == "" || fields == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid parameters",
		})
		return
	}

	var student []models.Student

	queryStr := fmt.Sprintf("deleteat IS NULL AND %s LIKE ?", fields)
	err := DB.Model(models.Student{}).WithContext(ctx).Where(queryStr, "%"+query+"%").Find(&student).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": student,
	})
}

func Signup(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed To Read Body",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed To Hash Password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	err = DB.Model(models.User{}).Create(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed To Create User",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed To Read Body",
		})
		return
	}

	var user models.User
	err := DB.Model(models.User{}).First(&user, "email = ?", body.Email).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Create Token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

// r.GET("/validate", middleware.RequireAuth, controller.Validate)
func Validate(c *gin.Context) {
	user, err := c.Get("user")
	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized User.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
