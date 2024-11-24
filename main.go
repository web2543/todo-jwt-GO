package main

import (
	"fmt"
	"net/http"
	"os"

	"todo/jwt/app/auth"
	"todo/jwt/app/todo"
	"todo/jwt/database"

	"github.com/gorilla/mux"
)

func main() {
	var (
		con = database.New(database.Option{
			Path: os.Getenv("DBPATH"),
		})
		authStore   = auth.NewAuthStore(con)
		authHandler = auth.NewAuth(authStore)
		todoStore   = todo.NewTodoStore(con)
		todoHandler = todo.NewTodo(todoStore)
	)

	router := mux.NewRouter()
	router.HandleFunc("/signup", authHandler.Signup).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/logout", authHandler.Logout).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo", todoHandler.CreateTodo).Methods("POST", "OPTIONS")
	router.HandleFunc("/todo/{todo_id}", todoHandler.ToggleTodo).Methods("GET", "OPTIONS")
	router.HandleFunc("/todos", todoHandler.GetTodo).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo/{todo_id}", todoHandler.DeleteTodo).Methods("DELETE", "OPTIONS")

	/*router.HandleFunc("/user",)
	 */
	PORT := os.Getenv("PORT")
	fmt.Printf("On port: %s", PORT)
	http.ListenAndServe(":"+PORT, router)

}
