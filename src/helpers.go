package main

import (
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/gologme/log"
)

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

func SortBy(jsonField string, arr []PeerSorting) []PeerSorting {
	if len(arr) < 1 {
		return []PeerSorting{}
	}

	// first we find the field based on the json tag
	valueType := reflect.TypeOf(arr[0])

	var field reflect.StructField

	for i := 0; i < valueType.NumField(); i++ {
		field = valueType.Field(i)

		if field.Tag.Get("json") == jsonField {
			break
		}
	}

	// then we sort based on the type of the field
	sort.Slice(arr, func(i, j int) bool {
		v1 := reflect.ValueOf(arr[i]).FieldByName(field.Name)
		v2 := reflect.ValueOf(arr[j]).FieldByName(field.Name)

		switch field.Type.Name() {
		case "int":
			return int(v1.Int()) < int(v2.Int())
		case "string":
			return v1.String() < v2.String()
		case "bool":
			return v1.Bool() // return small numbers first
		default:
			return false // return unmodified
		}
	})

	return arr
}

func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Errorln(err)
	}

	return dir
}
