package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

//Trainee Struct
type Trainee struct {
	Id     int    `json:id`
	Name   string `json:name`
	Batch  string `json:batch`
	Salary int    `json:salary`
	State  string `json:state`
}

const (

	//DBDriver name
	DBDriver = "mysql"

	//DBName Schema name
	DBName = "gin_curd"

	//DBUser for useranme of database
	DBUser = "root"

	//DBPassword for password
	DBPassword = "dhirajpatel"

	//DBURL for database connection url
	DBURL = DBUser + ":" + DBPassword + "@/" + DBName
)

func dbConn() (db *sql.DB) {

	db, err := sql.Open(DBDriver, DBURL)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(db)
	return db
}

type TraineeInputData struct {
	Name   string `json:"name"`
	Batch  string `json:"batch"`
	Salary int    `json:salary`
	State  string `json:state`
}

//GetTrainee function
func GetTrainee(c *gin.Context) {
	id := c.Param("id")
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM trainee WHERE id=?", id)
	exception := ""
	if err != nil {
		exception = err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	var salary int
	var name, state, batch string

	if selDB.Next() {
		err = selDB.Scan(&id, &name, &batch, &salary, &state)
		exception := ""
		if err != nil {
			exception = err.Error()
			c.JSON(500, gin.H{"exception": exception})
			return
		}
		fmt.Printf("name: %s; state: %s; batch: %s; salary: %d", name, state, batch, salary)

		c.JSON(200, gin.H{
			"id":     id,
			"name":   name,
			"batch":  batch,
			"state":  state,
			"salary": salary,
		})
		return
	}
	c.JSON(404, gin.H{"exception": "No data Found"})

}

//AddTrainee Function
func AddTrainee(c *gin.Context) {

	var traineeInputData TraineeInputData
	c.BindJSON(&traineeInputData)

	db := dbConn()
	_, err := db.Exec(`INSERT INTO trainee (name, batch, state, salary) VALUES(?,?,?,?)`, traineeInputData.Name, traineeInputData.Batch, traineeInputData.State, traineeInputData.Salary)
	exception := ""
	if err != nil {
		exception = err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return

	}

	fmt.Printf("Trainee Added:")
	fmt.Printf("name: %s; batch: %s; state: %s; salary: %d", traineeInputData.Name, traineeInputData.Batch, traineeInputData.State, traineeInputData.Salary)
	c.JSON(201, gin.H{"exception": exception, "data": traineeInputData})
}

//DeleteTrainee Function
func DeleteTrainee(c *gin.Context) {

	id := c.Param("id")
	db := dbConn()

	_, err := db.Exec(`DELETE FROM trainee WHERE id=?`, id)
	exception := ""
	if err != nil {
		c.JSON(500, gin.H{"exception": exception})
		return
	}
	defer db.Close()
	c.JSON(200, gin.H{"exception": exception})
}

//UpdateTrainee Function
func UpdateTrainee(c *gin.Context) {

	id := c.Param("id")
	var traineeInputData TraineeInputData
	c.BindJSON(&traineeInputData)

	db := dbConn()

	_, err := db.Exec(`UPDATE traine SET name=?, batch=?, salary=?, state=? WHERE id=?`, traineeInputData.Name, traineeInputData.Batch, traineeInputData.Salary, traineeInputData.State, id)
	exception := ""
	if err != nil {
		exception = err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return

	}

	fmt.Printf("name: %s; batch: %s; salary: %s; state: %s", traineeInputData.Name, traineeInputData.Batch, traineeInputData.State, traineeInputData.State)
	c.JSON(200, gin.H{"exception": exception, "data": traineeInputData})
}

func main() {
	router := gin.Default()
	api := router.Group("/trainee")

	{
		api.GET("/get/:id", GetTrainee)
		api.POST("/add", AddTrainee)
		api.PUT("/update/:id", UpdateTrainee)
		api.DELETE("/delete/:id", DeleteTrainee)
	}

	router.Run(":8080")
}
