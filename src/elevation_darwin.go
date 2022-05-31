package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func startNetworkServer() bool {

	fmt.Println("Server not running, starting it up")

	fmt.Println("Asking user for password")

	var password = ""
	var widget = widgets.NewQWidget(nil, 0)
	var dialog = widgets.NewQInputDialog(widget, core.Qt__Dialog)
	dialog.SetWindowTitle("ThreeFold Planetary Network")
	dialog.SetLabelText("ThreeFold Planetary Network would like to automatically\nset up your connection to the ThreeFold Network.\n\nTo do this, please provide the password for \"" + getUsername() + "\"")
	dialog.SetTextEchoMode(widgets.QLineEdit__Password)
	dialog.SetInputMethodHints(core.Qt__ImhNone)

	dialog.ConnectAccepted(func() {
		fmt.Println("Accepted")
		password = dialog.TextValue()
		dialog.Close()
	})

	dialog.ConnectRejected(func() {
		fmt.Println("Rejected")
		os.Exit(1)
	})

	dialog.Exec()

	fmt.Println("Starting server as root")
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
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	cmd := "echo " + password + " | sudo -S \"" + ex + "\" -server"

	rcmd := exec.Command("bash", "-c", cmd)
	err := rcmd.Start()

	if err != nil {
		fmt.Println(err)
		// os.Exit(1)
	}
}
