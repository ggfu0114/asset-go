package main

import (
	"asset-go/src/api"
	dbmodel "asset-go/src/models"
	assetserv "asset-go/src/services"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
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
		log.Println("Failed to parse token, ", err)
		return nil, err
	}

	// Checking token validity
	if !token.Valid {
		log.Println("invalid token")
		return nil, err
	}
	return &userClaim, nil
}

func CheckToken(c *gin.Context) (*UserClaim, error) {
	session := sessions.Default(c)
	token := session.Get("token")
	var claims *UserClaim
	var err error
	if token == nil {
		err = errors.New("token is empty")
		c.AbortWithStatus(http.StatusUnauthorized)

	} else {
		log.Println("token from session", token.(string))
		claims, err = ParseToken(token.(string))
		if err != nil {
			err = errors.New("failed to parse token")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
	return claims, err
}

func getUserFromContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Try to get user information by Middleware!")
		userClaim, err := CheckToken(c)
		if err != nil {
			log.Println("claims, err", userClaim, err)
			// Handle error, e.g., return unauthorized
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("userInfo", userClaim)
		c.Next()
	}
}

func getRequestUserInfo(c *gin.Context) (*UserClaim, bool) {
	env := os.Getenv("APP_ENV")
	if env == "dev" {
		user := UserClaim{
			Id:   "0",
			Name: "Sys dev",
		}
		return &user, true
	}
	user, err := c.MustGet("userInfo").(*UserClaim)
	if !err {
		c.AbortWithError(http.StatusInternalServerError, gin.Error{})
	}
	return user, err
}

func main() {
	env := os.Getenv("APP_ENV")
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	apiAuth := r.Group("/api/auth")
	store := cookie.NewStore([]byte("secret"))
	apiAuth.Use(sessions.Sessions("asset_sess", store))
	apiV1.Use(sessions.Sessions("asset_sess", store))

	apiAuth.GET("/oauth_login", func(c *gin.Context) {

		url := conf.AuthCodeURL("state")
		c.Redirect(http.StatusFound, url)
	})

	apiAuth.GET("/oauth_callback", func(c *gin.Context) {
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
		user, err := dbmodel.InsertUser(result)

		defer resp.Body.Close()

		token1 := jwt.NewWithClaims(jwt.SigningMethodHS256,
			UserClaim{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "PaulChen",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
				Id:   strconv.Itoa(user.Uid),
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

		c.Redirect(http.StatusFound, "/")
	})
	if env != "dev" {
		apiV1.Use(getUserFromContextMiddleware())
	}

	apiV1.GET("/assets", func(c *gin.Context) {
		user, err := getRequestUserInfo(c)
		if !err {
			c.AbortWithError(http.StatusInternalServerError, gin.Error{})
		}

		assets := dbmodel.ListMyAsset(user.Id)
		log.Println("assets", assets)

		c.JSON(200, assets)

	})
	apiV1.POST("/asset", func(c *gin.Context) {
		user, err := getRequestUserInfo(c)
		if !err {
			c.AbortWithError(http.StatusInternalServerError, gin.Error{})
		}
		reqAsset := new(dbmodel.Asset)
		err1 := c.Bind(reqAsset)
		if err1 != nil {
			log.Fatalln("Failed to bind request payload.", err)
		}
		aid := dbmodel.AddAsset(user.Id, *reqAsset)
		aa, err2 := strconv.Atoi(aid)
		reqAsset.Aid = aa
		if err2 != nil {
			log.Fatalln("Failed to convert asset id.", err2)
		}
		assetserv.QueryAssetValue(*reqAsset)
		c.JSON(200, gin.H{
			"aid": aid,
		})
	})
	apiV1.PUT("/asset/:aid", func(c *gin.Context) {

		user, userInfoErr := c.MustGet("userInfo").(*UserClaim)
		if !userInfoErr {
			c.AbortWithError(http.StatusInternalServerError, gin.Error{})
		}
		aid := c.Param("aid")
		reqAsset := new(dbmodel.Asset)
		err := c.Bind(reqAsset)
		if err != nil {
			log.Fatalln("Failed to bind request payload.", err)
		}
		dbmodel.UpdateAsset(user.Id, aid, *reqAsset)
		c.JSON(200, gin.H{
			"reqAsset": reqAsset,
		})
	})
	apiV1.DELETE("/asset/:aid", func(c *gin.Context) {
		user, userInfoErr := c.MustGet("userInfo").(*UserClaim)
		if !userInfoErr {
			c.AbortWithError(http.StatusInternalServerError, gin.Error{})
		}
		aid := c.Param("aid")
		dbmodel.DeleteAsset(user.Id, aid)
		c.JSON(200, gin.H{
			"deletedAsset": aid,
		})
	})

	r.Run()
}
