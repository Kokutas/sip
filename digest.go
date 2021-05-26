package sip

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"time"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-22.4
//

// 22.4 The Digest Authentication Scheme

// This section describes the modifications and clarifications required
// to apply the HTTP Digest authentication scheme to SIP.  The SIP
// scheme usage is almost completely identical to that for HTTP [17].

// Since RFC 2543 is based on HTTP Digest as defined in RFC 2069 [39],
// SIP servers supporting RFC 2617 MUST ensure they are backwards
// compatible with RFC 2069.  Procedures for this backwards
// compatibility are specified in RFC 2617.  Note, however, that SIP
// servers MUST NOT accept or request Basic authentication.

// The rules for Digest authentication follow those defined in [17],
// with "HTTP/1.1" replaced by "SIP/2.0" in addition to the following
// differences:

// 	1.  The URI included in the challenge has the following BNF:

// 		URI  =  SIP-URI / SIPS-URI

// 	2.  The BNF in RFC 2617 has an error in that the 'uri' parameter
// 		of the Authorization header field for HTTP Digest
// 		authentication is not enclosed in quotation marks.  (The
// 			example in Section 3.5 of RFC 2617 is correct.)  For SIP, the
// 			'uri' MUST be enclosed in quotation marks.

// 	3.  The BNF for digest-uri-value is:

// 		digest-uri-value  =  Request-URI ; as defined in Section 25

// 	4.  The example procedure for choosing a nonce based on Etag does
// 		not work for SIP.

// 	5.  The text in RFC 2617 [17] regarding cache operation does not
// 		apply to SIP.

// 	6.  RFC 2617 [17] requires that a server check that the URI in the
// 		request line and the URI included in the Authorization header
// 		field point to the same resource.  In a SIP context, these two
// 		URIs may refer to different users, due to forwarding at some
// 		proxy.  Therefore, in SIP, a server MAY check that the
// 		Request-URI in the Authorization header field value
// 		corresponds to a user for whom the server is willing to accept
// 		forwarded or direct requests, but it is not necessarily a
// 		failure if the two fields are not equivalent.

// 	7.  As a clarification to the calculation of the A2 value for
// 		message integrity assurance in the Digest authentication
// 		scheme, implementers should assume, when the entity-body is
// 		empty (that is, when SIP messages have no body) that the hash
// 		of the entity-body resolves to the MD5 hash of an empty
// 		string, or:

// 			H(entity-body) = MD5("") =
// 		"d41d8cd98f00b204e9800998ecf8427e"

// 	8.  RFC 2617 notes that a cnonce value MUST NOT be sent in an
// 		Authorization (and by extension Proxy-Authorization) header
// 		field if no qop directive has been sent.  Therefore, any
// 		algorithms that have a dependency on the cnonce (including
// 		"MD5-Sess") require that the qop directive be sent.  Use of
// 		the "qop" parameter is optional in RFC 2617 for the purposes
// 		of backwards compatibility with RFC 2069; since RFC 2543 was
// 		based on RFC 2069, the "qop" parameter must unfortunately
// 		remain optional for clients and servers to receive.  However,
// 		servers MUST always send a "qop" parameter in WWW-Authenticate
// 		and Proxy-Authenticate header field values.  If a client
// 		receives a "qop" parameter in a challenge header field, it
// 		MUST send the "qop" parameter in any resulting authorization
// 		header field.
// RFC 2543 did not allow usage of the Authentication-Info header field
// (it effectively used RFC 2069).  However, we now allow usage of this
// header field, since it provides integrity checks over the bodies and
// provides mutual authentication.  RFC 2617 [17] defines mechanisms for
// backwards compatibility using the qop attribute in the request.
// These mechanisms MUST be used by a server to determine if the client
// supports the new mechanisms in RFC 2617 that were not specified in
// RFC 2069.

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

	if p.Qop == "auth-int" {
		bytes = md5.Sum([]byte(fmt.Sprintf("%s:%s%s", p.Method, p.URI, p.EntityBody)))
	} else {
		bytes = md5.Sum([]byte(fmt.Sprintf("%s:%s", p.Method, p.URI)))
	}
	ha2 := fmt.Sprintf("%x", bytes)
	if p.Qop == "" {
		bytes = md5.Sum([]byte(ha1 + ":" + p.Nonce + ":" + ha2))
	} else {
		bytes = md5.Sum([]byte(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", ha1, p.Nonce, p.Nc, p.Cnonce, p.Qop, ha2)))
	}
	p.Response = fmt.Sprintf("%x", bytes)
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

// H(client-IP ":" time-stamp ":" private-key )
func GenNonce(clientIP string, privateKey string) string {
	bytes := md5.Sum([]byte(fmt.Sprintf("%v:%v:%v", clientIP, time.Now().UnixNano(), privateKey)))
	return fmt.Sprintf("%x", bytes)
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
