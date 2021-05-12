package line

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestNewRequestLine(t *testing.T) {
	uri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
			sip.NewHostPort(
				sip.NewHost("", net.IPv4(192, 168, 0, 26), nil).(*sip.Host),
				5060).(*sip.HostPort),
			sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060}).(*sip.Parameters),
			nil).(*sip.SipUri), nil).(*sip.RequestUri)
	version := sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	data, err := json.Marshal(rl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
	fmt.Print(rl.Raw())
}

func TestRequestLine_Raw(t *testing.T) {
	uri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
			sip.NewHostPort(
				sip.NewHost("", net.IPv4(192, 168, 0, 26), nil).(*sip.Host),
				5060).(*sip.HostPort),
			nil,
			nil).(*sip.SipUri), nil).(*sip.RequestUri)
	version := sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	fmt.Print(rl.Raw())
}

func TestRequestLine_JsonString(t *testing.T) {
	uri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
			sip.NewHostPort(
				sip.NewHost("3402000000", net.IPv4(192, 168, 0, 26), nil).(*sip.Host),
				5060).(*sip.HostPort),
			nil,
			nil).(*sip.SipUri), nil).(*sip.RequestUri)
	version := sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	if res := rl.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func TestRequestLine_Parser(t *testing.T) {
	raw := "REGISTER sip:34020000001320000001@192.168.0.26:5060 SIP/2.0\r\n"
	rl := CreateRequestLine().(*RequestLine)
	if err := rl.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(rl.Raw())
	fmt.Print(raw)
	if res := rl.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func TestRequestLine_Validator(t *testing.T) {

	uri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
			sip.NewHostPort(
				sip.NewHost("3402000000", net.IPv4(192, 168, 0, 26), nil).(*sip.Host),
				5060).(*sip.HostPort),
			nil,
			nil).(*sip.SipUri), nil).(*sip.RequestUri)
	version := sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	if err := rl.Validator(); err != nil {
		log.Fatal(err)
	}
	//rl:=CreateRequestLine()
	//if err:=rl.Validator(); err != nil {
	//	log.Fatal(err)
	//}
}
