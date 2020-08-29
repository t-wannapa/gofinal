package customer

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateTable() {
	createTb := `CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	)`

	_, err := db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
}

func CreateHandler(c *gin.Context) {
	cust := Customer{}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	GenerateTable()

	row := db.QueryRow("INSERT INTO customers (name, email, status) VALUES ($1, $2, $3) RETURNING id", cust.Name, cust.Email, cust.Status)
	err := row.Scan(&cust.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, cust)
}

func GetByIdHandler(c *gin.Context) {
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)
	cust := &Customer{}
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func GetAllHandler(c *gin.Context) {
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	customers := []Customer{}
	for rows.Next() {
		cust := Customer{}
		err := rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		customers = append(customers, cust)
	}
	c.JSON(http.StatusOK, customers)
}

func UpdateHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cust := Customer{}
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	stmt, err = db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if _, err := stmt.Exec(id, cust.Name, cust.Email, cust.Status); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, cust)

}

func DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}
