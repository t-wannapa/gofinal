package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/t-wannapa/gofinal/customer"
	"github.com/t-wannapa/gofinal/middleware"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Auth)

	r.POST("/customers", customer.CreateHandler)
	r.GET("/customers/:id", customer.GetByIdHandler)
	r.GET("/customers", customer.GetAllHandler)
	r.PUT("/customers/:id", customer.UpdateHandler)
	r.DELETE("/customers/:id", customer.DeleteHandler)
	return r
}

func main() {
	fmt.Println("customer service")

	r := setupRouter()
	r.Run(":2009")
	//run port ":2009"
}
