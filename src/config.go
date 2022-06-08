package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/gologme/log"
)

type Config struct {
	yggdrasil_config_location string
	SubnetAddress             string
	PublicKey                 string
	Error                     string
}

var APPLICATION_CONFIG Config

func InitializeConfig() {
	dir := GetCurrentDirectory() + "/threefold_yggdrasil.conf"
	APPLICATION_CONFIG.yggdrasil_config_location = dir

	log.Infoln("SETTING YGGDRASIL CONFIG ON LOCATION:", APPLICATION_CONFIG.yggdrasil_config_location)
}

func DeleteConfig() {
	err := os.Remove(APPLICATION_CONFIG.yggdrasil_config_location)

	if err != nil {
		log.Errorln("ERROR IN REMOVING YGG CONFIG")
	}

}

func CleanYggSockets() {
	if runtime.GOOS == WINDOWS {
		log.Infoln("CLEANED SOCKETS: WINDOWS")
		return
	}

	cmd := "rm -rf /var/run/yggdrasil.sock"
	stdout, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Errorln("ERROR IN CLEANING UP YGGDRASIL SOCKETS", err.Error())
	}

	log.Infoln("CLEANED UP YGGDRASIL SOCKETS", string(stdout))
}
