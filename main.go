package main

import (
	"fmt"
	"net/http"
	"os"
	"todo/jwt/database"

	//"todo/jwt/env"
	"todo/jwt/handler"

	"github.com/gorilla/mux"
)

func main() {
	//env.Setenv()
	database.Connectdatabase()
	database.Initdatabase()
	router := mux.NewRouter()
	router.HandleFunc("/signup", handler.Signup).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", handler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/logout", handler.Logout).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo", handler.CreateTodo).Methods("POST", "OPTIONS")
	router.HandleFunc("/todo/{todo_id}", handler.ToggleTodo).Methods("GET", "OPTIONS")
	router.HandleFunc("/todos", handler.GetTodo).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo/{todo_id}", handler.DeleteTodo).Methods("DELETE", "OPTIONS")

	/*router.HandleFunc("/user",)
	 */
	PORT := os.Getenv("PORT")
	fmt.Printf("On port: %s", PORT)
	http.ListenAndServe(":"+PORT, router)

}
