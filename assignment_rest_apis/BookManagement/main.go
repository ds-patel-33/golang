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

//Book Struct
type Book struct {
	Author string `json:"author"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
}

type Input struct {
	Author string `json:"author"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
}

const (
	//DBDriver name
	DBDriver = "mysql"

	//DBName Schema name
	DBName = "bookmanagement"

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

	router.HandleFunc("/book", addBook).Methods("POST")

	router.HandleFunc("/book", allBooks).Methods("GET")

	router.HandleFunc("/book/{id}", udpateBook).Methods("PUT")

	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")

	http.ListenAndServe(":5000", router)

}

func allBooks(w http.ResponseWriter, r *http.Request) {

	books := make([]Book, 0)
	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Query(
		`SELECT id,
				 author,
                  name,
                  price
                FROM book;
               `)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		b := Book{}
		rows.Scan(&b.ID, &b.Author, &b.Name, &b.Price)
		books = append(books, b)
	}
	log.Println("Books: ", books)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func addBook(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var inputBook Input
	json.Unmarshal(bodyBytes, &inputBook)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(
		`INSERT INTO book (author,name, price)
     VALUES (?, ?, ?)`,
		inputBook.Author, inputBook.Name, inputBook.Price)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(inputBook)

}

func udpateBook(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var bookStruct Book

	json.Unmarshal(bodyBytes, &bookStruct)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(
		`UPDATE book SET  author=? , name=? , price=?  WHERE   id=?`,
		bookStruct.Author, bookStruct.Name, bookStruct.Price, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	ID, _ := strconv.Atoi(id)

	bookStruct.ID = ID
	json.NewEncoder(w).Encode(bookStruct)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(`DELETE FROM book WHERE  id=?`, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte("Data Deleted !"))
	w.WriteHeader(http.StatusNoContent)

}
