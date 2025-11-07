package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Handler *gorm.DB

func Init() {
	dsn := "host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	var err error
	Handler, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = Handler.AutoMigrate(&CommandConfig{}, &CommandExecution{})

	if err != nil {
		return
	}

	seedDB(Handler)
}

func GetDB() *gorm.DB {
	return Handler
}

func seedDB(handler *gorm.DB) {
	// seedSetting(handler)
	// seedUser(handler)
}
