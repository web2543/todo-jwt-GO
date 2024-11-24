package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

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

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	Store Storage
}

type Storage interface {
	AddUser(credential Credentials) (Users, error)
	FindUser(username string) (Users, Password, error)
}

func NewAuth(s Storage) *Auth {
	return &Auth{
		Store: s,
	}
}

type Error struct {
	Massage string `json:"error"`
}
type UsersRequest struct {
	Token string `json:"jwt"`
	Users
}

func (a *Auth) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	var credential Credentials
	json.NewDecoder(r.Body).Decode(&credential)
	user, err := a.Store.AddUser(credential)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		err_massage := Error{
			Massage: "user exist",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&user)
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	var credential Credentials
	json.NewDecoder(r.Body).Decode(&credential)
	w.Header().Set("Content-Type", "application/json")
	user, password, err := a.Store.FindUser(credential.Username)
	if err != nil {
		err_massage := Error{
			Massage: "Not Found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	success := ComparePassword(credential.Password, password.Password)
	if !success {
		err_massage := Error{
			Massage: "Wrong password",
		}
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	token, exp_time, err := GenerateJWT(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err_massage := Error{
			Massage: "Internal Server Error",
		}
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "todo-token",
		Value:    token,
		Expires:  exp_time,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	var json_user UsersRequest
	json_user.Users = user
	json_user.Token = token
	json.NewEncoder(w).Encode(&json_user)

}
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	http.SetCookie(w, &http.Cookie{
		Name:     "todo-token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

}

func ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func GenerateJWT(users Users) (token string, expirationTime time.Time, err error) {
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

func Cors(w http.ResponseWriter, r *http.Request) {
	//url := strings.Split(r.Referer(), "/")
	Origin := os.Getenv("ORIGIN")
	w.Header().Set("Access-Control-Allow-Origin", Origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//fmt.Printf("%s//%s \n", url[0], url[2])
}
