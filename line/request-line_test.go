package line

import (
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
			sip.NewUserInfo("34020000001320000001", "", ""),
			sip.NewHostPort(
				sip.NewHost("", net.IPv4(192, 168, 0, 26), nil),
				5060),
			sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060}),
			nil), nil)
	version := sip.NewSipVersion(sip.SIP, 2.0)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	fmt.Printf("%s\r\n",rl)
}

func TestRequestLine_Method(t *testing.T) {

}

func TestRequestLine_Parser(t *testing.T) {
	raw := "REGISTER sip:34020000001320000001@192.168.0.26:5060 SIP/2.0\r\n"
	rl := new(RequestLine)
	if err := rl.Parser(raw); err != nil {
		log.Fatal(err)
	}
	//fmt.Print(rl.Raw())
	//fmt.Print(raw)
	fmt.Println(rl.String())
}

func TestRequestLine_Raw(t *testing.T) {
	uri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", ""),
			sip.NewHostPort(
				sip.NewHost("", net.IPv4(192, 168, 0, 26), nil),
				5060),
			nil,
			nil), nil)
	version := sip.NewSipVersion(sip.SIP, 2.0)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	fmt.Print(rl.Raw())
}

func TestRequestLine_SetMethod(t *testing.T) {

}

func TestRequestLine_String(t *testing.T) {
	uri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", ""),
			sip.NewHostPort(
				sip.NewHost("3402000000", net.IPv4(192, 168, 0, 26), nil),
				5060),
			nil,
			nil), nil)
	version := sip.NewSipVersion(sip.SIP, 2.0)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	fmt.Println(rl.String())
}

func TestRequestLine_Validator(t *testing.T) {
	uri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", ""),
			sip.NewHostPort(
				sip.NewHost("3402000000", net.IPv4(192, 168, 0, 26), nil),
				5060),
			nil,
			nil), nil)
	version := sip.NewSipVersion(sip.SIP, 2.0)
	rl := NewRequestLine(sip.REGISTER, uri, version)
	if err := rl.Validator(); err != nil {
		log.Fatal(err)
	}
	//rl:=CreateRequestLine()
	//if err:=rl.Validator(); err != nil {
	//	log.Fatal(err)
	//}
}
