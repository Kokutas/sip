package sip

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
)

func TestNewHost(t *testing.T) {
	host := NewHost("www.baidu.com", net.IPv4(192, 168, 0, 1), nil)
	data, err := json.Marshal(host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func Test_Host_Raw(t *testing.T) {
	host := NewHost("", nil, net.ParseIP("fe80::9134:a673:dd9d:9656"))
	fmt.Print(host.Raw())
	data, err := json.Marshal(host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func Test_Host_JsonString(t *testing.T) {
	host := NewHost("www.baidu.com", net.IPv4(192, 168, 0, 1), nil)
	fmt.Print(host.JsonString())
}

func Test_Host_Parser(t *testing.T) {
	raw := "www.baidu.com"
	h := CreateHost()
	if err := h.Parser(raw); err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(h)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func Test_Host_Validator(t *testing.T) {
	host := NewHost("", net.IPv4(192, 168, 0, 1), nil)
	if err := host.Validator(); err != nil {
		log.Fatal(err)
	}
}
