package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var info ConnectionInfo

func startServer() {
	fmt.Println("Running in server mode")
	router := gin.Default()
	router.GET("/connection", getConnection)
	router.POST("/connect", connect)
	router.POST("/disconnect", disconnect)
	router.GET("/health", health)
	router.POST("/exit", exit)
	router.POST("/reset", reset)
	router.Run("127.0.0.1:7070")
	return

}
func reset(c *gin.Context) {
	resetApplication()
	os.Exit(0)
}
func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, true)
}

func getConnection(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, info)
}

func connect(c *gin.Context) {
	info = yggdrasilConnect()
	c.IndentedJSON(http.StatusOK, info)
}
func disconnect(c *gin.Context) {
	n.shutdown()
}
func exit(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, "Shutting down")
	os.Exit(0)
}
