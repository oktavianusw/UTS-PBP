package main

import (
	"fmt"
	"log"
	"net/http"
	"week2/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/rooms", controllers.GetAllRooms).Methods("GET")
	router.HandleFunc("/rooms/{id}", controllers.GetDetailRoom).Methods("GET")
	router.HandleFunc("/rooms/{id}/join/{account_id}", controllers.InsertRoom).Methods("POST")
	router.HandleFunc("/rooms/{id}/leave/{account_id}", controllers.LeaveRoom).Methods("POST")

	http.Handle("/", router)
	fmt.Println("Connected to port 1234")
	log.Println("Connected to port 1234")
	log.Fatal(http.ListenAndServe(":1234", router))
}
