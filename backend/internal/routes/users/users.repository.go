package users

import (
	"command-dispatcher/internal/config/db"
	"strings"
)

type User db.User

type UserRepository struct{}

// func (u *UserRepository) Update(id *string, data *UpdateSettingDTO) (*Setting, error) {
// 	var setting Setting
// 	// Find the existing setting by ID
// 	if err := db.GetDB().Model(&Setting{}).Where("id = ?", id).First(&setting).Error; err != nil {
// 		log.Error(err.Error())
// 		return nil, err
// 	}

// 	// Update the setting with the provided data
// 	if err := db.GetDB().Model(&setting).Updates(data).Error; err != nil {
// 		log.Error(err.Error())
// 		return nil, err
// 	}

// 	// Return the updated setting
// 	return &setting, nil
// }

func (u *UserRepository) FindOne(email string) (*User, error) {
	var user User
	// Convert key to uppercase and find the setting by key
	if err := db.GetDB().Model(&User{}).Where("UPPER(email) = ?", strings.ToUpper(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
