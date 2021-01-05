package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

var err error

//Person structure
type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	City      string `json:"city"`
}

func main() {

	db, _ = gorm.Open("mysql", "root:dhirajpatel@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	db.AutoMigrate(&Person{})
	r := gin.Default()
	r.GET("/people/", GetPeople)
	r.GET("/people/:id", GetPerson)
	r.POST("/people", CreatePerson)
	r.PUT("/people/:id", UpdatePerson)
	r.DELETE("/people/:id", DeletePerson)
	r.Run(":8080")

}

//DeletePerson function
func DeletePerson(c *gin.Context) {

	id := c.Params.ByName("id")
	var person Person

	d := db.Where("id = ?", id).Delete(&person)

	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})

}

//UpdatePerson function
func UpdatePerson(c *gin.Context) {

	var person Person
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)

}

//CreatePerson function
func CreatePerson(c *gin.Context) {

	var person Person
	c.BindJSON(&person)
	db.Create(&person)
	c.JSON(200, person)

}

//GetPerson function
func GetPerson(c *gin.Context) {

	id := c.Params.ByName("id")
	var person Person

	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}

}

//GetPeople function
func GetPeople(c *gin.Context) {

	var people []Person

	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}

}
