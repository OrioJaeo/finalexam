package customer

import (
	"log"
	"net/http"

	"github.com/OrioJaeo/finalexam/database"
	"github.com/gin-gonic/gin"
)

func CreateTable() {
	ctb := `
	 CREATE TABLE IF NOT EXISTS customers (
		 id SERIAL PRIMARY KEY,
		 name TEXT,
		 email TEXT,
		 status TEXT
	 );`
	_, err := database.Conn().Exec(ctb)
	if err != nil {
		log.Fatal("can't create table ", err)
	}
}

func NewRouter() *gin.Engine {
	//database.Conn().Close()
	r := gin.Default()
	r.Use(loginMiddleWare)
	r.POST("/customers", createCustomersHandler)
	r.DELETE("/customers/:id", delCustByIdHander)
	r.PUT("/customers/:id", putCustByIdHander)
	r.GET("/customers/:id", getCustByIdHandler)
	r.GET("/customers", getCustAllHandler)

	return r
}
func loginMiddleWare(c *gin.Context) {

	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "UnAuthorize")
		c.Abort()
		return
	}
	c.Next()
}

//add record
func createCustomersHandler(c *gin.Context) {
	var item Customer
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	row := database.Conn().QueryRow("INSERT INTO Customers(name, email, status) values($1,$2,$3) RETURNING id", item.Name, item.Email, item.Status)

	err = row.Scan(&item.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't Scan row into variable" + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

//Delete
func delCustByIdHander(c *gin.Context) {
	id := c.Param("id")

	stmt, err := database.Conn().Prepare("DELETE FROM Customers WHERE ID=$1;")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if _, err := stmt.Exec(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "data not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
	return
}

//Update
func putCustByIdHander(c *gin.Context) {
	item := Customer{}
	id := c.Param("id")
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	stmt, err := database.Conn().Prepare("UPDATE Customers SET name=$2, email=$3, status=$4 WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if _, err := stmt.Exec(id, item.Name, item.Email, item.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, item)
	return
}

// get data by id
func getCustByIdHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM Customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	row := stmt.QueryRow(id)
	t := Customer{}
	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": t})
		return
	}
	c.JSON(http.StatusOK, t)
}

//get all

func getCustAllHandler(c *gin.Context) {

	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM Customers ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't query all todos" + err.Error()})
		return
	}
	items := []Customer{}
	for rows.Next() {
		t := Customer{}
		err = rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "can't Scan row into variable" + err.Error()})
			return
		}
		items = append(items, t)
	}

	c.JSON(http.StatusOK, items)
	return

}
