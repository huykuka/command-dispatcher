package db

import (
	"errors"
	"command-dispatcher/internal/core/services/hashing"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var hashingService = hashing.NewHashingService()

// Hash the password and handle the error
var usersSeed = []*User{
	{
		Email:    os.Getenv("SECRET_EMAIL"),
		Password: hashingService.HashPassword(os.Getenv("SECRET_PASSWORD")),
	},
}

func seedUser(db *gorm.DB) {
	log.Println("Seeding users...")

	for _, user := range usersSeed {
		var existingUser User
		err := db.Where("email = ?", user.Email).First(&existingUser).Error

		// Handle errors other than ErrRecordNotFound
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Failed to query user %s: %v", user.Email, err)
			continue
		}

		// If the user does not exist, create the new user
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", user.Email, err)
				continue
			}
			log.Printf("User %s created successfully", user.Email)
		} else {
			log.Printf("User %s already exists, skipping", user.Email)
		}
	}
}
