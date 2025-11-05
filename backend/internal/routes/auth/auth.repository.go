package auth

import (
	"command-dispatcher/internal/config/db"
	"command-dispatcher/internal/core/services/hashing"
	"command-dispatcher/internal/models"
	"command-dispatcher/internal/routes/users"

	log "github.com/sirupsen/logrus"
)

type User db.User

type AuthRepository struct {
	hashingService hashing.HashingService
	userRepository users.UserRepository
}

func (r *AuthRepository) Login(userDto *models.LoginDTO) (*users.User, error) {
	user, err := r.userRepository.FindOne(userDto.Email)
	if err != nil {
		return nil, err
	}

	// Compare the hashed password with the provided password.
	err = r.hashingService.ComparePasswords(user.Password, userDto.Password) // Ensure you compare the hashed password
	if err != nil {
		return nil, err
	}

	// Log the successful login event
	log.Infof("User with email %s has logged in successfully.", user.Email)

	return user, nil
}
