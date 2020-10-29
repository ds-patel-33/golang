package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

// Employee struct
type Employee struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  []byte `json:"-"`
}

func main() {
	log.Println("Server is up on 8080 port")
	router := NewRouter()
	log.Fatalln(http.ListenAndServe(":8080", router))
}

const (
	//DBDriver name
	DBDriver = "mysql"

	//DBName Schema name
	DBName = "employees"

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

// NewRouter return all router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/employees", FindAllEmployees)
	router.POST("/employees", SaveEmployees)
	router.GET("/employees/:id", FindUserByID)
	router.PUT("/employees/:id", UpdateEmployee)
	router.DELETE("/employees/:id", DeleteEmployee)
	return router
}

//NewEmployee to cretae new employee
func NewEmployee() *Employee {
	return new(Employee)
}

// GetAllEmployees return slice of Employee that contains all employees in database
func (e *Employee) GetAllEmployees() ([]Employee, error) {
	employees := make([]Employee, 0)
	db, err := GetDB()
	if err != nil {
		return employees, err
	}
	defer db.Close()
	rows, err := db.Query(
		`SELECT id,
                  username,
                  first_name,
                  last_name
                FROM employees;
               `)
	if err != nil {
		return employees, err
	}
	for rows.Next() {
		emp := Employee{}
		rows.Scan(&emp.ID, &emp.Username, &emp.FirstName, &emp.LastName)
		employees = append(employees, emp)
	}
	log.Println("employees: ", employees)
	return employees, nil
}

// Save Employee
func (e *Employee) Save() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()
	bPass, err := getDefaultPassword()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`INSERT INTO employees (id , username, first_name, last_name, password)
     VALUES (?,?, ?, ?, ?)`,
		e.ID, e.Username, e.FirstName, e.LastName, bPass)
	if err != nil {
		return err
	}
	return nil
}

// FindByID return Employee by ID
func (e *Employee) FindByID(ID string) (Employee, error) {
	if ID == "" {
		return Employee{}, errors.New("ID can not be empty")
	}
	_, err := strconv.Atoi(ID)
	if err != nil {
		return Employee{}, errors.New("ID must be a number")
	}
	db, err := GetDB()
	emp := Employee{}
	if err != nil {
		return emp, err
	}
	defer db.Close()
	rows, err := db.Query(
		`SELECT id,
                     username,
                     first_name,
                     last_name
                  FROM employees
                 WHERE id = ` + ID)
	if err != nil {
		return emp, err
	}
	if rows.Next() {
		rows.Scan(&emp.ID, &emp.Username, &emp.FirstName, &emp.LastName)
		return emp, nil
	}
	return Employee{}, errors.New("No data found")
}

// Update existing Employee
func (e *Employee) Update(newEmployee Employee, id string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	strconv.Atoi(id)

	_, err = db.Exec(
		`UPDATE employees SET username = ?, first_name = ?, last_name=? , id = ?  WHERE id = ?`,
		newEmployee.Username, newEmployee.FirstName, newEmployee.LastName, newEmployee.ID, id)
	if err != nil {
		return err
	}
	return nil
}

// Delete user from database
func (e *Employee) Delete(ID string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM employees WHERE id = ?`, ID)
	return err
}

func getDefaultPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte("123"), bcrypt.MinCost)
}

// Index Handler
func Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	_, err := w.Write([]byte("Welcome!"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// FindAllEmployees return all Employee in database
func FindAllEmployees(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	employees, err := NewEmployee().GetAllEmployees()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(employees)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//SaveEmployees create a Employee
func SaveEmployees(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body := req.Body
	emp := NewEmployee()
	err := json.NewDecoder(body).Decode(emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer body.Close()
	err = emp.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User Successfully added!"))
	w.WriteHeader(http.StatusOK)
}

// FindUserByID return a Employee
func FindUserByID(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	emp, err := NewEmployee().FindByID(id)
	log.Println(emp.ID)
	if emp.ID == 0 {
		w.Write([]byte("No Employee Found!"))
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateEmployee update an existing Employee
func UpdateEmployee(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ID := params.ByName("id")
	if ID == "" {
		http.Error(w, "ID can not be empty", http.StatusBadRequest)
		return
	}
	var e Employee
	err := json.NewDecoder(req.Body).Decode(&e)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if e.ID == 0 {
		http.Error(w, "Please set ID in User information", http.StatusBadRequest)
		return
	}
	_, err = NewEmployee().FindByID(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = NewEmployee().Update(e, ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Employee Successfully Updated!"))
	w.WriteHeader(http.StatusOK)
}

// DeleteEmployee delete a Employee from database
func DeleteEmployee(w http.ResponseWriter, req *http.Request, prm httprouter.Params) {
	ID := prm.ByName("id")
	if ID == "" {
		http.Error(w, "Please send ID of Employee", http.StatusBadRequest)
		return
	}
	emp, err := NewEmployee().FindByID(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if emp.ID == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = emp.Delete(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Employee Successfully deleted!"))
	w.WriteHeader(http.StatusOK)
}
