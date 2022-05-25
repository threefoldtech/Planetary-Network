package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gologme/log"

	"github.com/yggdrasil-network/yggdrasil-go/src/admin"
	"github.com/yggdrasil-network/yggdrasil-go/src/config"
	"github.com/yggdrasil-network/yggdrasil-go/src/core"

	"github.com/yggdrasil-network/yggdrasil-go/src/ipv6rwc"
	"github.com/yggdrasil-network/yggdrasil-go/src/multicast"
	"github.com/yggdrasil-network/yggdrasil-go/src/tuntap"
)

var n node

func yggdrasilDisconnect() ConnectionInfo {
	n.shutdown()
	fmt.Println("COMING HERE 1")
	connInfo := ConnectionInfo{
		IpAddress:       "",
		SubnetAddress:   "",
		ConnectionPeers: []string{},
		PublicKey:       "",
		Error:           "",
	}
	return connInfo
}

func cleanupYggdrasilSock() string {

	cmd := "rm -rf /var/run/yggdrasil.sock"
	// fmt.Println(cmd)
	stdout, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err)
		// os.Exit(1)
	}

	return strings.TrimSpace(string(stdout))
}

func yggdrasilConnect() ConnectionInfo {
	cleanupYggdrasilSock()
	if _, err := os.Stat("/var/run/yggdrasil.sock"); err == nil {
		connInfo := ConnectionInfo{
			IpAddress:       "N/A",
			SubnetAddress:   "N/A",
			PublicKey:       "N/A",
			ConnectionPeers: []string{},
			Error:           "Yggdrasil is already running",
		}
		return connInfo

	}

	var logger *log.Logger

	logger = log.New(os.Stdout, "", log.Flags())

	if logger == nil {
		logger = log.New(os.Stdout, "", log.Flags())
		logger.Warnln("Logging defaulting to stdout")
	}
	logger.EnableLevel("error")
	logger.EnableLevel("info")
	var cfg *config.NodeConfig
	var err error

	if !fileExists(app_config.yggdrasil_config_location) {
		generateConfigFile(cfg)
	}
	cfg = readConfig(logger, true, app_config.yggdrasil_config_location, false)

	logger.Errorln("An error occurred during startup")
	fmt.Println("Private key in config is ", cfg.PrivateKey)

	n = node{config: cfg}

	// Now start Yggdrasil - this starts the DHT, router, switch and other core
	// components needed for Yggdrasil to operate
	if err = n.core.Start(cfg, logger); err != nil {
		logger.Errorln("An error occurred during startup")
		panic(err)
	}
	fmt.Println("COMING HERE 4")
	// Register the session firewall gatekeeper function
	// Allocate our modules
	n.admin = &admin.AdminSocket{}
	n.multicast = &multicast.Multicast{}
	n.tuntap = &tuntap.TunAdapter{}

	fmt.Println("COMING HERE 5")
	// Start the admin socket
	if err := n.admin.Init(&n.core, cfg, logger, nil); err != nil {
		logger.Errorln("An error occurred initialising admin socket:", err)
	} else if err := n.admin.Start(); err != nil {
		logger.Errorln("An error occurred starting admin socket:", err)
	}
	n.admin.SetupAdminHandlers(n.admin)
	// Start the multicast interface
	if err := n.multicast.Init(&n.core, cfg, logger, nil); err != nil {
		logger.Errorln("An error occurred initialising multicast:", err)
	} else if err := n.multicast.Start(); err != nil {
		logger.Errorln("An error occurred starting multicast:", err)
	}
	n.multicast.SetupAdminHandlers(n.admin)
	// Start the TUN/TAP interface
	rwc := ipv6rwc.NewReadWriteCloser(&n.core)
	if err := n.tuntap.Init(rwc, cfg, logger, nil); err != nil {
		logger.Errorln("An error occurred initialising TUN/TAP:", err)
	} else if err := n.tuntap.Start(); err != nil {
		logger.Errorln("An error occurred starting TUN/TAP:", err)
	}

	fmt.Println("COMING HERE 6")
	n.tuntap.SetupAdminHandlers(n.admin)
	// Make some nice output that tells us what our IPv6 address and subnet are.
	// This is just logged to stdout for the user.
	address := n.core.Address()
	subnet := n.core.Subnet()
	public := n.core.GetSelf().Key

	fmt.Println(n.core.GetPaths())

	var listPeers []string

	peers, err := getYggdrasilPeers()
	if err != nil {
		info.Error = err.Error()
	}

	for _, peer := range peers {
		address := peer.Remote
		listPeers = append(listPeers, string(address))
	}

	logger.Infof("Your public key is %s", hex.EncodeToString(public[:]))
	logger.Infof("Your IPv6 address is %s", address.String())
	logger.Infof("Your IPv6 subnet is %s", subnet.String())
	connInfo := ConnectionInfo{
		IpAddress:       address.String(),
		SubnetAddress:   subnet.String(),
		ConnectionPeers: listPeers,
		PublicKey:       hex.EncodeToString(public[:]),
		Error:           "",
	}
	return connInfo

}

func resetApplication() {
	cleanupYggdrasilSock()
	err := os.Remove("/var/run/yggdrasil.sock")

	if err != nil {
		fmt.Println(err)
	}

	err = os.Remove(app_config.yggdrasil_config_location)

	if err != nil {
		fmt.Println(err)
	}
}

func getConnectionInfo() (data ConnectionDetails) {

	var info ConnectionDetails

	ipAddress, err := getAddress()
	if err != nil {
		fmt.Println(err)
		info.Error = err.Error()
		return info
	}
	info.IpAddress = ipAddress

	var listPeers []string

	peers, err := getYggdrasilPeers()
	if err != nil {
		info.Error = err.Error()
	}

	for _, peer := range peers {
		address := peer.Remote
		listPeers = append(listPeers, string(address))
	}

	info.ConnectionPeers = listPeers
	return info

}

func getAddress() (address string, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("No address found")
		}
	}()

	address = n.core.Address().String()
	return address, nil
}

func getYggdrasilPeers() (peers []core.Peer, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("No peers found")
		}
	}()

	peers = n.core.GetPeers()
	return peers, nil
}

func deleteConfig() {
	cmd := "rm -rf /etc/threefold_yggdrasil.conf"

	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error in deleting config")
		fmt.Println(err)
		// os.Exit(1)
	}
}
