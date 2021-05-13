package sip

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
)



func TestNewSipUri(t *testing.T) {
	su := NewSipUri("sip",
		NewUserInfo("34020000001320000001", "", ""),
		NewHostPort(NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		nil, nil)
	data, err := json.Marshal(su)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestSipUri_Headers(t *testing.T) {

}

func TestSipUri_Parser(t *testing.T) {
	// raw := "sip:34020000001320000001@192.168.0.1:5060"
	raw := "sip:34020000001320000001:Ali12345@192.168.0.1:5060;transport=udp;user=34020000001320000001;method=register;ttl=5;maddr=www.baidu.com;lr;hello=world?haha=1&heihei=4.5"
	su := new(SipUri)
	if err := su.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Println(su.String())
	fmt.Println(su.Raw())
	fmt.Println(raw)
}

func TestSipUri_Raw(t *testing.T) {
	su := NewSipUri("sip",
		NewUserInfo("34020000001320000001", "", "Ali12345"),
		NewHostPort(
			NewHost("www.163.c0m", net.IPv4(192, 168, 0, 1), nil), 5060),
		NewParameters("udp", "34020000001320000001", "register", 5, "www.baidu.com", true, map[string]interface{}{"hello": "world"}),
		map[string]interface{}{"haha": 1, "heihei": 4.5})
	fmt.Println(su.Raw())
}

func TestSipUri_Schema(t *testing.T) {

}

func TestSipUri_SetHeaders(t *testing.T) {

}

func TestSipUri_SetSchema(t *testing.T) {

}

func TestSipUri_String(t *testing.T) {
	su := NewSipUri("sip",
		NewUserInfo("34020000001320000001", "", ""),
		NewHostPort(NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		nil, nil)
	fmt.Println(su.String())
}

func TestSipUri_Validator(t *testing.T) {
	su := NewSipUri("sip",
		NewUserInfo("34020000001320000001", "", ""),
		NewHostPort(NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		nil, nil)
	if err := su.Validator(); err != nil {
		log.Fatal(err)
	}
}
