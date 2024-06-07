package main

import (
	"ginchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:mysql123$%^@(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.UserBasic{})
	// db.AutoMigrate(&models.Message{})
	// db.AutoMigrate(&models.GroupBasic{})
	// db.AutoMigrate(&models.Contact{})
	user := &models.UserBasic{}
	user.Name = "shen"
	db.Create(user)

	db.Model(user).Update("Password", 1234)
}
