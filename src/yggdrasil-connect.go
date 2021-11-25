package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gologme/log"

	"github.com/yggdrasil-network/yggdrasil-go/src/admin"
	"github.com/yggdrasil-network/yggdrasil-go/src/config"

	"github.com/yggdrasil-network/yggdrasil-go/src/ipv6rwc"
	"github.com/yggdrasil-network/yggdrasil-go/src/multicast"
	"github.com/yggdrasil-network/yggdrasil-go/src/tuntap"
)

var n node

func yggdrasilConnect() ConnectionInfo {

	// defer close(done)

	if _, err := os.Stat("/var/run/yggdrasil.sock"); err == nil {
		connInfo := ConnectionInfo{
			IpAddress:     "N/A",
			SubnetAddress: "N/A",
			PublicKey:     "N/A",
			Error:         "Yggdrasil is already running",
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

	if !fileExists("/etc/threefold_yggdrasil.conf") {
		generateConfigFile(cfg)
	}
	cfg = readConfig(logger, true, "/etc/threefold_yggdrasil.conf", false)

	logger.Errorln("An error occurred during startup")
	fmt.Println("Private key in config is ", cfg.PrivateKey)

	n = node{config: cfg}

	// Now start Yggdrasil - this starts the DHT, router, switch and other core
	// components needed for Yggdrasil to operate
	if err = n.core.Start(cfg, logger); err != nil {
		logger.Errorln("An error occurred during startup")
		panic(err)
	}
	// Register the session firewall gatekeeper function
	// Allocate our modules
	n.admin = &admin.AdminSocket{}
	n.multicast = &multicast.Multicast{}
	n.tuntap = &tuntap.TunAdapter{}
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
	n.tuntap.SetupAdminHandlers(n.admin)
	// Make some nice output that tells us what our IPv6 address and subnet are.
	// This is just logged to stdout for the user.
	address := n.core.Address()
	subnet := n.core.Subnet()
	public := n.core.GetSelf().Key
	logger.Infof("Your public key is %s", hex.EncodeToString(public[:]))
	logger.Infof("Your IPv6 address is %s", address.String())
	logger.Infof("Your IPv6 subnet is %s", subnet.String())

	connInfo := ConnectionInfo{
		IpAddress:     address.String(),
		SubnetAddress: subnet.String(),
		PublicKey:     hex.EncodeToString(public[:]),
		Error:         "",
	}
	return connInfo

}

func resetApplication() {
	err := os.Remove("/var/run/yggdrasil.sock")

	if err != nil {
		fmt.Println(err)
	}

	err = os.Remove("/etc/threefold_yggdrasil.conf")

	if err != nil {
		fmt.Println(err)
	}
}
