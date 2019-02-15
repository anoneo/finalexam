package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	//db, err = sql.Open("postgres", url)
	url := os.Getenv("DATABASE_URL")
	fmt.Printf("urlllll: %q \n", url)
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("can't connect database : ", err)
	}
	return db
}

func InsertTodo(name, email, status string) *sql.Row {
	return Conn().QueryRow("INSERT INTO customer (name,email,status) values ($1,$2,$3) RETURNING id", name, email, status)
}

// func SelectByID(id int) *sql.Row {
//     //return Conn().Prepare("SELECT id, name, email ,status FROM customer WHERE id=$1")
// 	stmt, err := Conn().Prepare("SELECT id, name, email ,status FROM customer WHERE id=$1")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	row := stmt.QueryRow(id)
// 	t := Todo{}
// 	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	return row
// }
