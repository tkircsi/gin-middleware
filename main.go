package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck is used for container HEALTHCHECK
// GET /healthcheck
// Response: "OK"
func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// HomeAPI returns the Demo API version number
// GET /
func HomeAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "version: 0.1",
	})
}

// CheckAppAPIKey checks API_KEY
// GET /api
func CheckAppAPIKey(c *gin.Context) {
	const demoKey string = "11AA22BB"
	apiKey := c.GetHeader("API_KEY")
	if len(apiKey) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "API_KEY is required!",
		})
		return
	}

	if apiKey != demoKey {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid API_KEY",
		})
		return
	}

	c.Next()
}

// GetStatus function
// GET /api/status
func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome at Demo API!",
	})
}

// GetTime function returns the current time at location
// Closure funtion with location initial value
// GET /api/time
func GetTime(location string) gin.HandlerFunc {
	return func(c *gin.Context) {
		loc, err := time.LoadLocation(location)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Now()
		t = t.In(loc)
		c.JSON(http.StatusOK, gin.H{
			"location": location,
			"time":     t.String(),
		})
	}
}

func Auth(c *gin.Context) {
	token := c.GetHeader("Bearer")
	if len(token) == 0 {
		c.Redirect(http.StatusTemporaryRedirect, "https://tkircsigmailb2c.b2clogin.com/tkircsigmailb2c.onmicrosoft.com/oauth2/v2.0/authorize?p=B2C_1_signupsignin1&client_id=7f8c9fc4-5d4e-4673-9985-2d0c157191e7&nonce=defaultNonce&redirect_uri=http%3A%2F%2Flocalhost%3A5000%2Fcallback&scope=openid&response_type=id_token&&response_mode=query&prompt=login")
	}
}

func CallBack(c *gin.Context) {
	idToken := c.Query("id_token")
	
	fmt.Printf("id_token: %s\n", idToken)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func main() {
	r := gin.Default()
	r.GET("/", HomeAPI)
	r.GET("/healthcheck", HealthCheck)
	r.GET("/userinfo", Auth)
	r.GET("/callback", CallBack)
	api := r.Group("/api")
	api.Use(CheckAppAPIKey)
	{
		api.GET("/status", GetStatus)
		api.GET("/time", GetTime("Europe/Budapest"))
	}
	r.Run(":5000")
}
