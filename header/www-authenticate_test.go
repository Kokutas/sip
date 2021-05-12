package header

import (
	"fmt"
	"log"
	"sip"
	"testing"
)

func TestWWWAuthenticate_Raw(t *testing.T) {
	wa := NewWWWAuthenticate(sip.Digest, "3402000000", nil, "6fe9ba44a76be22a", "", "", "md5", "", nil).(*WWWAuthenticate)
	fmt.Print(wa.Raw())
}

func TestWWWAuthenticate_JsonString(t *testing.T) {
	wa := NewWWWAuthenticate(sip.Digest, "3402000000", nil, "6fe9ba44a76be22a", "", "", "md5", "", nil).(*WWWAuthenticate)
	fmt.Println(wa.JsonString())
}

func TestWWWAuthenticate_Parser(t *testing.T) {
	raw := "WWW-Authenticate: Digest realm=\"3402000000\",nonce=\"6fe9ba44a76be22a\",algorithm=MD5\r\n"
	wa := CreateWWWAuthenticate()
	if err := wa.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(wa.Raw())
	fmt.Println(wa.JsonString())
	fmt.Println(raw)
}

func TestWWWAuthenticate_Validator(t *testing.T) {
	wa := NewWWWAuthenticate(sip.Digest, "3402000000", nil, "6fe9ba44a76be22a", "", "", "md5", "", nil).(*WWWAuthenticate)
	fmt.Println(wa.Validator())
}
