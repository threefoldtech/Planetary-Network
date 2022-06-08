package main

import (
	"net/http"
	"os"
	"time"
)

func checkToStartNetworkServer() {

	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 1
	_, err := http.Get("http://localhost:62853/health")
	if err != nil {
		startNetworkServer()
	}
}

func getUsername() string {
	return os.Getenv("USER")
}
