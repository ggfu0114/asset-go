package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Name struct {
	FirstName  string `json:"first_name" binding:"required"`
	SecondName string `json:"second_name" binding:"required"`
}

func main() {
	// Hello world, the web server

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/input", func(c *gin.Context) {
		var name Name
		c.BindJSON(&name)
		log.Println("Recv one request", name)
		c.JSON(200, gin.H{"first": name.FirstName, "second": name.SecondName})
	})
	r.Run()

}
