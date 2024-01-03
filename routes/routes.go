package routes

import (
	"github.com/Sahil-4555/go-crud-api/controller"
	"github.com/Sahil-4555/go-crud-api/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	// Public Routes
	public := r.Group("/api")
	public.POST("/login", controller.Login)
	public.POST("/register", controller.Signup)

	// Protected Routes
	protected := r.Group("/api/admin")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.POST("/create", controller.Create)
	protected.PUT("/update/:id", controller.Update)
	protected.GET("/getall", controller.Getall)
	protected.GET("/getbyid/:id", controller.Getbyid)
	protected.DELETE("/deletebyid/:id", controller.Deletebyid)
	protected.POST("/uploadimage/:id", controller.UploadImage)
	protected.PUT("/updateimage/:id", controller.UpdateImage)
	protected.GET("/getimage/:id", controller.GetImage)
	protected.GET("/pagination", controller.Pagination)
	protected.GET("/searchandler", controller.SearchHandler)
	protected.GET("/validator", controller.Validate)
}
