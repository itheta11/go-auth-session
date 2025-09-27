package router

import (
	"auth-session/controller"
	"auth-session/repository"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	repo := repository.NewAuthRepo(db)
	ctrl := controller.NewAuthController(repo)

	r.GET("/", func(ctx *gin.Context) { ctx.String(200, "API is running") })
	r.GET("/isLoggedin", ctrl.IsController)
	r.POST("/signup", ctrl.SignUp)
	r.GET("/users", ctrl.GetAllUsers)
	r.POST("/login", ctrl.Login)
	return r
}
