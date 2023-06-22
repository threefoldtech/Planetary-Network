package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gologme/log"
)

func main() {
	// Add Logging (for some reasons, this cannot be in seperate function)
	f, err := os.OpenFile(GetThreefoldDirectory()+"tf-planetary-connector.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)

	if err != nil {
		log.Errorln("Error opening file: %v", err)
		panic(err)
	}

	defer f.Close()

	log.SetOutput(f)

	log.EnableLevel("debug")
	log.EnableLevel("info")
	log.EnableLevel("warn")
	log.EnableLevel("error")

	log.EnableFormattedPrefix()
	log.Infoln("STARTING APPLICATION")

	// Set right Ygg configs
	InitializeConfig()

	// Makes events available
	InitializeEvents()

	// Ping all IPs and store in global var
	fillPeers()

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
			go startOneServer()

			args := getArgs()
			hup := make(chan os.Signal, 1)
			//signal.Notify(hup, os.Interrupt, syscall.SIGHUP)
			term := make(chan os.Signal, 1)
			signal.Notify(term, os.Interrupt, syscall.SIGTERM)
			for {
				done := make(chan struct{})
				ctx, cancel := context.WithCancel(context.Background())
				//TODO: this function is blocking as it calls [QApplication::exec()](https://doc.qt.io/qt-5/qapplication.html#exec)
				// so I don't quite get how this is supposed to work
				userInterface(args, ctx, done)
				log.Info("User interface created")
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
