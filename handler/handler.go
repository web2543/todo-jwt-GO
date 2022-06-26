package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo/jwt/auth"
	"todo/jwt/database"
	"todo/jwt/model"

	"github.com/gorilla/mux"
)

type Error struct {
	Massage string `json:"error"`
}
type Json_users struct {
	Token string `json:"jwt"`
	model.Users
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	var credential model.Credentials
	json.NewDecoder(r.Body).Decode(&credential)
	user, err := database.AddUser(credential)
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

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	var credential model.Credentials
	json.NewDecoder(r.Body).Decode(&credential)
	w.Header().Set("Content-Type", "application/json")
	user, password, err := database.FindUser(credential.Username)
	if err != nil {
		err_massage := Error{
			Massage: "Not Found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	success := auth.ComparePassword(credential.Password, password.Password)
	if !success {
		err_massage := Error{
			Massage: "Wrong password",
		}
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	token, exp_time, err := auth.GenerateJWT(user)
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
	var json_user Json_users
	json_user.Users = user
	json_user.Token = token
	json.NewEncoder(w).Encode(&json_user)

}
func Logout(w http.ResponseWriter, r *http.Request) {
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

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	cookie, err := r.Cookie("todo-token")
	if err != nil {
		unauthorized(w)
		//panic(err)
		return
	}
	user, err := auth.GetdataFromJWT(cookie.Value)
	if err != nil {
		unauthorized(w)
		//panic(err)
		return
	}
	var todo model.Todos
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		err_massage := Error{
			Massage: "Bad Request",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	todo.UserID = user.UserID
	err = database.AddTodo(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err_massage := Error{
			Massage: "Internal Server Error",
		}
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	todos, _ := database.GetTodobyuser(user.UserID)
	json.NewEncoder(w).Encode(&todos)
}

func ToggleTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	cookie, err := r.Cookie("todo-token")
	if err != nil {
		unauthorized(w)
		return
	}
	user, err := auth.GetdataFromJWT(cookie.Value)
	if err != nil {
		unauthorized(w)
		return
	}
	vars := mux.Vars(r)
	fmt.Println(user.UserID, vars["todo_id"])
	err = database.ToggleTodo(vars["todo_id"], user.UserID)
	if err != nil {
		unauthorized(w)
		return
	}
	user_todo, _ := database.GetTodobyid(vars["todo_id"], user.UserID)
	json.NewEncoder(w).Encode(&user_todo)

}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie("todo-token")
	if err != nil {
		unauthorized(w)
		return
	}
	user, err := auth.GetdataFromJWT(cookie.Value)
	if err != nil {
		unauthorized(w)
		return
	}
	user_todo, _ := database.GetTodobyuser(user.UserID)
	json.NewEncoder(w).Encode(&user_todo)

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		Cors(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	Cors(w, r)
	cookie, err := r.Cookie("todo-token")
	if err != nil {
		unauthorized(w)
		return
	}
	user, err := auth.GetdataFromJWT(cookie.Value)
	if err != nil {
		unauthorized(w)
		return
	}
	vars := mux.Vars(r)
	err = database.DeleteTodo(vars["todo_id"], user.UserID)
	if err != nil {
		unauthorized(w)
		return
	}
	user_todo, _ := database.GetTodobyuser(user.UserID)
	json.NewEncoder(w).Encode(&user_todo)
}
func Cors(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.Referer(), "/")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	fmt.Printf("%s//%s \n", url[0], url[2])
}
func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	err_massage := Error{
		Massage: "Unauthorized",
	}
	json.NewEncoder(w).Encode(&err_massage)

}
