package header

import (
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestCreateTo(t *testing.T) {
}

func TestNewTo(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}).(*sip.Parameters),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr.(*sip.SipUri), tag)
	fmt.Println(to)
}

func TestTo_JsonString(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}).(*sip.Parameters),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr.(*sip.SipUri), tag)
	fmt.Println(to.JsonString())
}

func TestTo_Parser(t *testing.T) {
	raw := "To: \"Alisa\" <sip:34020000001320000001:123@192.168.0.1:5060;rport=5060;transport=udp;received=192.168.0.1?hello=world>;tag=123456\n"
	to := CreateTo()
	if err := to.Parser(raw); err != nil {
		log.Fatalln(err)
	}
	fmt.Print(raw)
	fmt.Print(to.Raw())
	fmt.Println(to.JsonString())
}

func TestTo_Raw(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}).(*sip.Parameters),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr.(*sip.SipUri), tag)
	fmt.Println(to.Raw())
}

func TestTo_Validator(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123").(*sip.UserInfo),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil).(*sip.Host), 5060).(*sip.HostPort),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}).(*sip.Parameters),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr.(*sip.SipUri), tag)
	fmt.Println(to.JsonString())
	fmt.Println(to.Validator())
}
