package config

import (
	"log"

	models "auth-session/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}
	db.Logger.LogMode(logger.Info)

	///auto migrate
	err = migrate(db)
	if err != nil {
		log.Fatal("Migration failed", err)
	}
	log.Println("✅ Connected to SQLite database")
	return db
}

func migrate(db *gorm.DB) error {
	// Create User table
	if !db.Migrator().HasTable(&models.User{}) {
		if err := db.Migrator().CreateTable(&models.User{}); err != nil {
			return err
		}
	}

	// Create Application table
	if !db.Migrator().HasTable(&models.Application{}) {
		if err := db.Migrator().CreateTable(&models.Application{}); err != nil {
			return err
		}
	}

	// Create UserAppSession table (depends on User and Application)
	if !db.Migrator().HasTable(&models.UserAppSession{}) {
		if err := db.Migrator().CreateTable(&models.UserAppSession{}); err != nil {
			return err
		}
	}

	return nil
}
