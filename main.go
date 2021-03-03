package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Todo struct {
	task           string    `json:"task"`
	datetime_start time.Time `json:"datetime_start"`
	datetime_end   time.Time `json:"datetime_end"`
}

var database *sql.DB
var todo []Todo

func ShowTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := database.Query("SELECT task, datetime_start, datetime_end FROM todo")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	for result.Next() {
		var tasks Todo
		err := result.Scan(&tasks.task, &tasks.datetime_start, &tasks.datetime_end)
		if err != nil {
			panic(err.Error())
		}
		todo = append(todo, tasks)
	}
	json.NewEncoder(w).Encode(todo)
}

func OneTask(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT * FROM todo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}

func AddTask(w http.ResponseWriter, r *http.Request) {

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

}

func main() {
	connect := "user=postgres password=poMSta dbname=task sslmode=disable"
	db, err := sql.Open("postgres", connect) // Connect to DB
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/todo", ShowTask).Methods("GET")
	router.HandleFunc("/todo/{datatime}", OneTask).Methods("GET")
	router.HandleFunc("/todo/add", AddTask).Methods("POST")
	router.HandleFunc("/todo/delete/{}", DeleteTask).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}
