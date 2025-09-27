package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	FirstName  string    `gorm:"size:100;not null"`
	LastName   string    `gorm:"size:100"`
	Username   string    `gorm:"size:50;unique;not null"`
	Email      string    `gorm:"size:100;unique;not null"`
	Password   string    `gorm:"size:255;not null"`
	Created    time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}

type Application struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name string    `gorm:"size:100;not null;unique"`
}

// UserAppSession table
type UserAppSession struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID uuid.UUID
	AppID  uuid.UUID

	User        User        `gorm:"foreignKey:UserID;"`
	Application Application `gorm:"foreignKey:AppID;"`
}
