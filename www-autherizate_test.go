package sip

import (
	"fmt"
	"sync"
	"testing"
)

func TestWWWAuthenticate_Raw(t *testing.T) {
	was := []*WWWAuthenticate{
		NewWWWAuthenticate("3402000001", "", "69c1ad64c2e5323a883be2469838589ce", "", false, "MD5", "auth", sync.Map{}),
		NewWWWAuthenticate("3402000001", "", "f9e3df022ed622c0f886b9e2d0dad507", "", false, "MD5", "auth", sync.Map{}),
		NewWWWAuthenticate("3402000001", "", GetNonce("192.168.124.29", "ZRJOgEycUtwwPBSncBTPgElUUemRsiIJ"), "", false, "MD5", "auth", sync.Map{}),
		NewWWWAuthenticate("3402000000", "", GetNonce("192.168.0.1", "call-id"), "", false, "MD5", "auth", sync.Map{}),
		NewWWWAuthenticate("3402000000", "", GetNonce("192.168.0.1", "call-id"), "", false, "MD5", "auth", sync.Map{}),
	}
	for _, wa := range was {

		result := wa.Raw()
		fmt.Print(result.String())
	}
}

func TestWWWAuthenticate_Parse(t *testing.T) {
	raws := []string{
		`WWW-Authenticate: Digest realm="3402000001", nonce="69c1ad64c2e5323a883be2469838589ce", algorithm=MD5, qop="auth"`,
		`WWW-Authenticate: Digest realm="3402000001", nonce="f9e3df022ed622c0f886b9e2d0dad507", algorithm=MD5, qop="auth"`,
		`WWW-Authenticate: Digest realm="3402000001", nonce="974b9f3e5ff619b08e7e944bc2165736", algorithm=MD5, qop="auth"`,
		`WWW-Authenticate: Digest realm="3402000000", nonce="f899c760a1ee1a92a3e3ec85a5fc1a64", algorithm=MD5, qop="auth"`,
		`WWW-Authenticate: Digest realm="3402000000", nonce="3f7aa6601db1c95f9675d6a1889d2295", algorithm=MD5, qop="auth"`,
	}
	for _, raw := range raws {
		wa := new(WWWAuthenticate)
		wa.Parse(raw)
		if len(wa.GetSource()) > 0 {
			fmt.Println(wa.GetField(), wa.GetAuthSchema(), wa.GetRealm(), wa.GetDomain(), wa.GetNonce(), wa.GetStale(), wa.GetQop(), wa.GetAlgorithm())
			result := wa.Raw()
			fmt.Print(result.String())
		}
	}
}
