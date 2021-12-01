package main

import "runtime"

type Config struct {
	yggdrasil_config_location string
	SubnetAddress             string
	PublicKey                 string
	Error                     string
}

var app_config Config

func init_config() {
	if runtime.GOOS == "windows" {
		app_config.yggdrasil_config_location = "c:/threefold_yggdrasil.conf"
		return
	}
	app_config.yggdrasil_config_location = "/etc/threefold_yggdrasil.conf"

}
