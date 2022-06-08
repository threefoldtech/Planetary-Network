package main

import (
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gologme/log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func startNetworkServer() bool {
	log.Infoln("SERVER NOT RUNNING - STARTING")
	log.Infoln("ASKING USER FOR PASSWORD")

	var password = ""
	var widget = widgets.NewQWidget(nil, 0)
	var dialog = widgets.NewQInputDialog(widget, core.Qt__Dialog)
	dialog.SetWindowTitle("ThreeFold Planetary Network")
	dialog.SetLabelText("ThreeFold Planetary Network would like to automatically\nset up your connection.\n\nTo do this, please provide the password for \"" + getUsername() + "\"")
	dialog.SetTextEchoMode(widgets.QLineEdit__Password)
	dialog.SetInputMethodHints(core.Qt__ImhNone)

	dialog.ConnectAccepted(func() {
		log.Infoln("ACCEPTED PASSWORD")
		password = dialog.TextValue()
		dialog.Close()
	})

	dialog.ConnectRejected(func() {
		log.Errorln("REJECTED PASSWORD")
		os.Exit(1)
	})

	dialog.Exec()

	log.Infoln("STARTING SERVER AS ROOT")
	startNetworkServerAsRoot(password)

	time.Sleep(2 * time.Second)
	_, err2 := http.Get("http://localhost:62853/health")
	if err2 != nil {
		startNetworkServer()
	}
	return false
}

func startNetworkServerAsRoot(password string) {
	ex, errp := os.Executable()
	if errp != nil {
		panic(errp)
	}

	cmd := "echo " + password + " | sudo -S \"" + ex + "\" -server"

	rcmd := exec.Command("bash", "-c", cmd)
	err := rcmd.Start()

	if err != nil {
		log.Errorln("ERROR IN START NETWORK SERVER: ", err.Error())
	}
}
