package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"asset-go/src/models"
	_ "github.com/go-sql-driver/mysql"
)



func main() {
	// Hello world, the web server

	r := gin.Default()
	
	r.GET("/assets", func(c *gin.Context) {
		assets := dbmodel.ListMyAsset()
		log.Println("assets", assets)
		
		c.JSON(200, gin.H{
			"assets": assets,
		})

	})
	r.POST("/asset", func(c *gin.Context) {
		reqAsset := new(dbmodel.Asset)
		err := c.Bind(reqAsset)
		if err != nil{
			log.Fatalln("Failed to bind request payload.", err)
		}
		aid := dbmodel.AddAsset(*reqAsset)
		c.JSON(200, gin.H{
			"aid": aid,
		})
	})
	r.PUT("/asset/:aid", func(c *gin.Context) {
		aid := c.Param("aid")
		reqAsset := new(dbmodel.Asset)
		err := c.Bind(reqAsset)
		if err != nil{
			log.Fatalln("Failed to bind request payload.", err)
		}
		dbmodel.UpdateAsset(aid, *reqAsset)
		c.JSON(200, gin.H{
			"reqAsset": reqAsset,
		})
	})
	r.DELETE("/asset/:aid", func(c *gin.Context) {
		aid := c.Param("aid")
		dbmodel.DeleteAsset(aid)
		c.JSON(200, gin.H{
			"deletedAsset": aid,
		})
	})
	
	r.Run()
}
