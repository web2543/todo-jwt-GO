package model

import "gorm.io/gorm"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Password struct {
	UserID   uint
	User     Users
	Password string
}
type Users struct {
	gorm.Model
	User string  `gorm:"unique;->;<-:create"`
	Todo []Todos `gorm:"foreignKey:UserID"`
}

type Todos struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Done   bool   `json:"check" default:"false"`
	Todo   string `json:"text"`
}
