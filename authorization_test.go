package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestAuthorization_Raw(t *testing.T) {
	// Authorization: Digest username="bob",
	// 		realm="biloxi.com",
	// 		nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093",
	// 		uri="sip:bob@biloxi.com",
	// 		qop=auth,
	// 		nc=00000001,
	// 		cnonce="0a4f113b",
	// 		response="6629fae49393a05397450978507c4ef1",
	// 		opaque="5ccc069c403ebaf9f0171e9517f40e41"
	authorizations := []*Authorization{
		NewAuthorization("34020000001320000001",
			"3402000000", "",
			NewRequestUri(NewSipUri(NewUserInfo("34020000001320000001", "", ""),
				NewHostPort("", net.IPv4(192, 168, 0, 26), nil, 5060), nil, sync.Map{})),
			"6629fae49393a05397450978507c4ef1",
			"MD5",
			"0a4f113b",
			"5ccc069c403ebaf9f0171e9517f40e41",
			"auth",
			"00000001",
			sync.Map{}),
		NewAuthorization("bob",
			"biloxi.com", "",
			NewRequestUri(NewSipUri(NewUserInfo("bob", "", ""),
				NewHostPort("biloxi.com", nil, nil, 0), nil, sync.Map{})),
			"6629fae49393a05397450978507c4ef1",
			"MD5",
			"0a4f113b",
			"5ccc069c403ebaf9f0171e9517f40e41",
			"auth",
			"00000001",
			sync.Map{}),
	}
	for _, authorization := range authorizations {
		result := authorization.Raw()
		fmt.Print(result.String())

	}
}

func TestAuthorization_Parse(t *testing.T) {
	raws := []string{
		`Authorization: Digest username="34020000001320000001", realm="3402000000", uri="sip:34020000001320000001@192.168.0.26:5060", Response="6629fae49393a05397450978507c4ef1", algorithm=MD5,noNce="cesP", CNonce="0a4f113b", opaque="5ccc069c403ebaf9f0171e9517f40e41", qop=auth, nc=00000001`,
		`Authorization: Digest username="bob", realm="biloxi.com", uri="sip:bob@biloxi.com", response="6629fae49393a05397450978507c4ef1", algorithm=MD5,nonce=456, cnonce="0a4f113b", opaque="5ccc069c403ebaf9f0171e9517f40e41", qop=auth, nc=00000001,hello=word`,
	}
	for index, raw := range raws {
		authorization := new(Authorization)
		authorization.Parse(raw)
		if len(authorization.GetSource()) > 0 {
			fmt.Println(index, "field:", authorization.GetField())
			fmt.Println(index, "auth-schema:", authorization.GetAuthSchema())
			fmt.Println(index, "username:", authorization.GetUsername())
			fmt.Println(index, "realm:", authorization.GetRealm())
			fmt.Println(index, "response:", authorization.GetResponse())
			fmt.Println(index, "algorithm:", authorization.GetAlgorithm())
			fmt.Println(index, "nonce:", authorization.GetNonce())
			fmt.Println(index, "cnonce:", authorization.GetCNonce())
			fmt.Println(index, "opaque:", authorization.GetOpaque())
			fmt.Println(index, "qop:", authorization.GetQop())
			fmt.Println(index, "nc:", authorization.GetNc())
			fmt.Println(index, "uri->request-uri->sip-uri/sips-uri->schema:", authorization.GetUri().GetSipUri().GetSchema())
			fmt.Println(index, "uri->request-uri->sip-uri/sips-uri->userinfo->username:", authorization.GetUri().GetSipUri().GetUserInfo().GetUser())
			authorization.authParam.Range(func(key, value interface{}) bool {
				fmt.Println(index, "auth-param:", key, "=", value)
				return true
			})
			authorization.authParam.Store("hello", "www.baidu.com")
			result := authorization.Raw()
			fmt.Print(result.String())
		}
	}
}
