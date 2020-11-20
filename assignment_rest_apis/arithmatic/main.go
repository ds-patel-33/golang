package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

//InputData Struct
type InputData struct {
	A int `json:"a"`
	B int `json:"b"`
}

//OutputData Struct
type OutputData struct {
	Ans float64 `json:"ans"`
}

//Add Function
func Add(w http.ResponseWriter, r *http.Request) {

	var outputData OutputData
	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	var inputData InputData
	json.Unmarshal(bodyBytes, &inputData)

	fmt.Println(inputData.A)
	fmt.Println(inputData.B)

	a := float64(inputData.A)
	b := float64(inputData.B)

	add := a + b
	fmt.Print(add)
	outputData.Ans = add

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

//Sub function
func Sub(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	var outputData OutputData
	json.NewDecoder(r.Body).Decode(&inputData)

	a := float64(inputData.A)
	b := float64(inputData.B)

	fmt.Println(inputData.A)
	fmt.Println(inputData.B)

	add := a - b
	fmt.Print(add)
	outputData.Ans = add

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

//Mul function
func Mul(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	var outputData OutputData
	json.NewDecoder(r.Body).Decode(&inputData)

	a := float64(inputData.A)
	b := float64(inputData.B)

	fmt.Println(inputData.A)
	fmt.Println(inputData.B)

	add := a * b
	fmt.Print(add)
	outputData.Ans = add

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

//Div function
func Div(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	var outputData OutputData
	json.NewDecoder(r.Body).Decode(&inputData)

	a := float64(inputData.A)
	b := float64(inputData.B)

	fmt.Println(inputData.A)
	fmt.Println(inputData.B)

	add := a / b
	fmt.Print(add)
	outputData.Ans = add

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(outputData)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/add", Add).Methods("POST")
	r.HandleFunc("/sub", Sub).Methods("POST")
	r.HandleFunc("/mul", Mul).Methods("POST")
	r.HandleFunc("/div", Div).Methods("POST")
	http.ListenAndServe(":5000", r)
}
