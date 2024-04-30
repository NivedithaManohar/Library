package routes

import (
	"github.com/Library/controllers"
	"github.com/Library/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	auth := middleware.AuthMiddleware()
	r.POST("/login", controllers.Login)
	r.GET("/home", auth, controllers.Book)
	r.POST("/addBook", auth, controllers.Book)
	r.POST("/deleteBook", auth, controllers.RemoveBook)
}
