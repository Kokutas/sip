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
		NewUserInfo("34020000001320000001", "", "").(*UserInfo),
		NewHostPort(NewHost("", net.IPv4(192, 168, 0, 1), nil).(*Host), 5060).(*HostPort),
		nil, nil)
	data, err := json.Marshal(su)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func Test_SipUri_Raw(t *testing.T) {
	su := NewSipUri("sip",
		NewUserInfo("34020000001320000001", "", "Ali12345").(*UserInfo),
		NewHostPort(
			NewHost("www.163.c0m", net.IPv4(192, 168, 0, 1), nil).(*Host), 5060).(*HostPort),
		NewParameters("udp", "34020000001320000001", "register", 5, "www.baidu.com", true, map[string]interface{}{"hello": "world"}).(*Parameters),
		map[string]interface{}{"haha": 1, "heihei": 4.5})
	fmt.Println(su.Raw())
}
func Test_SipUri_JsonString(t *testing.T) {
	su := NewSipUri("sip",
		NewUserInfo("34020000001320000001", "", "").(*UserInfo),
		NewHostPort(NewHost("", net.IPv4(192, 168, 0, 1), nil).(*Host), 5060).(*HostPort),
		nil, nil)
	if res := su.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func Test_SipUri_Parser(t *testing.T) {
	// raw := "sip:34020000001320000001@192.168.0.1:5060"
	raw := "sip:34020000001320000001:Ali12345@192.168.0.1:5060;transport=udp;user=34020000001320000001;method=register;ttl=5;maddr=www.baidu.com;lr;hello=world?haha=1&heihei=4.5"
	su := CreateSipUri()
	if err := su.Parser(raw); err != nil {
		log.Fatal(err)
	}
	if res := su.JsonString(); res != "" {
		fmt.Println(res)
	}
	fmt.Println(su.Raw())
	fmt.Println(raw)
}

func Test_SipUri_Validator(t *testing.T) {
	su := NewSipUri("sip",
		NewUserInfo("34020000001320000001", "", "").(*UserInfo),
		NewHostPort(NewHost("", net.IPv4(192, 168, 0, 1), nil).(*Host), 5060).(*HostPort),
		nil, nil)
	if err := su.Validator(); err != nil {
		log.Fatal(err)
	}
}
