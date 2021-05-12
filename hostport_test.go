package sip

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
)

func TestNewHostPort(t *testing.T) {
	h := NewHost("", net.IPv4(192, 168, 0, 26), nil).(*Host)
	hp := NewHostPort(h, 5060)
	data, err := json.Marshal(hp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func Test_HostPort_Raw(t *testing.T) {
	h := NewHost("", net.IPv4(192, 168, 0, 26), nil).(*Host)
	hp := NewHostPort(h, 5060)
	fmt.Println(hp.Raw())
	data, err := json.Marshal(hp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func Test_HostPort_JsonString(t *testing.T) {
	h := NewHost("", net.IPv4(192, 168, 0, 26), nil).(*Host)
	hp := NewHostPort(h, 5060)
	if res := hp.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func Test_HostPort_Parser(t *testing.T) {
	raw := "192.168.0.1"
	hp := CreateHostPort()
	if err := hp.Parser(raw); err != nil {
		log.Fatal(err)
	}
	if res := hp.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func Test_HostPort_Validator(t *testing.T) {
	hp := NewHostPort(nil, 5060)
	if err := hp.Validator(); err != nil {
		log.Fatal(err)
	}
}
