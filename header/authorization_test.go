package header

import (
	"fmt"
	"log"
	"sip"
	"testing"
)

func TestAuthorization_Raw(t *testing.T) {
	au := NewAuthorization(sip.Digest, "34020000001320000001", "3402000000", "nonce456", sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000002000000001", "", "").(*sip.UserInfo),
		sip.NewHostPort(sip.NewHost("3402000000", nil, nil).(*sip.Host), 0).(*sip.HostPort), nil, nil).(*sip.SipUri), "response123", "", "", "", "", "xxxx", nil)
	fmt.Print(au.Raw())
}

func TestAuthorization_Parser(t *testing.T) {
	raw := "Authorization: Digest username=\"34020000001320000001\",realm=\"3402000000\",nonce=\"nonce456\",uri=\"sip:34020000002000000001@3402000000\",response=\"response123\",nc=\"78787878\",algorithm=MD5"
	au := CreateAuthorization()
	if err := au.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(au.Raw())
	fmt.Println(au.JsonString())
	fmt.Println(raw)
}

func TestAuthorization_JsonString(t *testing.T) {
	au := NewAuthorization(sip.Digest, "34020000001320000001", "3402000000", "nonce456", sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000002000000001", "", "").(*sip.UserInfo),
		sip.NewHostPort(sip.NewHost("3402000000", nil, nil).(*sip.Host), 0).(*sip.HostPort), nil, nil).(*sip.SipUri), "response123", "", "", "", "", "xxxx", nil)
	fmt.Println(au.JsonString())
}

func TestAuthorization_Validator(t *testing.T) {
	au := NewAuthorization(sip.Digest, "34020000001320000001", "3402000000", "nonce456", sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000002000000001", "", "").(*sip.UserInfo),
		sip.NewHostPort(sip.NewHost("3402000000", nil, nil).(*sip.Host), 0).(*sip.HostPort), nil, nil).(*sip.SipUri), "response123", "", "", "", "", "xxxx", nil)
	fmt.Println(au.Validator())
}
