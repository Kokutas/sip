package sip

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
)

func TestHost_HostName(t *testing.T) {

}

func TestHost_Ipv4Address(t *testing.T) {

}

func TestHost_Ipv6Reference(t *testing.T) {

}

func TestHost_Parser(t *testing.T) {
	raw := "www.baidu.com"
	h := new(Host)
	if err := h.Parser(raw); err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(h)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestHost_Raw(t *testing.T) {
	host := NewHost("", nil, net.ParseIP("fe80::9134:a673:dd9d:9656"))
	fmt.Print(host.Raw())
	data, err := json.Marshal(host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}
func TestHost_SetHostName(t *testing.T) {

}

func TestHost_SetIpv4Address(t *testing.T) {
}

func TestHost_SetIpv6Reference(t *testing.T) {
}

func TestHost_String(t *testing.T) {
	host := NewHost("www.baidu.com", net.IPv4(192, 168, 0, 1), nil)
	fmt.Print(host.String())
}

func TestHost_Validator(t *testing.T) {
	host := NewHost("", net.IPv4(192, 168, 0, 1), nil)
	if err := host.Validator(); err != nil {
		log.Fatal(err)
	}
}

func TestNewHost(t *testing.T) {
	host := NewHost("www.baidu.com", net.IPv4(192, 168, 0, 1), nil)
	data, err := json.Marshal(host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}
