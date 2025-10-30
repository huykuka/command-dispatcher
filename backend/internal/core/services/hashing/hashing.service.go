package hashing

import (
	"errors"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type HashingService struct {
	passwordKey []byte
}

var (
	instance *HashingService
	once     sync.Once
)

// NewHashingService returns a singleton instance of HashingService.
func NewHashingService() *HashingService {
	//singleton service
	once.Do(func() {
		instance = &HashingService{
			passwordKey: []byte(os.Getenv("HASH_PASSWORD_KEY")),
		}
	})
	return instance
}

func (s *HashingService) HashPassword(password string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedBytes)

}

func (s *HashingService) ComparePasswords(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}
