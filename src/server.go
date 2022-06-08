package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gologme/log"
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

	log.Infoln("RUNNING IN SERVER MODE")

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
	log.Infoln("API: RESET")

	resetApplication()
	os.Exit(0)
}
func health(c *gin.Context) {
	log.Infoln("API: HEALTH")

	c.IndentedJSON(http.StatusOK, true)
}

func getInfo(c *gin.Context) {
	log.Infoln("API: INFO")

	c.IndentedJSON(http.StatusOK, getConnectionInfo())
}

func connect(c *gin.Context) {
	log.Infoln("API: CONNECT")

	info = yggdrasilConnect()
	c.IndentedJSON(http.StatusOK, info)
}

func disconnect(c *gin.Context) {
	log.Infoln("API: DISCONNECT")

	info = yggdrasilDisconnect()
}

func exit(c *gin.Context) {
	log.Infoln("API: EXIT")

	if info.IpAddress != "" {
		n.shutdown()
	}
	c.IndentedJSON(http.StatusOK, "Shutting down")
	os.Exit(0)
}

func delete(c *gin.Context) {
	log.Infoln("API: DELETE")

	DeleteConfig()
}
