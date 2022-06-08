package main

import (
	"encoding/json"
	"github.com/gologme/log"
	"io/ioutil"
	"net/http"
	"strings"
)

func getPeers() []YggdrasilIPAddress {
	resp, err := http.Get("https://raw.githubusercontent.com/threefoldtech/planetary_network/main/nodelist")
	if err != nil || resp.StatusCode != 200 {
		log.Errorln("ERROR TO FETCH PEERS FROM GITHUB", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("ERROR IN GET PEERS READ BODY")
	}

	arrayAddresses := []string{}
	err = json.Unmarshal([]byte(string(body)), &arrayAddresses)

	if err != nil {
		log.Errorln("ERROR IN GET PEERS UNMARSHALL")
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
