package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getPeers() []YggdrasilIPAddress {
	resp, err := http.Get("https://raw.githubusercontent.com/threefoldtech/planetary_network/main/nodelist")
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error to fetch peers")
		fmt.Println("StatusCode: ", resp.StatusCode)
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error to read body")
		log.Fatal(err)
	}

	arrayAddresses := []string{}
	err = json.Unmarshal([]byte(string(body)), &arrayAddresses)

	if err != nil {
		fmt.Println("Error in parsing")
		log.Fatal(err)
	}

	var ipAddresses []YggdrasilIPAddress

	for _, s := range arrayAddresses {
		isThreefold := strings.HasPrefix(s, "tf-")
		fullAdd := strings.ReplaceAll(s, "tf-", "")
		result := strings.ReplaceAll(fullAdd, "tls://", "")
		result = strings.ReplaceAll(result, "tcp://", "")
		result = strings.ReplaceAll(result, "[", "")
		result = strings.ReplaceAll(result, "]", "")
		splitResult := strings.Split(result, ":")
		finalResult := strings.ReplaceAll(result, ":"+splitResult[len(splitResult)-1], "")

		ip := YggdrasilIPAddress{
			FullIPAddress:   fullAdd,
			IPAddress:       finalResult,
			latency:         9999,
			isThreefoldNode: isThreefold,
			RealIP:          finalResult,
		}

		ipAddresses = append(ipAddresses, ip)
	}

	return ipAddresses
}
