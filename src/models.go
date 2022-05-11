package main

type YggdrasilIPAddress struct {
	FullIPAddress   string
	IPAddress       string
	latency         float64
	isThreefoldNode bool
}

type ConnectionInfo struct {
	IpAddress       string
	SubnetAddress   string
	PublicKey       string
	ConnectionPeers []string
	Error           string
}

type ConnectionDetails struct {
	IpAddress       string
	ConnectionPeers []string
	Error           string
}
