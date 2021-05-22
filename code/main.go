package main

import (
	"fmt"
	"regexp"
)

func main() {
	raw := `Authorization: Digest username="bob", realm="biloxi.com", uri="sip:bob@biloxi.com", response="6629fae49393a05397450978507c4ef1", algorithm=MD5, cnonce="0a4f113b", opaque="5ccc069c403ebaf9f0171e9517f40e41", qop=auth, nc=00000001,hello=word`
	fmt.Println(regexp.MustCompile(`((?i)(response))( )*=`).FindString(raw))
}
