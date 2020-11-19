package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type InputData struct {
	a int `json:"a"`
	b int `json:"b"`
}

type OutputData struct {
	ans float64 `json:"ans"`
}

//Add Function
func Add(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	var outputData OutputData
	json.NewDecoder(r.Body).Decode(&inputData)

	a := float64(inputData.a)
	b := float64(inputData.b)

	fmt.Println(inputData.a)
	fmt.Println(inputData.b)

	add := a + b
	fmt.Print(add)
	outputData.ans = adds

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

func Sub(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	var outputData OutputData
	json.NewDecoder(r.Body).Decode(&inputData)

	a := float64(inputData.a)
	b := float64(inputData.b)

	fmt.Println(inputData.a)
	fmt.Println(inputData.b)

	add := a - b
	fmt.Print(add)
	outputData.ans = adds

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

func Mul(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	var outputData OutputData
	json.NewDecoder(r.Body).Decode(&inputData)

	a := float64(inputData.a)
	b := float64(inputData.b)

	fmt.Println(inputData.a)
	fmt.Println(inputData.b)

	add := a * b
	fmt.Print(add)
	outputData.ans = adds

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

func Div(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	var outputData OutputData
	json.NewDecoder(r.Body).Decode(&inputData)

	a := float64(inputData.a)
	b := float64(inputData.b)

	fmt.Println(inputData.a)
	fmt.Println(inputData.b)

	add := a / b
	fmt.Print(add)
	outputData.ans = adds

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

func main() {
	router := mux.NewRouter()
	r.HandleFunc("/add", Add).Methods(http.MethodPost)
	r.HandleFunc("/sub", Sub).Methods(http.MethodPost)
	r.HandleFunc("/mul", Mul).Methods(http.MethodPost)
	r.HandleFunc("/div", div).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}
