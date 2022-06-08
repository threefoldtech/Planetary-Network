package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/gologme/log"
	"github.com/yggdrasil-network/yggdrasil-go/src/config"
	"github.com/yggdrasil-network/yggdrasil-go/src/defaults"
)

var wg sync.WaitGroup
var ipAddresses []YggdrasilIPAddress

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func generateConfigFile(cfg *config.NodeConfig) {
	log.Infoln("GETTING CONFIG PEERS")
	var configPeers string
	configPeers = <-getConfigPeers()

	log.Infoln("GENERATING CONFIG FILE")
	cfg = defaults.GenerateConfig()
	configFile := doGenconf(true)
	log.Infoln("CONFIG FILE GENERATED")

	log.Infoln("FILLING IN PEERS")
	configFile = strings.ReplaceAll(configFile, "\"Peers\": []", configPeers)
	log.Infoln("PEERS REPLACED")

	f, err := os.Create(APPLICATION_CONFIG.yggdrasil_config_location)
	if err != nil {
		log.Errorln("ERROR CREATING FILE ON SYSTEM " + err.Error())
		return
	}

	l, err := f.WriteString(configFile)
	if err != nil {
		log.Errorln("ERROR WRITING TEXT INSIDE CONFIG " + err.Error())
		f.Close()
		return
	}

	log.Infoln("CONFIG SUCCESFULLY WRITTEN")
	log.Infoln(l)

	err = f.Close()
	if err != nil {
		log.Errorln("ERROR IN CLOSING FILE", err.Error())
		return
	}
}

func getConfigPeers() <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)

		ipAddresses = getPeers()

		for index := 0; index < len(ipAddresses); index++ {
			wg.Add(1)
			go getAddressInfo(ipAddresses[index], index)
		}

		wg.Wait()
		sort.Slice(ipAddresses, func(i, j int) bool {
			return ipAddresses[i].latency < ipAddresses[j].latency
		})

		var threefoldNodes []YggdrasilIPAddress
		var publicNodes []YggdrasilIPAddress

		for _, s := range ipAddresses {
			if s.isThreefoldNode {
				threefoldNodes = append(threefoldNodes, s)
			} else {
				publicNodes = append(publicNodes, s)
			}
		}

		r <- fmt.Sprintf("\"Peers\": [\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\"]", threefoldNodes[0].FullIPAddress, threefoldNodes[1].FullIPAddress, threefoldNodes[2].FullIPAddress, threefoldNodes[3].FullIPAddress, threefoldNodes[4].FullIPAddress, publicNodes[0].FullIPAddress, publicNodes[1].FullIPAddress, publicNodes[2].FullIPAddress, publicNodes[3].FullIPAddress, publicNodes[4].FullIPAddress)
	}()

	return r
}

// Making this function async in some magic go-syntax land.
func getPeerStats() <-chan []YggdrasilIPAddress {
	log.Infoln("DIDN'T GET PEERS YET, GETTING PEERS")
	r := make(chan []YggdrasilIPAddress)

	go func() {
		defer close(r)

		ipAddresses = getPeers()

		for index := 0; index < len(ipAddresses); index++ {
			wg.Add(1)
			go getAddressInfo(ipAddresses[index], index)
		}

		wg.Wait()
		sort.Slice(ipAddresses, func(i, j int) bool {
			return ipAddresses[i].latency < ipAddresses[j].latency
		})

		r <- ipAddresses
	}()

	return r
}

func fillPeers() {
	go func() {
		ipAddresses = getPeers()

		for index := 0; index < len(ipAddresses); index++ {
			wg.Add(1)
			go getAddressInfo(ipAddresses[index], index)
		}

		wg.Wait()
		sort.Slice(ipAddresses, func(i, j int) bool {
			return ipAddresses[i].latency < ipAddresses[j].latency
		})

	}()
}

func getAddressInfo(addr YggdrasilIPAddress, index int) {
	defer wg.Done()

	pinger, _ := ping.NewPinger(addr.IPAddress)

	if runtime.GOOS == WINDOWS {
		pinger.SetPrivileged(true)
	}

	pinger.Timeout = time.Second / 2

	pinger.Count = 2
	pinger.Run()

	stats := pinger.Statistics()
	if stats.AvgRtt.String() == "0s" {
		addr.latency = 9999
		ipAddresses[index] = addr
		log.Warnln("PING RESULT:", addr, " (0s)")
		return
	}

	addr.RealIP = stats.IPAddr.IP.String()
	addr.latency, _ = strconv.ParseFloat(strings.ReplaceAll(stats.AvgRtt.String(), "ms", ""), 64)

	log.Infoln("PING RESULT:", addr)
	ipAddresses[index] = addr
}
