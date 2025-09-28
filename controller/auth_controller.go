package controller

import (
	"auth-session/dto"
	"auth-session/repository"
	"net/http"

	authUtils "auth-session/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authRepo *repository.AuthRepo
}

func NewAuthController(repo *repository.AuthRepo) *AuthController {
	return &AuthController{authRepo: repo}
}

func (c *AuthController) IsController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.authRepo.IsLoggedIn())
}

func (c *AuthController) SignUp(ctx *gin.Context) {
	var newUser dto.CreateUser
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := authUtils.ValidatePassword(newUser.Password, newUser.Username, newUser.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.authRepo.SignUp(&newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusAccepted, res)

}

func (c *AuthController) GetAllUsers(ctx *gin.Context) {
	res, err := c.authRepo.GetAllusers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})

	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var userPayload dto.Login
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.authRepo.Login(ctx, userPayload.Username, userPayload.Password, userPayload.RedirectUrl, userPayload.AppCode)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successfully"})
}
