package repository

import (
	"auth-session/dto"
	"auth-session/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{DB: db}
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

func (repo *AuthRepo) Login(username string, password string, redirectUrl string, appCode string) error {
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

	return nil

}
