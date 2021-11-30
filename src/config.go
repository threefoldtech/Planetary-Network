package main

// +build windows
var yggdrasil_config_location = "%USERPROFILE%/threefold_yggdrasil.conf"

// +build !windows
var yggdrasil_config_location = "/etc/threefold_yggdrasil.conf"
