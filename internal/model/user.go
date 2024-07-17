package model

import (
	"gorm.io/gorm"
	"string_backend_0001/internal/database"
	"string_backend_0001/pkg"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

func GetUser(u User) (user *User, err error) {
	db := database.GetDB()
	err = db.Model(&User{}).Where(u).First(&user).Error
	return
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: pkg.Sha256(username, password),
	}
}

func (u *User) Create() error {
	return database.GetDB().Create(u).Error
}
