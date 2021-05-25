package sip

import (
	"fmt"
	"sync"
	"testing"
)

func TestParameters_Raw(t *testing.T) {
	parameters := NewParameters("udp", "34020000001320000001", "REGISTER", 5, "192.168.0.1", true, sync.Map{})
	result := parameters.Raw()
	fmt.Println(result.String())
}

func TestParameters_Parse(t *testing.T) {
	raws := []string{
		";transport=udp;user=34020000001320000001;method=REGISTER;maddr=192.168.0.1;lr;ttl=5;;",
		";transport=udp;user=34020000001320000001;maddr=192.168.0.1;lr;ttl=5;;method=REGISTER;",
		";transport=udp;user=34020000001320000001;maddr=192.168.0.1;lr;ttl=5;;method=REGISTER;token",
	}
	for index, raw := range raws {
		parameters := new(Parameters)
		parameters.Parse(raw)
		if len(parameters.GetSource()) > 0 {
			fmt.Print("index: ", index, ",transport: ", parameters.GetTransport(), ",user: ", parameters.GetUser(), ",ttl: ", parameters.GetTtl(), ",maddr: ", parameters.GetMaddr(), ",lr: ", parameters.GetLr(), ",method: ", parameters.GetMethod())
			other := parameters.GetOther()
			other.Range(func(key, value interface{}) bool {
				fmt.Print(",key= ", key, ",value =", value)
				return true
			})
			fmt.Println()
			result := parameters.Raw()
			fmt.Println(result.String())
		}

	}
}
