package header

import (
	"fmt"
	"log"
	"testing"

	"github.com/kokutas/sip"
)

func TestAuthorization_Algorithm(t *testing.T) {

}

func TestAuthorization_AuthParam(t *testing.T) {

}

func TestAuthorization_AuthSchema(t *testing.T) {

}

func TestAuthorization_Cnonce(t *testing.T) {

}

func TestAuthorization_Dresponse(t *testing.T) {

}

func TestAuthorization_Field(t *testing.T) {

}

func TestAuthorization_Nonce(t *testing.T) {

}

func TestAuthorization_NonceCount(t *testing.T) {

}

func TestAuthorization_Opaque(t *testing.T) {

}

func TestAuthorization_Parser(t *testing.T) {
	raw := "Authorization: Digest username=\"34020000001320000001\",realm=\"3402000000\",nonce=\"nonce456\",uri=\"sip:34020000002000000001@3402000000\",response=\"response123\",nc=\"78787878\",algorithm=MD5"
	au := new(Authorization)
	if err := au.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(au.Raw())
	fmt.Println(au.String())
	fmt.Println(raw)
}

func TestAuthorization_Qop(t *testing.T) {

}

func TestAuthorization_Raw(t *testing.T) {
	au := NewAuthorization(sip.Digest, "34020000001320000001", "3402000000", "nonce456",
		sip.NewSipUri(sip.SIP,
			sip.NewUserInfo("34020000002000000001", "", ""),
			sip.NewHostPort(sip.NewHost("3402000000", nil, nil), 0), nil, nil), "response123", "", "", "", "", "xxxx", nil)
	fmt.Print(au.Raw())
}

func TestAuthorization_Realm(t *testing.T) {

}

func TestAuthorization_SetAlgorithm(t *testing.T) {

}

func TestAuthorization_SetAuthParam(t *testing.T) {

}

func TestAuthorization_SetAuthSchema(t *testing.T) {

}

func TestAuthorization_SetCnonce(t *testing.T) {

}

func TestAuthorization_SetDresponse(t *testing.T) {

}

func TestAuthorization_SetField(t *testing.T) {

}

func TestAuthorization_SetNonce(t *testing.T) {

}

func TestAuthorization_SetNonceCount(t *testing.T) {

}

func TestAuthorization_SetOpaque(t *testing.T) {

}

func TestAuthorization_SetQop(t *testing.T) {

}

func TestAuthorization_SetRealm(t *testing.T) {

}

func TestAuthorization_SetUri(t *testing.T) {

}

func TestAuthorization_SetUserName(t *testing.T) {

}

func TestAuthorization_String(t *testing.T) {
	au := NewAuthorization(sip.Digest, "34020000001320000001", "3402000000", "nonce456", sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000002000000001", "", ""),
		sip.NewHostPort(sip.NewHost("3402000000", nil, nil), 0), nil, nil), "response123", "", "", "", "", "xxxx", nil)
	fmt.Println(au.String())
}

func TestAuthorization_Uri(t *testing.T) {

}

func TestAuthorization_UserName(t *testing.T) {

}

func TestAuthorization_Validator(t *testing.T) {
	au := NewAuthorization(sip.Digest, "34020000001320000001", "3402000000", "nonce456", sip.NewSipUri(sip.SIP,
		sip.NewUserInfo("34020000002000000001", "", ""),
		sip.NewHostPort(sip.NewHost("3402000000", nil, nil), 0), nil, nil), "response123", "", "", "", "", "xxxx", nil)
	fmt.Println(au.Validator())

}

func TestNewAuthorization(t *testing.T) {

}
