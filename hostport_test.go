package sip

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
)


func TestHostPort_Parser(t *testing.T) {
	raw := "192.168.0.1"
	hp := new(HostPort)
	if err := hp.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Println(hp.String())
}

func TestHostPort_Port(t *testing.T) {

}

func TestHostPort_Raw(t *testing.T) {
	h := NewHost("", net.IPv4(192, 168, 0, 26), nil)
	hp := NewHostPort(h, 5060)
	fmt.Println(hp.Raw())
	data, err := json.Marshal(hp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestHostPort_SetPort(t *testing.T) {

}

func TestHostPort_String(t *testing.T) {


}

func TestHostPort_Validator(t *testing.T) {
	hp := NewHostPort(new(Host), 5060)
	if err := hp.Validator(); err != nil {
		log.Fatal(err)
	}
}

func TestNewHostPort(t *testing.T) {
	h := NewHost("", net.IPv4(192, 168, 0, 26), nil)
	hp := NewHostPort(h, 5060)
	data, err := json.Marshal(hp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}
