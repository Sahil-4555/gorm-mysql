package routes

import (
	"github.com/Sahil-4555/go-crud-api/controller"
	"github.com/Sahil-4555/go-crud-api/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(incomingRoutes *gin.Engine) {
	// r.GET("/validate", middleware.RequireAuth, controller.Validate)
	incomingRoutes.POST("/create", middleware.RequireAuth, controller.Create)
	incomingRoutes.PUT("/update/:id", controller.Update)
	incomingRoutes.GET("/getall", controller.Getall)
	incomingRoutes.GET("/getbyid/:id", controller.Getbyid)
	incomingRoutes.DELETE("/deletebyid/:id", controller.Deletebyid)
	incomingRoutes.POST("/uploadimage/:id", controller.UploadImage)
	incomingRoutes.PUT("/updateimage/:id", controller.UpdateImage)
	incomingRoutes.GET("/getimage/:id", controller.GetImage)
	incomingRoutes.GET("/pagination", controller.Pagination)
	incomingRoutes.GET("/searchandler", controller.SearchHandler)
	incomingRoutes.POST("/signin", controller.Login)
	incomingRoutes.POST("/signup", controller.Signup)
	incomingRoutes.GET("/validator", middleware.RequireAuth, controller.Validate)
}
