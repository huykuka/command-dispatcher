package db

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var settingSeed = []*Setting{
	{
		Key:   "DB_HOST",
		Value: "postgres",
	},
	{
		Key:   "DB_PASSWORD",
		Value: "postgres",
	},
	{
		Key:   "DB_PORT",
		Value: "5432",
	},
	{
		Key:   "DB_USER",
		Value: "postgres",
	},
	{
		Key:   "DB_NAME",
		Value: "postgres-go",
	},
	{
		Key:   "DB_BACKUP_DIR",
		Value: "/home/thaiqu/backup",
	},
	{
		Key:   "GENERAL_DB_TYPE",
		Value: "postgresql",
	},
	{
		Key:   "GENERAL_BACKUP_INTERVAL",
		Value: "midnight",
		//"midnight": Backup at midnight.
		//"daily": Backup once every day.
		//"weekly": Backup once every week.
		//"hourly": Backup every hour.
	},
}

func seedSetting(handler *gorm.DB) {
	log.Println("Seeding settings ......")
	for _, setting := range settingSeed {
		var existingSetting Setting
		if err := handler.Where("key = ?", setting.Key).First(&existingSetting).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				handler.Create(setting)
			}
		}
	}
}
