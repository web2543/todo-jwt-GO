package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Option struct {
	Path string
}

func New(opt Option) *gorm.DB {
	dbpath := os.Getenv("DB_PATH")
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		fmt.Println("Fail to connect Database")
	}
	return db

}

func Connectdatabase() {
	var err error
	dbpath := os.Getenv("DB_PATH")
	db, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		fmt.Println("Fail to connect Database")
	}
	fmt.Println("Connect Success")
}

// func Initdatabase() {
// 	db.AutoMigrate(&model.Users{}, &model.Password{}, &model.Todos{})

// }
