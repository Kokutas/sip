package auth

import (
	"crypto/md5"
	"fmt"
	"log"
	"regexp"
)

type Digest struct {
	Realm    string
	UserName string
	Password string
}
type DigestParams struct {
	Digest
	Qop        string
	Qops       []string
	Algorithm  string
	Method     string
	URI        string
	Nonce      string
	Cnonce     string
	Nc         uint32
	EntityBody string
	Response   string
}

// if the algorithm directive's value is "MD5" or unspecified ,then HA1 is : HA1=MD5(username:realm:password)
// if the algorithm directive's value is "MD5-sess" , then HA1 is : HA1=MD5(MD5(username:realm:password):nonce:cnonce)
// if the qop directive's  value is "auth" or unspecified, then HA2 is : HA2=MD5(method:digest-uri)
// if the qop directive's value is "auth-int" , them HA2 is : HA2=MD5(method:digest-uri:MD5(entity-body)
// if the qop directive's value is "auth" or "auth-int" , then compute the response is : response=MD5(HA1:nonce:nonce-count:cnonce:qop:HA2)
// if the qop directive is unspecified , then compute the response  is : response=MD5(HA1:nonce:HA2)
// The above shows that when qop is not specified , the simpler RFC 2069 standard is followed

func GenDigestResponse(p *DigestParams) string {
	bytes := md5.Sum([]byte(p.Digest.UserName + ":" + p.Digest.Realm + ":" + p.Digest.Password))
	if p.Algorithm == "MD5-sess" {
		b := append(bytes[:], []byte(":"+p.Nonce+":"+p.Cnonce)...)
		bytes = md5.Sum(b)
	}
	ha1 := fmt.Sprintf("%x", bytes)
	log.Printf("HA1: %s\r\n", ha1)

	if p.Qop == "auth-int" {
		bytes = md5.Sum([]byte(fmt.Sprintf("%s:%s%s", p.Method, p.URI, p.EntityBody)))
	} else {
		bytes = md5.Sum([]byte(fmt.Sprintf("%s:%s", p.Method, p.URI)))
	}
	ha2 := fmt.Sprintf("%x", bytes)
	log.Printf("HA2: %s\r\n", ha2)
	if p.Qop == "" {
		bytes = md5.Sum([]byte(ha1 + ":" + p.Nonce + ":" + ha2))
	} else {
		bytes = md5.Sum([]byte(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", ha1, p.Nonce, p.Nc, p.Cnonce, p.Qop, ha2)))
	}
	p.Response = fmt.Sprintf("%x", bytes)
	log.Printf("response : %s\r\n", p.Response)
	return p.Response
}

func DigestCalculatorResponse(username, realm, password, nonce, uri string) []string {
	responses := make([]string, 0)
	response1 := getDigestResponse(username, realm, password, nonce, uri)
	response2 := getDigestResponse(username, username[:10], password, nonce, uri)
	response3 := getDigestResponse(username, username[:10], password, nonce, regexp.MustCompile("@.*").ReplaceAllString(uri, "@"+realm))
	response4 := getDigestResponse(username, realm, password, nonce, regexp.MustCompile("@.*").ReplaceAllString(uri, "@"+realm))
	responses = append(responses, response1, response2, response3, response4)
	return responses
}
func getDigestResponse(username, realm, password, nonce, uri string) string {
	dp := &DigestParams{
		Digest: Digest{
			Realm:    realm,
			UserName: username,
			Password: password,
		},
		Qop:       "",
		Algorithm: "MD5",
		Method:    "REGISTER",
		URI:       uri,
		Nonce:     nonce,
		Cnonce:    "",
		//Nc: 1,
	}
	response := GenDigestResponse(dp)
	return response
}
func GetDigestNonce(username, realm, password, uri, callid string) string {
	dp := &DigestParams{
		Digest: Digest{
			Realm:    realm,
			UserName: username,
			Password: password,
		},
		Qop:       "",
		Algorithm: "MD5",
		Method:    "REGISTER",
		URI:       uri,
		Cnonce:    "",
		Nc:        1,
	}
	dp.Digest.UserName = username
	bytes := md5.Sum([]byte(callid + username))
	dp.Nonce = fmt.Sprintf("%x", bytes)
	return dp.Nonce
}
