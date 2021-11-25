package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := flag.Bool("server", false, "Yggdrasil root server")
	flag.Parse()

	if *server { //check if started as server
		startServer()
	} else {

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
