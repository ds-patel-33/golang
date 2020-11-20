package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Todo Struct
type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

const (
	//DBDriver name
	DBDriver = "mysql"

	//DBName Schema name
	DBName = "todo"

	//DBUser for useranme of database
	DBUser = "root"

	//DBPassword for password
	DBPassword = "dhirajpatel"

	//DBURL for database connection url
	DBURL = DBUser + ":" + DBPassword + "@/" + DBName
)

// GetDB return DB
func GetDB() (*sql.DB, error) {
	db, err := sql.Open(DBDriver, DBURL)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/tasks", post).Methods("POST")

	router.HandleFunc("/tasks", get).Methods("GET")

	router.HandleFunc("/tasks/{uid}/{id}", put).Methods("PUT")

	router.HandleFunc("/tasks/{uid}/{id}", delete).Methods("DELETE")

	http.ListenAndServe(":5000", router)

}

func get(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	fmt.Println(bodyBytes)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var todoStruct Todo
	json.Unmarshal(bodyBytes, &todoStruct)

	todos := make([]Todo, 0)
	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Query(
		`SELECT userId,
                  id,
                  title,
                  completed
                FROM todo;
               `)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		t := Todo{}
		rows.Scan(&t.UserID, &t.ID, &t.Title, &t.Completed)
		todos = append(todos, t)
	}
	log.Println("todos: ", todos)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func post(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var todoStruct Todo
	json.Unmarshal(bodyBytes, &todoStruct)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(
		`INSERT INTO todo (userId, id, title, completed)
     VALUES (?, ?, ?, ?)`,
		todoStruct.UserID, todoStruct.ID, todoStruct.Title, todoStruct.Completed)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(todoStruct)

}

func put(w http.ResponseWriter, r *http.Request) {
	userid := mux.Vars(r)["uid"]
	id := mux.Vars(r)["id"]

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var todoStruct Todo

	json.Unmarshal(bodyBytes, &todoStruct)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(
		`UPDATE todo SET  title=? , completed=?  WHERE (userId=? AND id=?)`,
		todoStruct.Title, todoStruct.Completed, userid, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	ID, _ := strconv.Atoi(id)
	UID, _ := strconv.Atoi(id)
	todoStruct.ID = ID
	todoStruct.UserID = UID
	json.NewEncoder(w).Encode(todoStruct)

}

func delete(w http.ResponseWriter, r *http.Request) {

	userid := mux.Vars(r)["uid"]
	id := mux.Vars(r)["id"]

	fmt.Println(userid)
	fmt.Println(id)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(`DELETE FROM todo WHERE (userId=? AND id=?)`, userid, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte("Data Deleted !"))
	w.WriteHeader(http.StatusNoContent)

}
