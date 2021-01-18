package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (

	//DBDriver name
	DBDriver = "mysql"

	//DBName Schema name
	DBName = "social_network"

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

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{

		//Authentication
		api.POST("/signup", UserSignup)
		api.POST("/login", UserLogin)

		//post
		api.POST("/create_new_post", CreateNewPost)
		api.DELETE("/delete_post", DeletePost)
		api.PUT("/update_post", UpdatePost)
		api.POST("/follow", Follow)
		api.DELETE("/unfollow", Unfollow)
		api.POST("/like", Like)
		api.POST("/unlike", Unlike)

		//get
		api.GET("/index/:id", Index)
		api.GET("/profile/:id", Profile)

	}

	router.Run(":8080")
}

//CreateNewPostInput struct
type CreateNewPostInput struct {
	UserID  int    `json:"userid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

//CreateNewPost function
func CreateNewPost(c *gin.Context) {

	var input CreateNewPostInput
	c.BindJSON(&input)
	title := strings.TrimSpace(input.Content)
	content := strings.TrimSpace(input.Title)
	id := input.UserID

	db := dbConn()

	stmt, _ := db.Prepare("INSERT INTO posts(title, content, createdBy, createdAt) VALUES (?, ?, ?, ?)")
	rs, err := stmt.Exec(title, content, id, time.Now())
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	insertID, _ := rs.LastInsertId()

	resp := map[string]interface{}{
		"postID": insertID,
		"mssg":   "Post Created!!",
	}
	c.JSON(200, gin.H{"data": resp})
}

//UpdatePostInput struct
type UpdatePostInput struct {
	PostID  int    `json:"postid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePost route
func UpdatePost(c *gin.Context) {

	var input UpdatePostInput
	c.BindJSON(&input)
	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)
	postID := input.PostID

	db := dbConn()

	_, err := db.Exec("UPDATE posts SET title=?, content=? WHERE postID=?", title, content, postID)

	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	c.JSON(200, gin.H{"msg": "Post Upodated !"})
}

//DeletePostInput struct
type DeletePostInput struct {
	PostID int `json:"postid"`
}

//DeletePost function
func DeletePost(c *gin.Context) {
	var input DeletePostInput
	c.BindJSON(&input)
	postID := input.PostID

	db := dbConn()

	_, err := db.Exec("DELETE FROM posts WHERE postID=?", postID)
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}
	c.JSON(200, gin.H{"msg": "Post Deleted !"})
}

//UserInput struct
type UserInput struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordAgain string `json:"password_again"`
}

func hash(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash
}

//UserSignup function
func UserSignup(c *gin.Context) {

	var input UserInput
	c.BindJSON(&input)
	username := strings.TrimSpace(input.Username)
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)
	passwordAgain := strings.TrimSpace(input.PasswordAgain)

	mailErr := checkmail.ValidateFormat(email)

	resp := make(map[string]interface{})

	db := dbConn()

	var (
		userCount  int
		emailCount int
	)

	db.QueryRow("SELECT COUNT(id) AS userCount FROM users WHERE username=?", username).Scan(&userCount)
	db.QueryRow("SELECT COUNT(id) AS emailCount FROM users WHERE email=?", email).Scan(&emailCount)

	if username == "" || email == "" || password == "" || passwordAgain == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if len(username) < 4 || len(username) > 32 {
		resp["mssg"] = "Username should be between 4 and 32"
	} else if mailErr != nil {
		resp["mssg"] = "Invalid email format!!"
	} else if password != passwordAgain {
		resp["mssg"] = "Passwords don't match"
	} else if userCount > 0 {
		resp["mssg"] = "Username already exists!!"
	} else if emailCount > 0 {
		resp["mssg"] = "Email already exists!!"
	} else {

		stmt, _ := db.Prepare("INSERT INTO users(username, email, password, joined) VALUES (?, ?, ?, ?)")
		_, err := stmt.Exec(username, email, password, time.Now())

		if err != nil {
			exception := err.Error()
			c.JSON(500, gin.H{"exception": exception})
			return
		}

		resp["success"] = true
		resp["msg"] = "Hello, " + username + "You have Signedup " + "!!"

	}
	c.JSON(200, gin.H{"msg": resp})
}

//LoginInput struct
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//UserLogin function
func UserLogin(c *gin.Context) {
	resp := make(map[string]interface{})

	var input LoginInput
	c.BindJSON(&input)
	rusername := strings.TrimSpace(input.Username)
	rpassword := strings.TrimSpace(input.Password)

	db := dbConn()

	var (
		userCount int
		id        int
		username  string
		password  string
	)

	db.QueryRow("SELECT COUNT(id) AS userCount, id, username, password FROM users WHERE username=?", rusername).Scan(&userCount, &id, &username, &password)

	//err := bcrypt.CompareHashAndPassword([]byte(password), []byte(rpassword))
	// if err != nil {
	// 	exception := err.Error()
	// 	c.JSON(500, gin.H{"exception": exception})
	// 	return
	// }

	if rusername == "" || rpassword == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if userCount == 0 {
		resp["mssg"] = "Invalid username!!"

	} else {

		resp["msg"] = "Hello, " + username + "!!"
		resp["success"] = true
	}
	c.JSON(200, gin.H{"msg": resp})
}

//FollowInput struct
type FollowInput struct {
	FollowBy int `json:"followBy"`
	FollowTo int `json:"followTo"`
}

