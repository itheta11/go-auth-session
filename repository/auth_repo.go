package repository

import (
	"auth-session/dto"
	"auth-session/models"
	"auth-session/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func (repo *AuthRepo) IsLoggedIn(accessToken string, refreshToken string, loggedIn *dto.LoggedInDto, appCode string) (*dto.LoggedInDto, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return repo.JwtManager.RsaPublicKey, nil
	})

	if err != nil {
		return loggedIn, err
	}

	claims, _ := token.Claims.(jwt.Claims)
	expTime, _ := claims.GetExpirationTime()
	timeLeft := expTime.Sub(time.Now())
	if timeLeft < 0 {
		return loggedIn, err
	}

	username, err := claims.GetSubject()
	if username == "" || err != nil || appCode == "" {
		return loggedIn, err
	}

	sessionId, err := repo.GetSessionId(username, appCode, int(timeLeft), refreshToken)
	if err != nil {
		return loggedIn, err
	}

	loggedIn.IsLoggedIn = true
	loggedIn.SessionId = sessionId

	return loggedIn, nil
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

func (repo *AuthRepo) Login(ctx *gin.Context, username string, password string, redirectUrl string, appCode string) (dto.LoginTokenResponse, error) {
	var response dto.LoginTokenResponse
	var user models.User
	var app models.Application

	res := repo.DB.Find(&user, "username = ?", username)

	if res.Error != nil {
		return response, res.Error
	}

	check := CheckPasswordHash(user.Password, password)

	if !check {
		panic("wrong password")
	}

	resApp := repo.DB.Find(&app, "name = ?", appCode)

	if resApp.Error != nil {
		return response, res.Error
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
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	return response, nil
}

func (repo *AuthRepo) GetSessionId(username string, appCode string, expSeconds int, tokenHash string) (string, error) {
	var user models.User
	var app models.Application

	res := repo.DB.Find(&user, "username = ?", username)

	if res.Error != nil {
		return "", res.Error
	}

	resApp := repo.DB.Find(&app, "name = ?", appCode)
	if resApp.Error != nil {
		return "", resApp.Error
	}

	var session models.UserAppSession
	resSession := repo.DB.Model(&models.UserAppSession{}).
		Where("user_id = ?", user.ID).
		Where("app_id = ?", app.ID).
		Where("is_Active = ?", true).
		Where("strftime('%s','now') - strftime('%s', last_accessed_time) <= ?", expSeconds+5).
		Order("last_accessed_time desc").
		First(&session)

	if resSession.Error != nil {
		newSession := models.UserAppSession{
			ID:               uuid.New(),
			UserID:           user.ID,
			AppID:            app.ID,
			StartTime:        time.Now(),
			LastAccessedTime: time.Now(),
			IsActive:         true,
			TokenHash:        tokenHash,
		}

		result := repo.DB.Create(&newSession)

		if result.Error != nil {
			return "", result.Error
		}
		return newSession.ID.String(), nil
	}

	return session.ID.String(), nil

}
