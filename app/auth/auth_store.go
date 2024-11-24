package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthStore struct {
	db *gorm.DB
}

type Password struct {
	UserID   uint
	User     Users
	Password string
}
type Users struct {
	gorm.Model
	User string `gorm:"unique;->;<-:create"`
}

func NewAuthStore(db *gorm.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}

func (s *AuthStore) AddUser(credential Credentials) (Users, error) {
	var user Users
	var password Password
	var err error
	user.User = credential.Username
	result := s.db.Create(&user)
	if result.Error != nil {
		return Users{}, errors.New("not unique")
	}
	result = s.db.Where("user = ?", credential.Username).First(&user)
	if result.Error != nil {
		return Users{}, result.Error
	}
	password.Password, err = hashPassword(credential.Password)
	if err != nil {
		return Users{}, err
	}
	password.UserID = user.ID
	s.db.Create(&password)
	return user, err
}
func (s *AuthStore) FindUser(username string) (Users, Password, error) {
	var user Users
	var password Password
	result := s.db.Where("user = ?", username).First(&user)
	if result.Error != nil {
		return Users{}, Password{}, errors.New("not found")
	}
	result = s.db.Where("user_id = ?", user.ID).First(&password)
	if result.Error != nil {
		return Users{}, Password{}, errors.New("not found")
	}
	return user, password, nil
}
