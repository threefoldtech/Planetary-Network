package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var info ConnectionInfo

func startServer() {
	info = ConnectionInfo{
		IpAddress:       "",
		SubnetAddress:   "",
		PublicKey:       "",
		ConnectionPeers: []string{},
		Error:           "",
	}

	fmt.Println("Running in server mode")
	router := gin.Default()
	router.GET("/info", getInfo)
	router.POST("/connect", connect)
	router.POST("/disconnect", disconnect)
	router.GET("/health", health)
	router.POST("/exit", exit)
	router.POST("/reset", reset)
	router.POST("/delete", delete)
	router.Run("127.0.0.1:62853")
	return

}
func reset(c *gin.Context) {
	resetApplication()
	os.Exit(0)
}
func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, true)
}

func getInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getConnectionInfo())
}

func connect(c *gin.Context) {
	// addCronIfNotExists()

	info = yggdrasilConnect()
	c.IndentedJSON(http.StatusOK, info)
}
func disconnect(c *gin.Context) {
	info = yggdrasilDisconnect()

}
func exit(c *gin.Context) {
	if info.IpAddress != "" {
		n.shutdown()
	}
	c.IndentedJSON(http.StatusOK, "Shutting down")
	os.Exit(0)
}

func delete(c *gin.Context) {
	deleteConfig()
}
