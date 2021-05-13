package header

import (
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestNewTo(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123"),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr, tag)
	fmt.Println(to)

}

func TestTo_Addr(t *testing.T) {
}

func TestTo_DisplayName(t *testing.T) {
}

func TestTo_Field(t *testing.T) {
}

func TestTo_Parser(t *testing.T) {
	raw := "To: \"Alisa\" <sip:34020000001320000001:123@192.168.0.1:5060;rport=5060;transport=udp;received=192.168.0.1?hello=world>;tag=123456\n"
	to := new(To)
	if err := to.Parser(raw); err != nil {
		log.Fatalln(err)
	}
	fmt.Print(raw)
	fmt.Print(to.Raw())
	fmt.Println(to.String())
}

func TestTo_Raw(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123"),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr, tag)
	fmt.Println(to.Raw())
}

func TestTo_SetAddr(t *testing.T) {
}

func TestTo_SetDisplayName(t *testing.T) {
}

func TestTo_SetField(t *testing.T) {
}

func TestTo_SetTag(t *testing.T) {
}

func TestTo_String(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123"),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr, tag)
	fmt.Println(to.String())
}

func TestTo_Tag(t *testing.T) {
}

func TestTo_Validator(t *testing.T) {
	addr := sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000001320000001", "", "123"),
		sip.NewHostPort(
			sip.NewHost("", net.IPv4(192, 168, 0, 1), nil), 5060),
		sip.NewParameters(sip.UDP, "", "", 0, "", false, map[string]interface{}{"rport": 5060, "received": net.IPv4(192, 168, 0, 1)}),
		map[string]interface{}{"hello": "world"})
	tag := "123456"
	to := NewTo("Alisa", addr, tag)
	fmt.Println(to.String())
	fmt.Println(to.Validator())
}
