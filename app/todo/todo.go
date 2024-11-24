package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"todo/jwt/app/auth"

	"github.com/gorilla/mux"
)

type Storage interface {
	AddTodo(todo TblTodo) error
	ToggleTodo(id string, userId uint) error
	DeleteTodo(id string, userId uint) error
	GetTodobyuser(userId uint) ([]TblTodo, error)
	GetTodobyid(id string, userId uint) (TblTodo, error)
}

type Todo struct {
	db Storage
}

type Error struct {
	Massage string `json:"error"`
}

func NewTodo(s Storage) *Todo {
	return &Todo{
		db: s,
	}
}

func (t *Todo) CreateTodo(w http.ResponseWriter, r *http.Request) {
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
	var todo TblTodo
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
	err = t.db.AddTodo(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err_massage := Error{
			Massage: "Internal Server Error",
		}
		json.NewEncoder(w).Encode(&err_massage)
		return
	}
	todos, _ := t.db.GetTodobyuser(user.UserID)
	json.NewEncoder(w).Encode(&todos)
}

func (t *Todo) ToggleTodo(w http.ResponseWriter, r *http.Request) {
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
	err = t.db.ToggleTodo(vars["todo_id"], user.UserID)
	if err != nil {
		unauthorized(w)
		return
	}
	user_todo, _ := t.db.GetTodobyid(vars["todo_id"], user.UserID)
	json.NewEncoder(w).Encode(&user_todo)

}

func (t *Todo) GetTodo(w http.ResponseWriter, r *http.Request) {
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
	user_todo, _ := t.db.GetTodobyuser(user.UserID)
	json.NewEncoder(w).Encode(&user_todo)

}

func (t *Todo) DeleteTodo(w http.ResponseWriter, r *http.Request) {
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
	err = t.db.DeleteTodo(vars["todo_id"], user.UserID)
	if err != nil {
		unauthorized(w)
		return
	}
	user_todo, _ := t.db.GetTodobyuser(user.UserID)
	json.NewEncoder(w).Encode(&user_todo)
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
func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	err_massage := Error{
		Massage: "Unauthorized",
	}
	json.NewEncoder(w).Encode(&err_massage)

}
