package main

import (
	"bytes"
	"encoding/json"
	"github.com/gologme/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SendPostAsync(url string, rc chan *http.Response) {

	response, err := http.Post(url, "application/json", bytes.NewBuffer(nil))
	if err != nil {
		panic(err)
	}
	rc <- response
}

func SendGetAsync(url string, rc chan *http.Response) {

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	rc <- response
}

func RetrievePeers() {
	// Default time.sleep otherwise peers are not initialized yet
	time.Sleep(time.Second)

	go func() {
		peersChannel := make(chan *http.Response)

		go SendGetAsync("http://localhost:62853/info", peersChannel)

		peersResponse := <-peersChannel
		defer peersResponse.Body.Close()

		body, err := ioutil.ReadAll(peersResponse.Body)
		if err != nil {
			log.Errorln("ERROR IN RETREIVE PEERS IN READALL", err)
		}

		if body == nil {
			log.Errorln("ERROR IN RETREIVE PEERS BODY IS NULL")
		}

		var details = &ConnectionDetails{}
		err = json.Unmarshal([]byte(string(body)), details)

		if err != nil {
			log.Errorln("ERROR IN RETREIVE PEERS IN UNMARSHALL", err)
		}

		Emit("PEERS_RECEIVED", *details)
	}()
}

func ConnectToServer() {
	go func() {
		connectionChannel := make(chan *http.Response)

		http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 15
		go SendPostAsync("http://localhost:62853/connect", connectionChannel)

		connectionResponse := <-connectionChannel
		defer connectionResponse.Body.Close()

		body, _ := ioutil.ReadAll(connectionResponse.Body)
		var connectionInfo = &ConnectionInfo{}
		json.Unmarshal([]byte(body), connectionInfo)

		Emit("EVENT_CONNECTED", *connectionInfo)
	}()
}

func DisconnectToServer() {
	http.Post("http://localhost:62853/disconnect", "application/json", bytes.NewBuffer(nil))

	// Catch interrupts from the operating system to exit gracefully.
	c := make(chan os.Signal, 1)
	r := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(r, os.Interrupt, syscall.SIGHUP)
}

func DeleteConfigOnFileSystem() {
	http.Post("http://localhost:62853/delete", "application/json", bytes.NewBuffer(nil))
}

func GetCurrentPeers() []string {
	resp, err := http.Get("http://localhost:62853/info")
	if err != nil {
		panic(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("ERROR IN  GET CURRENT PEERS IOUTIL", err)
	}

	var details *ConnectionDetails
	err = json.Unmarshal([]byte(string(body)), &details)

	if err != nil {
		log.Errorln("ERROR IN  GET CURRENT PEERS UNMARSHALL", err)
	}

	return details.ConnectionPeers
}

func fetchConnectionData() *ConnectionDetails {
	resp, err := http.Get("http://localhost:62853/info")
	if err != nil {
		panic(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("ERROR IN FETCH CONNECTION DATA IOUTIL", err)
	}

	if body == nil {
		log.Errorln("ERROR IN FETCH CONNECTION DATA BODY NULL", err)
	}

	var details *ConnectionDetails
	err = json.Unmarshal([]byte(string(body)), &details)

	if err != nil {
		log.Errorln("ERROR IN FETCH CONNECTION DATA UNMARSHALL", err)
	}

	return details
}
