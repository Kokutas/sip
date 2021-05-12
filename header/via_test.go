package header

import (
	"fmt"
	"log"
	"net"
	"sip"
	"sip/util"
	"testing"
)

func TestVia_Raw(t *testing.T) {
	branch := util.GenerateUnixNanoBranch()
	via := NewVia(
		sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion),
		sip.UDP,
		sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 108), nil).(*sip.Host), 5060).(*sip.HostPort),
		1,
		sip.NewParameters(sip.UDP, "", sip.REGISTER, 0, "192.168.0.101", false, nil).(*sip.Parameters),
		branch,
		"192.168.0.26",
	)
	fmt.Print(via.Raw())
}

func TestVia_JsonString(t *testing.T) {
	branch := util.GenerateUnixNanoBranch()
	via := NewVia(
		sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion),
		sip.UDP,
		sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 108), nil).(*sip.Host), 5060).(*sip.HostPort),
		1,
		nil,
		branch,
		"",
	)
	fmt.Println(via.JsonString())
}

func TestVia_Parser(t *testing.T) {
	// raw := "Via: SIP/2.0/UDP 192.168.0.108:5060;rport;transport=udp;method=register;maddr=192.168.0.101;branch=z9hG4bK5f11bd69d0a510cf5d1c462312cfaff7;received=192.168.0.26\r\n"
	// raw := "Via: SIP/2.0/UDP 192.168.0.108:5060;branch=z9hG4bK5f11bd69d0a510cf5d1c462312cfaff7\r\n"
	raw := "Via: SIP/2.0/UDP www.biadu.com:5060;branch=z9hG4bK5f11bd69d0a510cf5d1c462312cfaff7\r\n"
	via := CreateVia()
	if err := via.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(raw)
	fmt.Println("-----------------------------------------------")
	fmt.Print(via.Raw())
	fmt.Println(via.JsonString())
}

func TestVia_Validator(t *testing.T) {
	branch := util.GenerateUnixNanoBranch()
	via := NewVia(
		sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion),
		sip.UDP,
		sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 108), nil).(*sip.Host), 5060).(*sip.HostPort),
		1,
		nil,
		branch,
		"",
	)
	fmt.Println(via.Validator())
}
