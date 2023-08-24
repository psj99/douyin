package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string
	UserID  uint
	VideoID uint
}