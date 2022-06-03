package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Misc
var connectionState bool

// Windows
var application *widgets.QApplication
var window *widgets.QMainWindow
var mainWindow *widgets.QWidget
var peersWindow *widgets.QWidget
var systray *QSystemTrayIconWithCustomSlot
var systrayMenu *widgets.QMenu
var gridLayoutMainWindow *widgets.QGridLayout
var gridLayoutPeerWindow *widgets.QGridLayout
var groupBoxMainWindow *widgets.QGroupBox
var groupBoxPeerWindow *widgets.QGroupBox

// Labels
var connectionTitleLabel *widgets.QLabel
var subnetTitleLabel *widgets.QLabel
var ipTitleLabel *widgets.QLabel
var peerTitleLabel *widgets.QLabel
var statusLabel *widgets.QLabel

var connectionDataLabel *widgets.QLabel
var ipDataLabel *widgets.QLabel
var subnetDataLabel *widgets.QLabel
var peersCountDataLabel *widgets.QLabel

// Buttons
var connectButton *widgets.QPushButton
var ipButton *widgets.QPushButton
var subnetButton *widgets.QPushButton
var showPeersButton *widgets.QPushButton

var resetPeerButton *widgets.QPushButton

// Lists
var peersList *widgets.QListWidget

type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon

	_ func(f func()) `slot:"triggerSlot,auto"` // Create a slot that takes a function and automatically connect it
}

func (tray *QSystemTrayIconWithCustomSlot) triggerSlot(f func()) { f() } //the slot just needs to call the passed function to execute it inside the main thread

func SetDefaultTexts() {
	connectionTitleLabel = widgets.NewQLabel2("Connection: ", nil, 0)
	subnetTitleLabel = widgets.NewQLabel2("Subnet: ", nil, 0)
	ipTitleLabel = widgets.NewQLabel2("IPv6: ", nil, 0)
	peerTitleLabel = widgets.NewQLabel2("Peers: ", nil, 0)

	connectButton = widgets.NewQPushButton2("Connect", nil)
	ipButton = widgets.NewQPushButton2("Copy Ipv6", nil)
	subnetButton = widgets.NewQPushButton2("Copy Subnet", nil)
	showPeersButton = widgets.NewQPushButton2("Show Peers", nil)

	resetPeerButton = widgets.NewQPushButton2("Search new peers", nil)

	statusLabel = widgets.NewQLabel2("Connection status: ", nil, 0)
	connectionDataLabel = widgets.NewQLabel2("Disconnected", nil, 0)
	ipDataLabel = widgets.NewQLabel2("N/A", nil, 0)
	subnetDataLabel = widgets.NewQLabel2("N/A", nil, 0)
	peersCountDataLabel = widgets.NewQLabel2("No peers found", nil, 0)
}

