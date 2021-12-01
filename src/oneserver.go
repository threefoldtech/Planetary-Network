//oneserver provides functionality to raise window if the application is already running

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func startOneServer() {
	fmt.Println("Running in the single instance server")
	router := gin.Default()
	router.GET("/raise", raise)
	router.Run("127.0.0.1:62854")
	return

}
func raise(c *gin.Context) {
	raiseWindow()
}
