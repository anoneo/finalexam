package customerhandler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/anoneo/finalexam/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var id int = 0
var todos []Todo
var todosMap = make(map[int]*Todo)

func deleteTodosHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}

	stmt, err := database.Conn().Prepare("DELETE FROM customer WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err = stmt.Exec(pId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func updateTodosHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	var temp Todo
	err = c.ShouldBindJSON(&temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("DEBUG>>>>", temp)
	stmt, err := database.Conn().Prepare("UPDATE customer SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(pId, &temp.Name, &temp.Email, &temp.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	temp.ID = pId
	c.JSON(http.StatusOK, temp)
}

func getTodosHandler(c *gin.Context) {
	pStatus := c.Query("status")
	var temp []Todo
	sql := "SELECT id, name, email, status FROM customer"
	if pStatus != "" {
		sql += " WHERE status=$1"
	}

	stmt, err := database.Conn().Prepare(sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := stmt.Query()
	if pStatus != "" {
		rows, err = stmt.Query(pStatus)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			log.Fatal("can't scan : ", err)
		}
		temp = append(temp, t)
	}
	c.JSON(http.StatusOK, temp)
}

func getTodosByIdHandler(c *gin.Context) { // GET:id
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	stmt, err := database.Conn().Prepare("SELECT id, name, email ,status FROM customer WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	row := stmt.QueryRow(pId)
	t := Todo{}
	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, t)

}

func createTodosHandler(c *gin.Context) { // POST
	var item Todo
	err := c.ShouldBindJSON(&item)
	fmt.Println("ITEM =", item.Name, item.Email, item.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// row := database.Conn().QueryRow("INSERT INTO todos (title,status) values ($1,$2) RETURNING id", item.Title, "active")
	row := database.InsertTodo(item.Name, item.Email, item.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	item.ID = id
	c.JSON(http.StatusCreated, item)
}

func CreateTb() {
	var err error
	createTb := `
 CREATE TABLE IF NOT EXISTS customer (
  id SERIAL PRIMARY KEY,
  name TEXT,
  email TEXT,
  status TEXT
 );
 `
	_, err = database.Conn().Exec(createTb)
	if err != nil {
		log.Fatal("can't create table : ", err)
	}
}

func loginMiddleware(c *gin.Context) {
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(loginMiddleware)
	r.GET("/customers", getTodosHandler)
	r.GET("/customers/:id", getTodosByIdHandler)
	r.POST("/customers", createTodosHandler)
	r.PUT("/customers/:id", updateTodosHandler)
	r.DELETE("/customers/:id", deleteTodosHandler)
	return r
}
