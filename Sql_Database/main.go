package main
 
import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)
 
type Post struct {
    Id   int
    Name string
    Text string
}
 
func main() {
    db, e := sql.Open("mysql", "rootuser:dhirajpatel@/posts")
    ErrorCheck(e)
 
    // close database after all work is done
    defer db.Close()
 
    PingDB(db)
 
    // INSERT INTO DB
    // prepare
    stmt, e := db.Prepare("insert into posts(id, Name, Text) values (?, ?, ?)")
    ErrorCheck(e)
 
    //execute
    res, e := stmt.Exec("5", "Post five", "Contents of post 5")
    ErrorCheck(e)
 
    id, e := res.LastInsertId()
    ErrorCheck(e)
 
    fmt.Println("Insert id", id)
 
    //Update db
    stmt, e := db.Prepare("update posts set Text=? where id=?")
    ErrorCheck(e)
 
    // execute
    res, e := stmt.Exec("This is post five", "5")
    ErrorCheck(e)
 
    a, e := res.RowsAffected()
    ErrorCheck(e)
 
    fmt.Println(a)
 
    // query all data
    rows, e := db.Query("select * from posts")
    ErrorCheck(e)
 
    var post = Post{}
 
    for rows.Next() {
        e = rows.Scan(&post.Id, &amp;post.Name, &post.Text)
        ErrorCheck(e)
        fmt.Println(post)
    }
 
    // delete data
    stmt, e := db.Prepare("delete from posts where id=?")
    ErrorCheck(e)
 
    // delete 5th post
    res, e := stmt.Exec("5")
    ErrorCheck(e)
 
    // affected rows
    a, e := res.RowsAffected()
    ErrorCheck(e)
 
    fmt.Println(a) // 1
}
 
func ErrorCheck(err error) {
    if err != nil {
        panic(err.Error())
    }
}
 
func PingDB(db *sql.DB) {
    err := db.Ping()
    ErrorCheck(err)
}