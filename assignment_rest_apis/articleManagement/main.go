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

//Article Struct
type Article struct {
	Writer  string `json:"writer"`
	ID      int    `json:"id"`
	Subject string `json:"subject"`
	Type    string `json:"type"`
}

//Input struct
type Input struct {
	Writer  string `json:"writer"`
	Subject string `json:"subject"`
	Type    string `json:"type"`
}

const (
	//DBDriver name
	DBDriver = "mysql"

	//DBName Schema name
	DBName = "articlemanagement"

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

	router.HandleFunc("/article", addArticle).Methods("POST")

	router.HandleFunc("/article", allArticles).Methods("GET")

	router.HandleFunc("/article/{id}", udpateArticle).Methods("PUT")

	router.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")

	http.ListenAndServe(":5000", router)

}

func allArticles(w http.ResponseWriter, r *http.Request) {

	articles := make([]Article, 0)
	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Query(
		`SELECT id,
				 writer,
                  subject,
                  type
                FROM article;
               `)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		a := Article{}
		rows.Scan(&a.ID, &a.Writer, &a.Subject, &a.Type)
		articles = append(articles, a)
	}
	log.Println("Articles: ", articles)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func addArticle(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var inputArticle Input
	json.Unmarshal(bodyBytes, &inputArticle)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(
		`INSERT INTO article (writer,subject,type)
     VALUES (?, ?, ?)`,
		inputArticle.Writer, inputArticle.Subject, inputArticle.Type)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(inputArticle)

}

func udpateArticle(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var articleStruct Article

	json.Unmarshal(bodyBytes, &articleStruct)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(
		`UPDATE article SET  writer=? , subject=? , type=?  WHERE   id=?`,
		articleStruct.Writer, articleStruct.Subject, articleStruct.Type, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	ID, _ := strconv.Atoi(id)

	articleStruct.ID = ID
	json.NewEncoder(w).Encode(articleStruct)

}

func deleteArticle(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(`DELETE FROM article WHERE  id=?`, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte("Data Deleted !"))
	w.WriteHeader(http.StatusNoContent)

}
