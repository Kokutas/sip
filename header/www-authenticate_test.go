package header

import (
	"fmt"
	"log"
	"sip"
	"testing"
)

func TestNewWWWAuthenticate(t *testing.T) {}

func TestWWWAuthenticate_Algorithm(t *testing.T) {}

func TestWWWAuthenticate_AuthParam(t *testing.T) {}

func TestWWWAuthenticate_AuthSchema(t *testing.T) {}

func TestWWWAuthenticate_Domain(t *testing.T) {}

func TestWWWAuthenticate_Field(t *testing.T) {}

func TestWWWAuthenticate_Nonce(t *testing.T) {}

func TestWWWAuthenticate_Opaque(t *testing.T) {}

func TestWWWAuthenticate_Parser(t *testing.T) {
	raw := "WWW-Authenticate: Digest realm=\"3402000000\",nonce=\"6fe9ba44a76be22a\",algorithm=MD5\r\n"
	wa := new(WWWAuthenticate)
	if err := wa.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(wa.Raw())
	fmt.Println(wa.String())
	fmt.Println(raw)
}

func TestWWWAuthenticate_QopOptions(t *testing.T) {}

func TestWWWAuthenticate_Raw(t *testing.T) {
	wa := NewWWWAuthenticate(sip.Digest, "3402000000", nil, "6fe9ba44a76be22a", "", "", "md5", "", nil)
	fmt.Print(wa.Raw())
}

func TestWWWAuthenticate_Realm(t *testing.T) {}

func TestWWWAuthenticate_SetAlgorithm(t *testing.T) {}

func TestWWWAuthenticate_SetAuthParam(t *testing.T) {}

func TestWWWAuthenticate_SetAuthSchema(t *testing.T) {}

func TestWWWAuthenticate_SetDomain(t *testing.T) {}

func TestWWWAuthenticate_SetField(t *testing.T) {}

func TestWWWAuthenticate_SetNonce(t *testing.T) {}

func TestWWWAuthenticate_SetOpaque(t *testing.T) {}

func TestWWWAuthenticate_SetQopOptions(t *testing.T) {}

func TestWWWAuthenticate_SetRealm(t *testing.T) {}

func TestWWWAuthenticate_SetStale(t *testing.T) {}

func TestWWWAuthenticate_Stale(t *testing.T) {}

func TestWWWAuthenticate_String(t *testing.T) {
	wa := NewWWWAuthenticate(sip.Digest, "3402000000", nil, "6fe9ba44a76be22a", "", "", "md5", "", nil)
	fmt.Println(wa.String())
}

func TestWWWAuthenticate_Validator(t *testing.T) {
	wa := NewWWWAuthenticate(sip.Digest, "3402000000", nil, "6fe9ba44a76be22a", "", "", "md5", "", nil)
	fmt.Println(wa.Validator())
}
