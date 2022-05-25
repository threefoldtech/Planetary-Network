package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	init_config()

	InitializeEvents()

	server := flag.Bool("server", false, "Yggdrasil root server")
	flag.Parse()

	if *server { //check if started as server
		startServer()

	} else {
		// This will cause the application to only have one instance. If the application is already running (eg listening on port 62854),
		// the application will not start another instance.
		http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 1
		_, err := http.Get("http://localhost:62854/raise")

		if err != nil { //only start if this is the single instance, in onther case the existing instance is raised
			fmt.Println("Err on raise")

			go startOneServer()

			args := getArgs()
			hup := make(chan os.Signal, 1)
			//signal.Notify(hup, os.Interrupt, syscall.SIGHUP)
			term := make(chan os.Signal, 1)
			signal.Notify(term, os.Interrupt, syscall.SIGTERM)
			for {
				done := make(chan struct{})
				ctx, cancel := context.WithCancel(context.Background())
				userInterface(args, ctx, done)
				select {
				case <-hup:
					cancel()
					<-done
				case <-term:
					cancel()
					<-done
					return
				case <-done:
					return
				}
			}
		}
	}
}