func getCurrentConnectionInfo() ConnectionInfo {
	resp, err := http.Get("http://localhost:62853/connection")

	if err != nil {
		fmt.Println("Err on connection info")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var connectionInfo = &ConnectionInfo{}
	err = json.Unmarshal([]byte(body), connectionInfo)
	return *connectionInfo

}
func raiseWindow() {
	window.Show()
	window.ActivateWindow()
	window.Raise()
}

func DisconnectToYggdrasil() {
	DisconnectToServer()
	DisconnectUserInterface()
}

func ConnectToYggdrasil() {
	ConnectToServer()
}

func DisconnectUserInterface() {
	// Remove state
	connectionState = false

	// Connection button
	connectButton.SetText("Connect")

	// Connection label
	connectionDataLabel.SetText("Disconnected")
	connectionDataLabel.SetStyleSheet("QLabel {color: red;}")

	// Remove IP & Subnet
	ipDataLabel.SetText("N/A")
	subnetDataLabel.SetText("N/A")

	// Clear peers
	peersList.Clear()
	peersCountDataLabel.SetText("No peers found")
}

func ConnectUserInterface(info ConnectionInfo) {
	// Add state
	connectionState = true

	// Connection label
	connectionDataLabel.SetText("Connected")
	connectionDataLabel.SetStyleSheet("QLabel {color: green;}")

	// Connection button
	connectButton.SetText("Disconnect")

	// IP & Subnet
	ipDataLabel.SetText(info.IpAddress)
	subnetDataLabel.SetText(info.SubnetAddress)

	// Update window
	connectButton.SetEnabled(true)
	connectButton.Repaint()

	// Updating peer text
	peersCountDataLabel.SetText("Searching peers...")

	// Fill peers
	RetrievePeers()
}

func UpdateWindowPeers(details ConnectionDetails) {
	resetPeerButton.BlockSignals(false)

	if len(details.ConnectionPeers) == 0 {
		peersCountDataLabel.SetText("No peers found")
		return
	}

	peersCountDataLabel.SetText("Count: " + strconv.Itoa(len(details.ConnectionPeers)))

	ShowPeersInUserInterface()
}

func SetWindowConnectedError(info ConnectionInfo) {
	widgets.QMessageBox_Critical(nil, "Yggdrasil already running", " You already have an Yggdrasil client running. Can't connect.", widgets.QMessageBox__Ok, 0)
}

func ShowPeersInUserInterface() {
	peersList.Clear()
	info := fetchConnectionData()

	if len(ipAddresses) <= 0 {
		ipAddresses = <-getPeerStats()
	}

	customizedPeers := []PeerSorting{}
	for _, v := range info.ConnectionPeers {
		isThreefoldNode := IsThreefoldNode(ipAddresses, v)
		customPeer := PeerSorting{
			Peer:            v,
			isThreefoldNode: isThreefoldNode,
		}

		customizedPeers = append(customizedPeers, customPeer)
	}

	sortedPeers := SortBy("isThreefoldNode", customizedPeers)

	for i, v := range sortedPeers {
		var item = widgets.NewQListWidgetItem2(v.Peer, peersList, i)

		if v.isThreefoldNode {
			item.SetIcon(gui.NewQIcon5(":/qml/icon.ico"))
		}

		// item.SetIcon(gui.NewQIcon5(":/qml/icon.ico"))

		peersList.AddItem2(item)
	}
}

func AddLayout() {
	// Fill layout components
	SetDefaultTexts()

	// Create layout components
	CreateLayout()

	// Listen to systray triggers
	ListenToSysTrayTriggers()

	// Listen to copy buttons
	ListenToCopyButtons()

	// Listen to action buttons
	ListenToButtonClicks()

}

func userInterface(args yggArgs, ctx context.Context, done chan struct{}) {
	// Initialize windows
	StartUserInterface()

	// Add layout + helper functions
	AddLayout()

	// Default sleep untill API started
	time.Sleep(2 * time.Second)
	checkToStartNetworkServer()

	// Get API Status
	CheckCurrentConnection()

	window.Show()
	application.Exec()
}

func CheckCurrentConnection() {
	connInfo := getCurrentConnectionInfo()

	if connInfo.IpAddress == "" {
		connectionState = false
		return
	}

	connectButton.SetText("Disconnect")
	ipDataLabel.SetText(connInfo.IpAddress)
	subnetDataLabel.SetText(connInfo.SubnetAddress)

	connectionState = true
	connectionDataLabel.SetText("Connected")
	connectionDataLabel.SetStyleSheet("QLabel {color: green;}")

	ipDataLabel.SetText(connInfo.IpAddress)
	subnetDataLabel.SetText(connInfo.SubnetAddress)

	if len(connInfo.ConnectionPeers) != 0 {
		peersCountDataLabel.SetText("Count: " + strconv.Itoa(len(connInfo.ConnectionPeers)))
	}
}

func StartUserInterface() {
	application = widgets.NewQApplication(len(os.Args), os.Args)
	application.SetWindowIcon(gui.NewQIcon5(":/qml/icon.ico"))

	mainWindow = widgets.NewQWidget(nil, 0)
	peersWindow = widgets.NewQWidget(nil, 0)

	peersList = widgets.NewQListWidget(nil)
	peersList.SetFixedHeight(300)
	peersList.SetFixedWidth(450)
	peersList.SetStyleSheet("QListWidget {padding: 10px;} QListWidget::item { margin: 10px; }")

	window = widgets.NewQMainWindow(nil, 0)

	window.SetMinimumSize2(600, 140)

	window.SetFixedSize(core.NewQSize2(600, 200))
	window.SetWindowTitle("ThreeFold Planetary Network")

	window.SetCentralWidget(mainWindow)

	systray = NewQSystemTrayIconWithCustomSlot(nil)
	systray.SetIcon(gui.NewQIcon5(":/qml/icon.ico"))

	systrayMenu = widgets.NewQMenu(nil)

	systray.SetContextMenu(systrayMenu)
	systray.Show()
}

func CreateLayout() {
	gridLayoutMainWindow = widgets.NewQGridLayout2()
	gridLayoutPeerWindow = widgets.NewQGridLayout2()

	groupBoxMainWindow = widgets.NewQGroupBox2("Status", nil)
	groupBoxPeerWindow = widgets.NewQGroupBox2("Peers", nil)

	gridLayoutMainWindow.AddWidget2(statusLabel, 0, 0, core.Qt__AlignLeft)
	gridLayoutMainWindow.AddWidget2(connectionDataLabel, 0, 1, core.Qt__AlignCenter)
	gridLayoutMainWindow.AddWidget2(connectButton, 0, 2, core.Qt__AlignRight)

	gridLayoutMainWindow.AddWidget2(ipTitleLabel, 2, 0, core.Qt__AlignLeft)
	gridLayoutMainWindow.AddWidget2(ipDataLabel, 2, 1, core.Qt__AlignCenter)
	gridLayoutMainWindow.AddWidget2(ipButton, 2, 2, core.Qt__AlignRight)

	gridLayoutMainWindow.AddWidget2(subnetTitleLabel, 3, 0, core.Qt__AlignLeft)
	gridLayoutMainWindow.AddWidget2(subnetDataLabel, 3, 1, core.Qt__AlignCenter)
	gridLayoutMainWindow.AddWidget2(subnetButton, 3, 2, core.Qt__AlignRight)

	gridLayoutMainWindow.AddWidget2(peerTitleLabel, 4, 0, core.Qt__AlignLeft)
	gridLayoutMainWindow.AddWidget2(peersCountDataLabel, 4, 1, core.Qt__AlignCenter)
	gridLayoutMainWindow.AddWidget2(showPeersButton, 4, 2, core.Qt__AlignRight)

	gridLayoutPeerWindow.AddWidget2(peersList, 0, 0, core.Qt__AlignLeft)
	gridLayoutPeerWindow.AddWidget2(resetPeerButton, 1, 0, core.Qt__AlignCenter)

	groupBoxMainWindow.SetLayout(gridLayoutMainWindow)
	groupBoxPeerWindow.SetLayout(gridLayoutPeerWindow)

	mainWindow.SetLayout(widgets.NewQVBoxLayout())
	peersWindow.SetLayout(widgets.NewQVBoxLayout())

	mainWindow.Layout().AddWidget(groupBoxMainWindow)
	peersWindow.Layout().AddWidget(groupBoxPeerWindow)

	connectionDataLabel.SetStyleSheet("QLabel {color: red;}")
}

func ListenToSysTrayTriggers() {
	settingsMenuAction := systrayMenu.AddAction("Status")
	settingsMenuAction.ConnectTriggered(func(bool) {
		window.Show()
		window.ActivateWindow()
		window.Raise()
	})

	yggdrasilVersionMenuAction := systrayMenu.AddAction("Reset")
	yggdrasilVersionMenuAction.ConnectTriggered(func(bool) {
		http.Post("http://localhost:62853/reset", "application/json", bytes.NewBuffer(nil))
		widgets.QMessageBox_Information(nil, "ThreeFold Planetary Network", "All the settings have been reset.\n The application will close itself. \n\n You can simply open it again.", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		os.Exit(0)
	})

	quitMenuAction := systrayMenu.AddAction("Quit")
	quitMenuAction.ConnectTriggered(func(bool) {
		println("Exiting application ... ")
		http.Post("http://localhost:62853/exit", "application/json", bytes.NewBuffer(nil))
		application.Exit(0)
		os.Exit(0)
	})

	window.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		window.Hide()
		event.Ignore()
	})
}

