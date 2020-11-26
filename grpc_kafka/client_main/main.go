package main

import (
	"GRPC/proto"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

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

func dbConn() (db *sql.DB) {
	//DBDriver name

	db, err := sql.Open(DBDriver, DBURL)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

	g := gin.Default()
	g.GET("/add/:username/:name", func(ctx *gin.Context) {
		a := ctx.Param("username")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Username"})
			return
		}

		b := ctx.Param("name")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
			return
		}

		req := &proto.Request{Username: string(a), Name: string(b)}
		if response, err := client.AddtoKafka(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(response.Status),
			})
			fmt.Println("receiving from Kafka")
			c, err := kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": "localhost:9092",
				"group.id":          "group-id-1",
				"auto.offset.reset": "earliest",
			})

			if err != nil {
				panic(err)
			}

			c.SubscribeTopics([]string{"jobs-topic1"}, nil)

			for {
				msg, err := c.ReadMessage(-1)

				if err == nil {
					fmt.Printf("Received from Kafka %s: %s\n", msg.TopicPartition, string(msg.Value))
					data := string(msg.Value)
					ctx.JSON(http.StatusOK, gin.H{
						"result1": fmt.Sprint(data),
					})
					addUserInMySQL(data)
				} else {
					fmt.Printf("Consumer error: %v (%v)\n", err, msg)
					break
				}
			}

			c.Close()

		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		}
	})

	if err := g.Run(":10000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

//addUserInMySql function
func addUserInMySQL(data string) {
	fmt.Println("Save to MySQL")
	db := dbConn()
	s := strings.Split(data, "&")
	in, err := db.Prepare("INSERT INTO user(username,name) VALUES(?,?)")
	if err != nil {
		panic(err)
	}
	in.Exec(s[0], s[1])
	fmt.Printf("data %s: %s\n", s[0], s[1])

	fmt.Printf("added to Mysql : %s", data)
}
