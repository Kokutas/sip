package header

import (
	"fmt"
	"log"
	"net"
	"github.com/kokutas/sip"
	"testing"
)

func TestContact_Addr(t *testing.T) {

}

func TestContact_CpExpires(t *testing.T) {

}

func TestContact_Cpq(t *testing.T) {

}

func TestContact_DisplayName(t *testing.T) {

}

func TestContact_Extension(t *testing.T) {

}

func TestContact_Field(t *testing.T) {

}

func TestContact_Parser(t *testing.T) {
	raw := "Contact: \"34020000001310000001\" <sip:34020000001320000001@192.168.0.26:5060>;q=0.7;expires=3600\r\n"
	contact := new(Contact)
	if err := contact.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(contact.Raw())
	fmt.Println(contact.String())
}

func TestContact_Raw(t *testing.T) {
	contact := NewContact("34020000001310000001",
		sip.NewSipUri(sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", ""),
			sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 26), nil), 5060),
			nil, nil), 0.7, 3600, nil)
	fmt.Print(contact.Raw())
}

func TestContact_SetAddr(t *testing.T) {

}

func TestContact_SetCpExpires(t *testing.T) {

}

func TestContact_SetCpq(t *testing.T) {

}

func TestContact_SetDisplayName(t *testing.T) {

}

func TestContact_SetExtension(t *testing.T) {

}

func TestContact_SetField(t *testing.T) {

}

func TestContact_String(t *testing.T) {
	contact := NewContact("34020000001310000001",
		sip.NewSipUri(sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", ""),
			sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 26), nil), 5060),
			nil, nil), 0.7, 3600, nil)
	fmt.Println(contact.String())
}

func TestContact_Validator(t *testing.T) {
	contact := NewContact("34020000001310000001",
		sip.NewSipUri(sip.SIP,
			sip.NewUserInfo("34020000001320000001", "", ""),
			sip.NewHostPort(sip.NewHost("", net.IPv4(192, 168, 0, 26), nil), 5060),
			nil, nil), 0.7, 3600, nil)
	fmt.Println(contact.Validator())
}

func TestNewContact(t *testing.T) {

}