// Follow route
func Follow(c *gin.Context) {
	var input FollowInput
	c.BindJSON(&input)

	followBy := input.FollowBy
	followTo := input.FollowTo

	db := dbConn()
	var Username string
	db.QueryRow("SELECT username AS Username FROM users WHERE id=?", followTo).Scan(&Username)
	username := Username

	stmt, _ := db.Prepare("INSERT INTO follow(followBy, followTo, followTime) VALUES(?, ?, ?)")
	_, err := stmt.Exec(followBy, followTo, time.Now())
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	c.JSON(200, gin.H{
		"mssg": "Followed " + username + "!!",
	})
}

// Unfollow route
func Unfollow(c *gin.Context) {
	var input FollowInput
	c.BindJSON(&input)

	followBy := input.FollowBy
	followTo := input.FollowTo

	db := dbConn()
	var Username string
	db.QueryRow("SELECT username AS Username FROM users WHERE id=?", followTo).Scan(&Username)
	username := Username

	stmt, _ := db.Prepare("DELETE FROM follow WHERE followBy=? AND followTo=?")
	_, err := stmt.Exec(followBy, followTo)
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	c.JSON(200, gin.H{
		"mssg": "Unfollowed " + username + "!!",
	})
}

//LikeInput struct
type LikeInput struct {
	PostID int `json:"postid"`
	UserID int `json:"userid"`
}

// Like post route
func Like(c *gin.Context) {

	var input LikeInput
	c.BindJSON(&input)
	post := input.PostID
	id := input.UserID

	db := dbConn()

	stmt, _ := db.Prepare("INSERT INTO likes(postID, likeBy, likeTime) VALUES (?, ?, ?)")
	_, err := stmt.Exec(post, id, time.Now())
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	c.JSON(200, gin.H{
		"mssg": "Post Liked!!",
	})
}

// Unlike post route
func Unlike(c *gin.Context) {

	var input LikeInput
	c.BindJSON(&input)
	post := input.PostID
	id := input.UserID

	db := dbConn()

	stmt, _ := db.Prepare("DELETE FROM likes WHERE postID=? AND likeBy=?")
	_, err := stmt.Exec(post, id)
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	c.JSON(200, gin.H{
		"mssg": "Post Unliked!!",
	})
}

//get requests

//Index Page
func Index(c *gin.Context) {

	id := c.Param("id")
	db := dbConn()
	var (
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)
	feeds := []interface{}{}

	stmt, _ := db.Prepare("SELECT posts.postID, posts.title, posts.content, posts.createdBy, posts.createdAt from posts, follow WHERE follow.followBy=? AND follow.followTo = posts.createdBy ORDER BY posts.postID DESC")
	rows, err := stmt.Query(id)
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	for rows.Next() {
		rows.Scan(&postID, &title, &content, &createdBy, &createdAt)

		//likedby
		var likeBy int
		LikedBy := []interface{}{}

		stmt, _ := db.Prepare("SELECT likeBy FROM likes WHERE postID=?")
		rows, err := stmt.Query(postID)
		if err != nil {
			exception := err.Error()
			c.JSON(500, gin.H{"exception": exception})
			return
		}

		for rows.Next() {
			rows.Scan(&likeBy)
			LikedBy = append(LikedBy, &likeBy)
		}

		feed := map[string]interface{}{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
			"likes":     len(LikedBy),
			"likedBy":   LikedBy,
		}
		feeds = append(feeds, feed)
	}

	c.JSON(200, gin.H{"posts": feeds})
}

// Profile Page
func Profile(c *gin.Context) {

	user := c.Param("id")
	db := dbConn()

	// VARS FOR USER DETAILS
	var (
		userCount int
		userID    int
		username  string
		email     string
		bio       string
	)

	// VARS FOR POSTS
	var (
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)

	posts := []interface{}{}

	var (
		followers  int
		followings int
	)

	// USER DETAILS
	db.QueryRow("SELECT COUNT(id) AS userCount, id AS userID, username, email, bio FROM users WHERE id=?", user).Scan(&userCount, &userID, &username, &email, &bio)

	fmt.Println(userCount)

	// POSTS
	stmt, err := db.Prepare("SELECT * FROM posts WHERE createdBy=? ORDER BY postID DESC")
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	for rows.Next() {
		rows.Scan(&postID, &title, &content, &createdBy, &createdAt)

		//likedby
		var likeBy int
		LikedBy := []interface{}{}

		stmt, _ := db.Prepare("SELECT likeBy FROM likes WHERE postID=?")
		rows, err := stmt.Query(postID)
		if err != nil {
			exception := err.Error()
			c.JSON(500, gin.H{"exception": exception})
			return
		}

		for rows.Next() {
			rows.Scan(&likeBy)
			LikedBy = append(LikedBy, &likeBy)
		}

		post := map[string]interface{}{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
			"likes":     len(LikedBy),
			"likedBy":   LikedBy,
		}
		posts = append(posts, post)
	}

	db.QueryRow("SELECT COUNT(followID) AS followers FROM follow WHERE followTo=?", user).Scan(&followers)  // FOLLOWERS
	db.QueryRow("SELECT COUNT(followID) AS followers FROM follow WHERE followBy=?", user).Scan(&followings) // FOLLOWINGS

	if err != nil {
		exception := err.Error()
		c.JSON(500, gin.H{"exception": exception})
		return
	}

	c.JSON(200, gin.H{
		"name": username,
		"user": gin.H{
			"id":       strconv.Itoa(userID),
			"username": username,
			"email":    email,
			"bio":      bio,
		},
		"posts":      posts,
		"followers":  followers,
		"followings": followings,
	})

}
