package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/gocolly/colly"
	"github.com/yggdrasil-network/yggdrasil-go/src/config"
	"github.com/yggdrasil-network/yggdrasil-go/src/defaults"
)

var publicYggdrasilPeersURL = "https://publicpeers.neilalexander.dev/"
var wg sync.WaitGroup

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func generateConfigFile(cfg *config.NodeConfig) {
	// if !fileExists("/etc/threefold_yggdrasil.conf") {
	fmt.Println("[info]: Generating config file")
	var configPeers string
	configPeers = <-getConfigPeers()
	fmt.Println("[info]: No config file")
	fmt.Println("Config file doesnt exist ...")
	cfg = defaults.GenerateConfig()
	fmt.Println(cfg)

	fmt.Println("[info]: config generated")

	configFile := doGenconf(true)
	fmt.Println("Config file created")
	fmt.Println("[info]: Created config file")
	configFile = strings.ReplaceAll(configFile, "\"Peers\": []", configPeers)
	fmt.Println("Peers replaced")
	fmt.Println("[info]: Peers replaced")

	f, err := os.Create("/etc/threefold_yggdrasil.conf")
	if err != nil {
		fmt.Println(err)
		fmt.Println("[err01]: " + err.Error())
		return
	}
	l, err := f.WriteString(configFile)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[err02]: " + err.Error())
		f.Close()
		return
	}
	fmt.Println("[info]: Config written")
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	// }
}

// Making this function async in some magic go-syntax land.
func getConfigPeers() <-chan string {
	fmt.Sprintf("Finding Peers ")
	r := make(chan string)

	go func() {
		defer close(r)

		c := colly.NewCollector()
		var ipAddresses []YggdrasilIPAddress

		c.OnHTML(".statusgood #address", func(e *colly.HTMLElement) {
			fmt.Print("Finding addresses")
			// Filtering out all tcp and ipv6 addresses

			if strings.Contains(e.Text, "tcp://") {
				fmt.Print("Skip tcp")
				return
			}

			if strings.Contains(e.Text, "[") && strings.Contains(e.Text, "]") {
				fmt.Print("Skip ipv6")
				return
			}

			// This also filters ipv6 incase we want it in the future.

			result := strings.ReplaceAll(e.Text, "tls://", "")
			result = strings.ReplaceAll(result, "tcp://", "")
			result = strings.ReplaceAll(result, "[", "")
			result = strings.ReplaceAll(result, "]", "")
			splitResult := strings.Split(result, ":")
			finalResult := strings.ReplaceAll(result, ":"+splitResult[len(splitResult)-1], "")
			fmt.Print(finalResult)
			ipAddr := YggdrasilIPAddress{
				FullIPAddress: e.Text,
				IPAddress:     finalResult,
				latency:       9999,
			}

			ipAddresses = append(ipAddresses, ipAddr)
		})

		c.Visit(publicYggdrasilPeersURL)

		for index := 0; index < len(ipAddresses); index++ {
			fmt.Print("found result")
			wg.Add(1)
			go pingAddress(ipAddresses[index])
		}

		wg.Wait()
		fmt.Print("wait done")
		sort.Slice(ipAddresses, func(i, j int) bool {
			return ipAddresses[i].latency < ipAddresses[j].latency
		})

		r <- fmt.Sprintf("\"Peers\": [\"%s\", \"%s\", \"%s\"]", ipAddresses[0].FullIPAddress, ipAddresses[1].FullIPAddress, ipAddresses[2].FullIPAddress)
	}()

	return r
}
func pingAddress(addr YggdrasilIPAddress) {
	fmt.Println("pinging")
	pinger, err := ping.NewPinger(addr.IPAddress)
	pinger.Timeout = time.Second / 2

	if err != nil {
		//panic(err)
	}
	pinger.Count = 2
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		// panic(err)
	}
	stats := pinger.Statistics() // get send/receive/rtt stats

	if stats.AvgRtt.String() == "0s" {
		fmt.Println("0s so skipped")
		addr.latency = 9999
		defer wg.Done()
		return
	}
	fmt.Println("Found host", addr.FullIPAddress)
	addr.latency, _ = strconv.ParseFloat(strings.ReplaceAll(stats.AvgRtt.String(), "ms", ""), 64)
	defer wg.Done()
}
