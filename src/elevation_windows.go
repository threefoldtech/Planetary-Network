package main

import (
	"github.com/gologme/log"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func startNetworkServer() {
	startNetworkServerAsAdmin()
}

func startNetworkServerAsAdmin() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := " -server"

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		log.Errorln("ERROR IN START NETWORK SERVER: ", err.Error())
	}
}
