package main

import (
	dbmodel "asset-go/src/models"
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf = &oauth2.Config{
	ClientID:     "<ClientID>",
	ClientSecret: "<ClientSecret>",
	RedirectURL:  "http://127.0.0.1:3000/oauth_callback",
	Scopes:       []string{"profile", "email"},
	Endpoint:     google.Endpoint,
}

func CheckToken(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		log.Println("Session token is empty.")
		url := conf.AuthCodeURL("state")
		c.Redirect(http.StatusFound, url)
	}
	log.Println("Session ID:", token)
}

func main() {

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/oauth_callback", func(c *gin.Context) {
		state := c.Query("state")
		if state != "state" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("invalid csrf token"))
			return
		}
		code := c.Query("code")
		var err error
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		log.Println("Google Oauth token", token)

		session := sessions.Default(c)
		session.Set("token", "pass")
		session.Save()
		c.Redirect(http.StatusFound, "/assets")
	})
	r.Use(CheckToken)
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
		if err != nil {
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
		if err != nil {
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
