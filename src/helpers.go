package main

import "strings"

func IsThreefoldNode(a []YggdrasilIPAddress, x string) bool {
	result := strings.ReplaceAll(x, "tls://", "")
	result = strings.ReplaceAll(result, "tcp://", "")
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	splitResult := strings.Split(result, ":")
	finalResult := strings.ReplaceAll(result, ":"+splitResult[len(splitResult)-1], "")

	for _, n := range a {
		if finalResult == n.RealIP {
			return n.isThreefoldNode
		}
	}
	return false
}
