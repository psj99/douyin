package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName        string `gorm:"unique"`
	Password        string
	Avatar          string
	BackgroundImage string
	Signature       string
	Videos          []Video   `gorm:"foreignKey:UserID"`
	Comments        []Comment `gorm:"foreignKey:UserID"`
	Likes           []*Video  `gorm:"many2many:like;"`
	Follows         []*User   `gorm:"many2many:follow;"`
}

const passWordCost = 12 //密码加密难度

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
