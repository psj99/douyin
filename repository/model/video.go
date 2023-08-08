package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	PlayUrl    string
	CoverUrl   string
	Title      string
	UserID     uint
	Comments   []Comment `gorm:"foreignKey:VideoID"`
	LikedUsers []*User   `gorm:"many2many:like;"`
}
