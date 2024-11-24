package todo

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type TodoStore struct {
	db *gorm.DB
}

func NewTodoStore(s *gorm.DB) *TodoStore {
	return &TodoStore{
		db: s,
	}
}

type TblTodo struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Done   bool   `json:"check" default:"false"`
	Todo   string `json:"text"`
}

func (s *TodoStore) AddTodo(todo TblTodo) error {
	result := s.db.Create(&todo)
	if result.Error != nil {
		return errors.New("fail to add new list")
	}
	return nil
}

func (s *TodoStore) DeleteTodo(todo_id string, user_id uint) error {
	id, err := strconv.Atoi(todo_id)
	if err != nil {
		return errors.New("input number only")
	}
	var todo TblTodo
	result := s.db.Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return errors.New("fail to delete")
	}
	if todo.UserID != user_id {
		return errors.New("unauthorized")
	}
	result = s.db.Unscoped().Delete(&todo)
	if result.Error != nil {
		return errors.New("fail to delete")
	}
	return nil
}

func (s *TodoStore) ToggleTodo(todo_id string, user_id uint) error {
	id, err := strconv.Atoi(todo_id)
	if err != nil {
		return errors.New("input number only")
	}
	var todo TblTodo
	result := s.db.Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return errors.New("fail to update")
	}
	fmt.Println(user_id, todo.UserID, id)
	if user_id != todo.UserID {
		return errors.New("unauthorized")
	}
	todo.Done = !(todo.Done)
	s.db.Save(&todo)
	return nil
}

func (s *TodoStore) GetTodobyuser(user_id uint) ([]TblTodo, error) {
	/*id, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, errors.New("input number only")
	}*/
	var todos []TblTodo
	result := s.db.Where("user_id = ?", user_id).Find(&todos)
	if result.Error != nil {
		return nil, errors.New("not found")
	}
	return todos, nil
}

func (s *TodoStore) GetTodobyid(todo_id string, user_id uint) (TblTodo, error) {
	id, err := strconv.Atoi(todo_id)
	if err != nil {
		return TblTodo{}, errors.New("input number only")
	}
	var todo TblTodo
	result := s.db.Where("id = ?", id).First(&todo)
	if result.Error != nil {
		return TblTodo{}, errors.New("not found")
	}
	if user_id != todo.UserID {
		return TblTodo{}, errors.New("unauthorized")
	}
	return todo, nil
}
