package jwttoken

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt"
)

type JWTService struct {
	JWTKey []byte
}

var (
	instance *JWTService
	once     sync.Once
)

func NewJWTService() *JWTService {
	//singleton service
	once.Do(func() {
		instance = &JWTService{
			JWTKey: []byte(os.Getenv("HASH_JWT_KEY")),
		}
	})
	return instance
}

func (j *JWTService) GenerateAccessToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": email,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(instance.JWTKey)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return tokenString, nil
}

func (j *JWTService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.JWTKey), nil
	})
}
