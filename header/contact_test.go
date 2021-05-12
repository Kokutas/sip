package header

import (
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestContact_Raw(t *testing.T) {
	contact := NewContact("34020000001310000001",
		sip.NewSipUri(sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
			sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 26), nil).(*sip.Host), 5060).(*sip.HostPort),
			nil, nil).(*sip.SipUri), 0.7, 3600, nil)
	fmt.Print(contact.Raw())
}

func TestContact_JsonString(t *testing.T) {
	contact := NewContact("34020000001310000001",
		sip.NewSipUri(sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
			sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 26), nil).(*sip.Host), 5060).(*sip.HostPort),
			nil, nil).(*sip.SipUri), 0.7, 3600, nil)
	fmt.Println(contact.JsonString())
}

func TestContact_Parser(t *testing.T) {
	raw := "Contact: \"34020000001310000001\" <sip:34020000001320000001@192.168.0.26:5060>;q=0.7;expires=3600\r\n"
	contact := CreateContact()
	if err := contact.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(contact.Raw())
	fmt.Println(contact.JsonString())
}

func TestContact_Validator(t *testing.T) {
	contact := NewContact("34020000001310000001",
		sip.NewSipUri(sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", "").(*sip.UserInfo),
			sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 26), nil).(*sip.Host), 5060).(*sip.HostPort),
			nil, nil).(*sip.SipUri), 0.7, 3600, nil)
	fmt.Println(contact.Validator())
}
