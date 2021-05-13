package header

import (
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestFrom_Addr(t *testing.T) {

}

func TestFrom_DisplayName(t *testing.T) {

}

func TestFrom_Field(t *testing.T) {

}

func TestFrom_Parser(t *testing.T) {
	raw := "From: \"Alisa\" <sip:34020000001320000001:123@192.168.0.1:5060;rport=5060;transport=udp;received=192.168.0.1?hello=world>;tag=123456\r\n"
	//raw := "From: \"Alisa\" <sip:34020000001320000001:123@192.168.0.1:5060;rport=5060;transport=udp;received=192.168.0.1?hello=world>\r\n"
	fmt.Print(raw)
	from := new(From)
	if err := from.Parser(raw); err != nil {
		log.Fatalln(err)
	}
	fmt.Print(from.Raw())
	fmt.Println(from.String())
}

func TestFrom_Raw(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123"),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		sip.NewParameters(sip.SIP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	from := NewFrom("Alisa", addr, tag)
	fmt.Println(from.Raw())
}

func TestFrom_SetAddr(t *testing.T) {

}

func TestFrom_SetDisplayName(t *testing.T) {

}

func TestFrom_SetField(t *testing.T) {

}

func TestFrom_SetTag(t *testing.T) {

}

func TestFrom_String(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123"),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		nil, nil)
	tag := "123456"
	from := NewFrom("Alisa", addr, tag)
	if res := from.String(); res != "" {
		fmt.Println(res)
	}
}

func TestFrom_Tag(t *testing.T) {

}

func TestFrom_Validator(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123"),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	from := NewFrom("Alisa", addr, tag)
	//from = new(From)
	//from.SetField("hello")
	if err := from.Validator(); err != nil {
		log.Fatalln(err)
	}
}

func TestNewFrom(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", ""),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		nil, nil)
	tag := "123456"
	from := NewFrom("", addr, tag)
	fmt.Printf("%s\r\n", from)
}
