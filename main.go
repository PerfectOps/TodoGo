package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Todo struct model Todo of db
type Todo struct {
	Task           string `json:"task"`
	Datetime_start string `json:"datetime_start"`
	Datetime_end   string `json:"datetime_end"`
}

// DataTask Data that will be added to the table
var DataTasks = Todo{Task: "five task", Datetime_start: "2021-03-11", Datetime_end: "2021-03-12"}

// ShowTask function show all tasks
func ShowTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo []Todo

	result, err := db.Query("SELECT * FROM todo")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var tasks Todo

		err := result.Scan(&tasks.Task, &tasks.Datetime_start, &tasks.Datetime_end)
		if err != nil {
			panic(err.Error())
		}
		todo = append(todo, tasks)
	}
	json.NewEncoder(w).Encode(todo)
}

// OneTask show one task on a given time range in URL
func OneTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	dateOne := r.URL.Query().Get("dateOne")
	dateTwo := r.URL.Query().Get("dateTwo")
	// read value from URL. Example: localhost:8000/todo/one?dateOne='2021-03-02'&dateTwo='2021-03-03'
	result := db.QueryRow(`SELECT task, datetime_start, datetime_end FROM todo 
						WHERE datetime_start BETWEEN $1 and $2 OR datetime_end BETWEEN $1 and $2`, dateOne, dateTwo)
	err := result.Scan(&todo.Task, &todo.Datetime_start, &todo.Datetime_end)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(todo)
}

// AddTask added data in table
func AddTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// DataTask := Todo{}

	// DataTask.Task = r.FormValue(DataTasks.Task)
	// DataTask.Datetime_start = r.FormValue(DataTasks.Datetime_start)
	// DataTask.Datetime_end = r.FormValue(DataTasks.Datetime_end)
	// output, err := json.Marshal(DataTask)
	// fmt.Println(string(output))

	result, err := db.Exec(`INSERT INTO todo (task, datetime_start, datetime_end) VALUES ($1, $2, $3)`, &DataTasks.Task, &DataTasks.Datetime_start, &DataTasks.Datetime_end)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(result)
}

func main() {
	var err error

	db, err = sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=postgres password=poMSta dbname=task sslmode=disable"))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/todo", ShowTask).Methods("GET")
	router.HandleFunc("/todo/one", OneTask).Methods("GET")
	router.HandleFunc("/todo/add", AddTask).Methods("POST")
	http.ListenAndServe(":8000", router)
}
