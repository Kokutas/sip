package header

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestCreateFrom(t *testing.T) {
	from := CreateFrom()
	fmt.Println(from)
}

func TestFrom_JsonString(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		nil, nil)
	tag := "123456"
	from := NewFrom("Alisa", addr.(*sip.SipUri), tag)
	if res := from.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func TestFrom_Parser(t *testing.T) {
	raw := "From: \"Alisa\" <sip:34020000001320000001:123@192.168.0.1:5060;rport=5060;transport=udp;received=192.168.0.1?hello=world>;tag=123456\r\n"
	//raw := "From: \"Alisa\" <sip:34020000001320000001:123@192.168.0.1:5060;rport=5060;transport=udp;received=192.168.0.1?hello=world>\r\n"
	fmt.Print(raw)
	from := CreateFrom()
	if err := from.Parser(raw); err != nil {
		log.Fatalln(err)
	}
	fmt.Print(from.(*From).Raw())
	fmt.Println(from.JsonString())
}

func TestFrom_Raw(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		sip.NewParameters(sip.SIP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}).(*sip.Parameters),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	from := NewFrom("Alisa", addr.(*sip.SipUri), tag)
	fmt.Println(from.Raw())
}

func TestFrom_Validator(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}).(*sip.Parameters),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	from := NewFrom("Alisa", addr.(*sip.SipUri), tag)
	//from = CreateFrom()
	//from.(*From).Field = "hello"
	if err := from.Validator(); err != nil {
		log.Fatalln(err)
	}
}

func TestNewFrom(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		nil, nil)
	tag := "123456"
	from := NewFrom("", addr.(*sip.SipUri), tag)
	data, err := json.Marshal(from)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\r\n", data)
}
