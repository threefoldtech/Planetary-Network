package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/robfig/cron/v3"
)

var cr *cron.Cron = cron.New()

func addCronIfNotExists() {
	if len(cr.Entries()) <= 0 {
		cr.AddFunc("@every 10s", func() {
			info := fetchConnectionData()
			fmt.Println(info)
			// UpdatePeersUi(info)
		})
		cr.Start()
	}
}

func removeCrons() {
	entries := cr.Entries()
	for _, entry := range entries {
		cr.Remove(entry.ID)
	}
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
		fmt.Errorf("Error in ioutil", err)
	}

	if body == nil {
		fmt.Println("BODY IS NULL")
	}

	var details *ConnectionDetails
	err = json.Unmarshal([]byte(string(body)), &details)

	if err != nil {
		fmt.Errorf("Error in unmarshall", err)
	}

	return details
}
