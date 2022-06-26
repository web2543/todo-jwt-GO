package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"todo/jwt/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connectdatabase() {
	var err error
	dbpath := os.Getenv("DB_PATH")
	db, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		fmt.Println("Fail to connect Database")
	}
	fmt.Println("Connect Success")
}
func Initdatabase() {
	db.AutoMigrate(&model.Users{}, &model.Password{}, &model.Todos{})

}
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}

func AddUser(credential model.Credentials) (model.Users, error) {
	var user model.Users
	var password model.Password
	var err error
	user.User = credential.Username
	result := db.Create(&user)
	if result.Error != nil {
		return model.Users{}, errors.New("not unique")
	}
	result = db.Where("user = ?", credential.Username).First(&user)
	if result.Error != nil {
		return model.Users{}, result.Error
	}
	password.Password, err = hashPassword(credential.Password)
	if err != nil {
		return model.Users{}, err
	}
	password.UserID = user.ID
	db.Create(&password)
	return user, err
}
func FindUser(username string) (model.Users, model.Password, error) {
	var user model.Users
	var password model.Password
	result := db.Where("user = ?", username).First(&user)
	if result.Error != nil {
		return model.Users{}, model.Password{}, errors.New("not found")
	}
	result = db.Where("user_id = ?", user.ID).First(&password)
	if result.Error != nil {
		return model.Users{}, model.Password{}, errors.New("not found")
	}
	return user, password, nil
}
func AddTodo(todo model.Todos) error {
	result := db.Create(&todo)
	if result.Error != nil {
		return errors.New("fail to add new list")
	}
	return nil
}

func DeleteTodo(todo_id string, user_id uint) error {
	id, err := strconv.Atoi(todo_id)
	if err != nil {
		return errors.New("input number only")
	}
	var todo model.Todos
	result := db.Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return errors.New("fail to delete")
	}
	if todo.UserID != user_id {
		return errors.New("unauthorized")
	}
	result = db.Unscoped().Delete(&todo)
	if result.Error != nil {
		return errors.New("fail to delete")
	}
	return nil
}

func ToggleTodo(todo_id string, user_id uint) error {
	id, err := strconv.Atoi(todo_id)
	if err != nil {
		return errors.New("input number only")
	}
	var todo model.Todos
	result := db.Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return errors.New("fail to update")
	}
	fmt.Println(user_id, todo.UserID, id)
	if user_id != todo.UserID {
		return errors.New("unauthorized")
	}
	todo.Done = !(todo.Done)
	db.Save(&todo)
	return nil
}

func GetTodobyuser(user_id uint) ([]model.Todos, error) {
	/*id, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, errors.New("input number only")
	}*/
	var todos []model.Todos
	result := db.Where("user_id = ?", user_id).Find(&todos)
	if result.Error != nil {
		return nil, errors.New("not found")
	}
	return todos, nil
}

func GetTodobyid(todo_id string, user_id uint) (model.Todos, error) {
	id, err := strconv.Atoi(todo_id)
	if err != nil {
		return model.Todos{}, errors.New("input number only")
	}
	var todo model.Todos
	result := db.Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return model.Todos{}, errors.New("not found")
	}
	if user_id != todo.UserID {
		return model.Todos{}, errors.New("unauthorized")
	}
	return todo, nil
}