func ListenToCopyButtons() {
	ipButton.ConnectClicked(func(bool) {
		application.Clipboard().SetText(ipDataLabel.Text(), gui.QClipboard__Clipboard)
	})

	subnetButton.ConnectClicked(func(bool) {
		application.Clipboard().SetText(subnetDataLabel.Text(), gui.QClipboard__Clipboard)
	})
}

func ListenToButtonClicks() {
	resetPeerButton.ConnectClicked(func(bool) {
		resetPeerButton.BlockSignals(true)
		resetPeerButton.SetEnabled(false)
		resetPeerButton.Repaint()

		DeleteConfigOnFileSystem()

		if connectionState == true {
			DisconnectToYggdrasil()
		}

		ConnectToYggdrasil()

		resetPeerButton.SetEnabled(true)
		resetPeerButton.Repaint()
	})

	showPeersButton.ConnectClicked(func(bool) {
		showPeersButton.SetEnabled(false)
		showPeersButton.Repaint()

		peersList.Clear()
		ShowPeersInUserInterface()

		showPeersButton.SetEnabled(true)
		showPeersButton.Repaint()

		peersWindow.SetWindowTitle("ThreeFold Planetary Network Peers 222")
		peersWindow.SetFixedSize(core.NewQSize2(500, 400))
		peersWindow.Show()

	})

	connectButton.ConnectClicked(func(bool) {
		connectButton.SetEnabled(false)
		connectButton.Repaint()

		if connectionState == false {
			ConnectToYggdrasil()
			return
		}

		DisconnectToYggdrasil()

		connectButton.SetEnabled(true)
		connectButton.Repaint()
	})
}
