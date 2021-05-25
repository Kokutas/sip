package sip

import (
	"regexp"
	"strings"
	"sync"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.44
//
// 20.44 WWW-Authenticate
//
// A WWW-Authenticate header field value contains an authentication
// challenge.  See Section 22.2 for further details on its usage.

// Example:

// 	WWW-Authenticate: Digest realm="atlanta.com",
// 					domain="sip:boxesbybob.com", qop="auth",
// 					nonce="f84f1cec41e6cbe5aea9c8e88d359",
// 					opaque="", stale=FALSE, algorithm=MD5
//
// WWW-Authenticate: Digest
//               realm="biloxi.com",
//               qop="auth,auth-int",
//               nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093",
//               opaque="5ccc069c403ebaf9f0171e9517f40e41"

//
//
//https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// WWW-Authenticate  =  "WWW-Authenticate" HCOLON challenge
// challenge           =  ("Digest" LWS digest-cln *(COMMA digest-cln))
//                        / other-challenge
// other-challenge     =  auth-scheme LWS auth-param
//                        *(COMMA auth-param)
// digest-cln          =  realm / domain / nonce
//                         / opaque / stale / algorithm
//                         / qop-options / auth-param
// realm               =  "realm" EQUAL realm-value
// realm-value         =  quoted-string
// domain              =  "domain" EQUAL LDQUOT URI
//                        *( 1*SP URI ) RDQUOT
// URI                 =  absoluteURI / abs-path
// nonce               =  "nonce" EQUAL nonce-value
// nonce-value         =  quoted-string
// opaque              =  "opaque" EQUAL quoted-string
// stale               =  "stale" EQUAL ( "true" / "false" )
// algorithm           =  "algorithm" EQUAL ( "MD5" / "MD5-sess"
//                        / token )
// qop-options         =  "qop" EQUAL LDQUOT qop-value
//                        *("," qop-value) RDQUOT
// qop-value           =  "auth" / "auth-int" / token
//
type WWWAuthenticate struct {
	field      string      // "WWW-Authenticate"
	authSchema string      // auth-schema: Basic / Digest
	realm      string      // realm =  "realm" EQUAL realm-value,realm-value =  quoted-string
	domain     string      // domain =  "domain" EQUAL LDQUOT URI,*( 1*SP URI ) RDQUOT, URI =  absoluteURI / abs-path
	nonce      string      // nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
	opaque     string      // opaque =  "opaque" EQUAL quoted-string
	stale      bool        // stale =  "stale" EQUAL ( "true" / "false" )
	algorithm  string      // algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
	qop        string      // qop-options =  "qop" EQUAL LDQUOT qop-value,*("," qop-value) RDQUOT,qop-value =  "auth" / "auth-int" / token
	authParam  sync.Map    // auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
	isOrder    bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order      chan string // It is convenient to record the order of the original parameter fields when parsing
	source     string      // source string
}

// "WWW-Authenticate"
func (wa *WWWAuthenticate) SetField(field string) {
	if regexp.MustCompile(`^(?i)(www-authenticate)$`).MatchString(field) {
		wa.field = strings.Title(field)
	} else {
		wa.field = "WWW-Authenticate"
	}
}
func (wa *WWWAuthenticate) GetField() string {
	return wa.field
}

// auth-schema: Basic / Digest
func (wa *WWWAuthenticate) SetAuthSchema(authSchema string) {
	if regexp.MustCompile(`(?i)(basic|digest)`).MatchString(authSchema) {
		wa.authSchema = strings.Title(authSchema)
	}
	wa.authSchema = "Digest"
}
func (wa *WWWAuthenticate) GetAuthSchema() string {
	return wa.authSchema
}

// realm = "realm" EQUAL realm-value,realm-value = quoted-string
func (wa *WWWAuthenticate) SetRealm(realm string) {
	wa.realm = realm
}
func (wa *WWWAuthenticate) GetRealm() string {
	return wa.realm
}

// domain =  "domain" EQUAL LDQUOT URI,*( 1*SP URI ) RDQUOT, URI =  absoluteURI / abs-path
func (wa *WWWAuthenticate) SetDomain(domain string) {
	wa.domain = domain
}
func (wa *WWWAuthenticate) GetDomain() string {
	return wa.domain
}

// nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
func (wa *WWWAuthenticate) SetNonce(nonce string) {
	wa.nonce = nonce
}
func (wa *WWWAuthenticate) GetNonce() string {
	return wa.nonce
}

// opaque =  "opaque" EQUAL quoted-string
func (wa *WWWAuthenticate) SetOpaque(opaque string) {
	wa.opaque = opaque
}
func (wa *WWWAuthenticate) GetOpaque() string {
	return wa.opaque
}

// stale =  "stale" EQUAL ( "true" / "false" )
func (wa *WWWAuthenticate) SetStale(stale bool) {
	wa.stale = stale
}
func (wa *WWWAuthenticate) GetStale() bool {
	return wa.stale
}

// algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
func (wa *WWWAuthenticate) SetAlgorithm(algorithm string) {
	wa.algorithm = algorithm
}
func (wa *WWWAuthenticate) GetAlgorithm() string {
	return wa.algorithm
}

// qop-options =  "qop" EQUAL LDQUOT qop-value,*("," qop-value) RDQUOT,qop-value =  "auth" / "auth-int" / token
func (wa *WWWAuthenticate) SetQop(qop string) {
	wa.qop = qop
}
func (wa *WWWAuthenticate) GetQop() string {
	return wa.qop
}

// auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
func (wa *WWWAuthenticate) SetAuthParam(authParam sync.Map) {
	wa.authParam = authParam
}
func (wa *WWWAuthenticate) GetAuthParam() sync.Map {
	return wa.authParam
}

// source string
func (wa *WWWAuthenticate) GetSource() string {
	return wa.source
}
func NewWWWAuthenticate(realm string, domain string, nonce string, opaque string, stale bool, algorithm string, qop string, authParam sync.Map) *WWWAuthenticate {
	return &WWWAuthenticate{
		field:      "WWW-Authenticate",
		authSchema: "Digest",
		realm:      realm,
		domain:     domain,
		nonce:      nonce,
		opaque:     opaque,
		stale:      stale,
		algorithm:  algorithm,
		qop:        qop,
		authParam:  authParam,
		isOrder:    false,
	}
}
func (wa *WWWAuthenticate) Raw() (result strings.Builder) {
	return
}
func (wa *WWWAuthenticate) Parse(raw string) {

}
func (wa *WWWAuthenticate) wwwAuthenticateOrder(raw string) {
	wa.isOrder = true
	wa.order = make(chan string, 1024)
	defer close(wa.order)

}
