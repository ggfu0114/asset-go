package main

import (
	api "asset-go/src/api"
	dbmodel "asset-go/src/models"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserClaim struct {
	jwt.RegisteredClaims
	Id   string `json:"id"`
	Name string `json:"name"`
}

const key = "<key>"

var conf = &oauth2.Config{
	ClientID:     "<client-id>",
	ClientSecret: "<client-security>",
	RedirectURL:  "http://127.0.0.1:3000/oauth_callback",
	Scopes:       []string{"profile", "email"},
	Endpoint:     google.Endpoint,
}

func ParseToken(jwtToken string) (*UserClaim, error) {
	var userClaim UserClaim
	token, err := jwt.ParseWithClaims(jwtToken, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Checking token validity
	if !token.Valid {
		log.Fatal("invalid token")
		return nil, err
	}
	return &userClaim, nil
}

func CheckToken(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("token")
	url := conf.AuthCodeURL("state")
	if token == nil {
		log.Println("Session token is empty.")

		c.Redirect(http.StatusFound, url)
	} else {
		log.Println("token from session", token.(string))
		claims, err := ParseToken(token.(string))
		if err != nil {
			log.Println("errerrerr", err)
			c.Redirect(http.StatusFound, url)
		} else {
			log.Println("claims", claims)
		}
	}
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
		client := conf.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		var result api.UserInfo
		if err := json.Unmarshal(data, &result); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		log.Println("result", result)
		user, err := dbmodel.InsertUser(result)
		log.Println("user", user)
		log.Println("err", err)

		defer resp.Body.Close()

		token1 := jwt.NewWithClaims(jwt.SigningMethodHS256,
			UserClaim{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "PaulChen",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
				Id:   strconv.FormatInt(user.Uid, 10),
				Name: user.Name,
			})
		jwtToken, err := token1.SignedString([]byte(key))

		log.Println("generated session", jwtToken)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		session := sessions.Default(c)
		session.Set("token", jwtToken)
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
