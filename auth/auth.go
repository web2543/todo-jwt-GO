package auth

import (
	"errors"
	"fmt"
	"os"
	"time"
	"todo/jwt/model"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   uint
	Username string
}
type UserData struct {
	UserID   uint
	Username string
}

func ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func GenerateJWT(users model.Users) (token string, expirationTime time.Time, err error) {
	expirationTime = time.Now().Add(time.Hour)
	SECRET_KEY := []byte(os.Getenv("SECRET_KEY"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: users.User,
		UserID:   users.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	})
	/*t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    users.User,
		//ID:        users.ID,
	},
	)*/
	token, err = t.SignedString(SECRET_KEY)

	//panic(err)
	return
}

func GetdataFromJWT(string_token string) (UserData, error) {
	var claim Claims
	SECRET_KEY := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.ParseWithClaims(string_token, &claim,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return UserData{}, errors.New("error signature")
		}
		//fmt.Println(err)
	}
	//fmt.Println(token)
	if token.Valid {
		var user UserData
		user.UserID = claim.UserID
		user.Username = claim.Username
		fmt.Println(claim.UserID, claim.Username)
		return user, nil
	}
	return UserData{}, errors.New("token invalid")

}
