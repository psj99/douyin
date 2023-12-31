package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username        string `gorm:"uniqueIndex"`
	Password        string
	Avatar          string
	BackgroundImage string
	Signature       string
	Works           []Video   `gorm:"foreignKey:UserID"`
	Favorites       []*Video  `gorm:"many2many:favorite"`
	Comments        []Comment `gorm:"foreignKey:UserID"`
	Follows         []*User   `gorm:"many2many:follow;joinForeignKey:user_id;JoinReferences:follow_id"`
	Followers       []*User   `gorm:"many2many:follow;joinForeignKey:follow_id;JoinReferences:user_id"`
	Messages        []Message `gorm:"foreignKey:FromUserID"`
}

const passwordCost = 12 //密码加密难度

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
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
