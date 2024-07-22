package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"asset-go/src/models"
	_ "github.com/go-sql-driver/mysql"
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
	r.GET("/assets", func(c *gin.Context) {
		assets := dbmodel.ListMyAsset()
		log.Println("assets", assets)
		
		c.JSON(200, gin.H{
			"assets": assets,
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
