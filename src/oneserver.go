//oneserver provides functionality to raise window if the application is already running

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gologme/log"
)

func startOneServer() {
	log.Infoln("RUNNING SERVER IN SINGLE INSTANCE MODE")

	router := gin.Default()
	router.GET("/raise", raise)
	router.Run("127.0.0.1:62854")
	return

}
func raise(c *gin.Context) {
	log.Infoln("RAISING WINDOW")

	raiseWindow()
}
