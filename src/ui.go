package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var ipLabel *widgets.QLabel
var subnetLabel *widgets.QLabel
var debugLabel *widgets.QLabel
var connectionLabel *widgets.QLabel
var connectButton *widgets.QPushButton

type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon

	_ func(f func()) `slot:"triggerSlot,auto"` //create a slot that takes a function and automatically connect it
}

func (tray *QSystemTrayIconWithCustomSlot) triggerSlot(f func()) { f() } //the slot just needs to call the passed function to execute it inside the main thread

func uiConnect() ConnectionInfo {
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 10
	resp, err := http.Post("http://localhost:7070/connect", "application/json", bytes.NewBuffer(nil))

	if err != nil {
		fmt.Println("Err on connect")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var connectionInfo = &ConnectionInfo{}
	err = json.Unmarshal([]byte(body), connectionInfo)
	fmt.Println("Connected, ip", connectionInfo.IpAddress)
	return *connectionInfo
}

func uiDisconnect() {
	http.Post("http://localhost:7070/disconnect", "application/json", bytes.NewBuffer(nil))
}

func userInterface(args yggArgs, ctx context.Context, done chan struct{}) {

	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetWindowIcon(gui.NewQIcon5(":/qml/icon.ico"))
	time.Sleep(2 * time.Second)
	checkRoot()

	window := widgets.NewQMainWindow(nil, 0)

	window.SetMinimumSize2(600, 140)
	window.SetWindowTitle("ThreeFold network connector")

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	systray := NewQSystemTrayIconWithCustomSlot(nil)
	systray.SetIcon(gui.NewQIcon5(":/qml/icon.ico"))

	systrayMenu := widgets.NewQMenu(nil)

	settingsMenuAction := systrayMenu.AddAction("Settings")
	settingsMenuAction.ConnectTriggered(func(bool) {
		println("Showing window ...")
		window.Show()
	})

	yggdrasilVersionMenuAction := systrayMenu.AddAction("Reset")
	yggdrasilVersionMenuAction.ConnectTriggered(func(bool) {
		http.Post("http://localhost:7070/reset", "application/json", bytes.NewBuffer(nil))
		widgets.QMessageBox_Information(nil, "ThreeFold network connector", "All the settings have been reset.\n The application will close itself. \n\n You can simply open it again.", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		os.Exit(0)
	})

	quitMenuAction := systrayMenu.AddAction("Quit")
	quitMenuAction.ConnectTriggered(func(bool) {
		println("Exiting application ... ")
		http.Post("http://localhost:7070/exit", "application/json", bytes.NewBuffer(nil))
		app.Exit(0)
		os.Exit(0)
	})

	systray.SetContextMenu(systrayMenu)
	systray.Show()

	connectionState := false
	groupBox := widgets.NewQGroupBox2("Settings", nil)

	// println(window.Type())
	gridLayout := widgets.NewQGridLayout2()

	statusLabel := widgets.NewQLabel2("Connection status: ", nil, 0)
	connectionLabel = widgets.NewQLabel2("Disconnected", nil, 0)
	connectionLabel.SetStyleSheet("QLabel {color: red;}")

	connectButton = widgets.NewQPushButton2("Connect", nil)

	CopyIPButton := widgets.NewQPushButton2("Copy Ipv6", nil)
	copySubnetButton := widgets.NewQPushButton2("Copy Subnet", nil)

	CopyIPButton.ConnectClicked(func(bool) {
		app.Clipboard().SetText(ipLabel.Text(), gui.QClipboard__Clipboard)
	})

	copySubnetButton.ConnectClicked(func(bool) {

		app.Clipboard().SetText(subnetLabel.Text(), gui.QClipboard__Clipboard)
	})

	connectButton.ConnectClicked(func(bool) {
		if !connectionState {
			connectButton.SetText("Disconnect")
			ipLabel.SetText("...")
			subnetLabel.SetText("...")

			connInfo := uiConnect()
			ipLabel.SetText(connInfo.IpAddress)
			subnetLabel.SetText(connInfo.SubnetAddress)

			connectionState = true
			connectionLabel.SetText("Connected")
			connectionLabel.SetStyleSheet("QLabel {color: green;}")

			return
		}
		uiDisconnect()
		connectButton.SetText("Connect")

		connectionLabel.SetText("Disconnected")
		connectionLabel.SetStyleSheet("QLabel {color: red;}")

		connectionState = false
		// defer n.shutdown()
		// go submain()
		// widgets.QMessageBox_Information(nil, "OK", "Connecting ...", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)

		ipLabel.SetText("N/A")
		subnetLabel.SetText("N/A")

		// Catch interrupts from the operating system to exit gracefully.
		c := make(chan os.Signal, 1)
		r := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		signal.Notify(r, os.Interrupt, syscall.SIGHUP)
		// Capture the service being stopped on Windows.
		// minwinsvc.SetOnExit(n.shutdown)
		// defer n.shutdown()
		// Wait for the terminate/interrupt signal. Once a signal is received, the
		// deferred Stop function above will run which will shut down TUN/TAP.
	})

	gridLayout.AddWidget2(statusLabel, 0, 0, core.Qt__AlignLeft)
	gridLayout.AddWidget2(connectionLabel, 0, 1, core.Qt__AlignCenter)
	gridLayout.AddWidget2(connectButton, 0, 2, core.Qt__AlignRight)

	ipLabelInfo := widgets.NewQLabel2("Ipv6: ", nil, 0)
	subnetLabelInfo := widgets.NewQLabel2("Subnet: ", nil, 0)
	//	debugLabelInfo := widgets.NewQLabel2("Debug: ", nil, 0)

	ipLabel = widgets.NewQLabel2("N/A", nil, 0)
	subnetLabel = widgets.NewQLabel2("N/A", nil, 0)
	//	debugLabel = widgets.NewQLabel2("Debug info", nil, 0)

	gridLayout.AddWidget2(ipLabelInfo, 2, 0, core.Qt__AlignLeft)
	gridLayout.AddWidget2(ipLabel, 2, 1, core.Qt__AlignCenter)
	gridLayout.AddWidget2(CopyIPButton, 2, 2, core.Qt__AlignRight)

	gridLayout.AddWidget2(subnetLabelInfo, 3, 0, core.Qt__AlignLeft)
	gridLayout.AddWidget2(subnetLabel, 3, 1, core.Qt__AlignCenter)
	gridLayout.AddWidget2(copySubnetButton, 3, 2, core.Qt__AlignRight)

	// Debugging purposes

	// gridLayout.AddWidget2(debugLabelInfo, 3, 0, core.Qt__AlignCenter)
	// gridLayout.AddWidget2(debugLabel, 3, 1, core.Qt__AlignCenter)

	groupBox.SetLayout(gridLayout)
	widget.Layout().AddWidget(groupBox)

	window.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		widgets.QMessageBox_Information(nil, "ThreeFold network connector", "The ThreeFold network connector will be minimized.", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		window.Hide()
		event.Ignore()
	})

	window.Show()
	app.Exec()
}
