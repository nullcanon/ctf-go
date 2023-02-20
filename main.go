package main

import (
	"ctf/api"
	"ctf/blockscan"
	"ctf/config"
	"ctf/core"
	"ctf/jsonrpc"
	"ctf/middleware/Cors"
	"ctf/models"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()

	r := gin.Default()
	r.Use(Cors.Cors())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	var err error
	core.InviterHandle, err = core.NewInviter()
	go blockscan.ScanTradeVolume()
	go blockscan.ScanLpRewards()
	if err != nil {
		fmt.Println("inviterHandle init error", err.Error())
	}
	// core.InviterHandle.ProcessPresellUsersRewards()

	apiHandle := api.Api{}

	r.POST("/", func(c *gin.Context) { jsonrpc.ProcessJsonRPC(c, &apiHandle) })

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	config.InitConfig()
	models.Setup()

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8888")
}
