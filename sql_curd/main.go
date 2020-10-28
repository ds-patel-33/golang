package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Customers type
type Customers struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {

	// inserting a rows
	insert(Customers{101, "Dinesh"})
	insert(Customers{102, "Ramesh"})

	// updating the customers by id
	updateById(Customers{101, "Krishnan"})

	// select all customerss
	results := selectAll()

	// iterating a results
	fmt.Println("All Records :-")
	for results.Next() {
		var customers Customers
		results.Scan(&customers.Id, &customers.Name)
		fmt.Println(customers.Id, customers.Name)
	}

	// select customers by id
	result := selectById(101)
	fmt.Println("Selected Record:-")
	var customers Customers
	result.Scan(&customers.Id, &customers.Name)
	fmt.Println(customers.Id, customers.Name)

	// delete a customers by id
	delete(101)
}

// function to get a database connection
func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:dhirajpatel@tcp(127.0.0.1:3306)/customers")

	if err != nil {
		fmt.Println("Error! Getting connection...")
	}
	return db
}

// function to insert a row in customers table
func insert(customers Customers) {
	db := connect()
	insert, err := db.Query("INSERT INTO customers VALUES (?, ?)", customers.Id, customers.Name)
	if err != nil {
		fmt.Println("Error! Inserting records...")
	}
	defer insert.Close()
	fmt.Println("Inserted record...")
	defer db.Close()
}

// function to select all records from customers table
func selectAll() *sql.Rows {
	db := connect()
	results, err := db.Query("SELECT * FROM customers")
	if err != nil {
		fmt.Println("Error! Getting records...")
	}
	defer db.Close()
	return results
}

// function to select a customers record from table by customers id
func selectById(id int) *sql.Row {
	db := connect()
	result := db.QueryRow("SELECT * FROM customers WHERE id=?", id)
	defer db.Close()
	return result
}

// function to update a customers record by customers id
func updateById(customers Customers) {
	db := connect()
	db.QueryRow("UPDATE customers SET name=? WHERE id=?", customers.Name, customers.Id)
	fmt.Println("Record Updated....")
}

// function to delete a customers by customers id
func delete(id int) {
	db := connect()
	db.QueryRow("DELETE FROM customers WHERE id=?", id)
	fmt.Println("Record Deleted....")
}
