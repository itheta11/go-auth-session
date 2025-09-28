package repository

import (
	"auth-session/dto"
	"auth-session/models"
	"auth-session/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB         *gorm.DB
	JwtManager *utils.JwtManager
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	jwtManager := utils.NewJwtManager()
	return &AuthRepo{DB: db, JwtManager: jwtManager}
}

func (repo *AuthRepo) IsLoggedIn() bool {
	return true
}

func (repo *AuthRepo) SignUp(user *dto.CreateUser) (dto.User, error) {
	password, err := HashPassword(user.Password)
	if err != nil {
		log.Fatal("Error in creating password")
		return dto.User{}, err
	}
	newUser := models.User{
		ID:         uuid.New(),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		Password:   password,
		Email:      user.Email,
		Created:    time.Now(),
		ModifiedAt: time.Now(),
	}

	result := repo.DB.Create(&newUser)

	res := dto.User{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Username:  newUser.Username,
		Email:     newUser.Email,
	}
	return res, result.Error
}

func (repo *AuthRepo) GetAllusers() ([]dto.User, error) {
	var allUsers []models.User
	var res []dto.User
	result := repo.DB.Find(&allUsers)

	if result.Error != nil {
		return res, result.Error
	}

	for _, u := range allUsers {
		res = append(res, dto.User{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Username:  u.Username,
			Email:     u.Email,
		})
	}
	return res, result.Error
}

func (repo *AuthRepo) Login(ctx *gin.Context, username string, password string, redirectUrl string, appCode string) error {
	var user models.User
	var app models.Application

	res := repo.DB.Find(&user, "username = ?", username)

	if res.Error != nil {
		return res.Error
	}

	check := CheckPasswordHash(user.Password, password)

	if !check {
		panic("wrong password")
	}

	resApp := repo.DB.Find(&app, "name = ?", appCode)

	if resApp.Error != nil {
		return res.Error
	}

	accessToken, err := repo.JwtManager.CreateAccessToken(user.Username, time.Now())
	if err != nil {
		panic(err)
	}

	refreshToken, err := repo.JwtManager.CreateRefreshToken(user.Username, time.Now())
	if err != nil {
		panic(err)
	}
	ctx.SetCookie("access_token",
		accessToken,
		int(repo.JwtManager.AccessTokenExpiry.Seconds()),
		"/",
		"",
		true,
		true)

	ctx.SetCookie("refresh_token",
		refreshToken,
		int(repo.JwtManager.RefreshTokenExpiry.Seconds()),
		"/",
		"",
		true,
		true)

	return nil
}
