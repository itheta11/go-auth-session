package controller

import (
	"auth-session/dto"
	"auth-session/repository"
	"net/http"

	authUtils "auth-session/utils"

	"github.com/gin-gonic/gin"
)

const ACCESS_TOKEN = "access_token"
const REFRESH_TOKEN = "refresh_token"

type AuthController struct {
	authRepo *repository.AuthRepo
}

func NewAuthController(repo *repository.AuthRepo) *AuthController {
	return &AuthController{authRepo: repo}
}

func (c *AuthController) IsLoggedIn(ctx *gin.Context) {
	/// check tokens
	var loggedInUser dto.LoggedInDto
	// accessToken, err := ctx.Cookie(ACCESS_TOKEN)
	// refreshToken, er := ctx.Cookie(REFRESH_TOKEN)
	// if err != nil || accessToken == "" || er != nil || refreshToken == "" {
	// 	loggedInUser.IsLoggedIn = false
	// 	ctx.JSON(401, loggedInUser)
	// }

	var tempLogedInResponse dto.LoginTokenResponse

	if err := ctx.ShouldBindJSON(&tempLogedInResponse); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appName := ctx.Query("appCode")

	res, err := c.authRepo.IsLoggedIn(tempLogedInResponse.AccessToken, tempLogedInResponse.RefreshToken, &loggedInUser, appName)
	if err != nil {
		ctx.JSON(401, loggedInUser)
	}

	ctx.JSON(http.StatusOK, res)

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
	var userPayload dto.LoginPayload
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.authRepo.Login(ctx,
		userPayload.Username,
		userPayload.Password,
		userPayload.RedirectUrl,
		userPayload.AppCode)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": response.AccessToken, "refresh_token": response.RefreshToken})
}
