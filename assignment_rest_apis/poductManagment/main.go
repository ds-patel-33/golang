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

//Product Struct
type Product struct {
	ID       int    `json:"id"`
	Quantity string `json:"author"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
}

//Input Struct
type Input struct {
	Quantity string `json:"author"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
}

const (
	DBDriver   = "mysql"
	DBName     = "productmanagement"
	DBUser     = "root"
	DBPassword = "dhirajpatel"
	DBURL      = DBUser + ":" + DBPassword + "@/" + DBName
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
	router.HandleFunc("/product", addProduct).Methods("POST")
	router.HandleFunc("/product", allProducts).Methods("GET")
	router.HandleFunc("/product/{id}", udpateProduct).Methods("PUT")
	router.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")
	http.ListenAndServe(":5000", router)

}

func allProducts(w http.ResponseWriter, r *http.Request) {

	products := make([]Product, 0)
	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Query(
		`SELECT id,
				 quantity,
                  name,
                  price
                FROM product;
               `)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		p := Product{}
		rows.Scan(&p.ID, &p.Quantity, &p.Name, &p.Price)
		products = append(products, p)
	}
	log.Println("products: ", products)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func addProduct(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var inputProduct Input
	json.Unmarshal(bodyBytes, &inputProduct)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(
		`INSERT INTO product (quantity,name, price)
     VALUES (?, ?, ?)`,
		inputProduct.Quantity, inputProduct.Name, inputProduct.Price)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(inputProduct)

}

func udpateProduct(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var productStruct Product

	json.Unmarshal(bodyBytes, &productStruct)

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(
		`UPDATE product SET  quantity=? , name=? , price=?  WHERE   id=?`,
		productStruct.Quantity, productStruct.Name, productStruct.Price, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	ID, _ := strconv.Atoi(id)

	productStruct.ID = ID
	json.NewEncoder(w).Encode(productStruct)

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	db, err := GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(`DELETE FROM product WHERE  id=?`, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte("Data Deleted !"))
	w.WriteHeader(http.StatusNoContent)

}
