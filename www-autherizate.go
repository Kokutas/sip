package sip

import (
	"fmt"
	"reflect"
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
// https://www.rfc-editor.org/rfc/rfc2069.html#section-2.1.1
//
// 2.1.1 The WWW-Authenticate Response Header

// If a server receives a request for an access-protected object, and an
// acceptable Authorization header is not sent, the server responds with
// a "401 Unauthorized" status code, and a WWW-Authenticate header,
// which is defined as follows:

// 	WWW-Authenticate    = "WWW-Authenticate" ":" "Digest"
// 							digest-challenge

// 	digest-challenge    = 1#( realm | [ domain ] | nonce |
// 						[ digest-opaque ] |[ stale ] | [ algorithm ] )

// 	realm               = "realm" "=" realm-value
// 	realm-value         = quoted-string
// 	domain              = "domain" "=" <"> 1#URI <">
// 	nonce               = "nonce" "=" nonce-value
// 	nonce-value         = quoted-string
// 	opaque              = "opaque" "=" quoted-string
// 	stale               = "stale" "=" ( "true" | "false" )
// 	algorithm           = "algorithm" "=" ( "MD5" | token )

// The meanings of the values of the parameters used above are as
// follows:

// 	realm
// 	A string to be displayed to users so they know which username and
// 	password to use.  This string should contain at least the name of
// 	the host performing the authentication and might additionally
// 	indicate the collection of users who might have access.  An example
// 	might be "registered_users@gotham.news.com".  The realm is a
// 	"quoted-string" as specified in section 2.2 of the HTTP/1.1
// 	specification [2].

// 	domain
// 	A comma-separated list of URIs, as specified for HTTP/1.0.  The
// 	intent is that the client could use this information to know the
// 	set of URIs for which the same authentication information should be
// 	sent.  The URIs in this list may exist on different servers.  If
// 	this keyword is omitted or empty, the client should assume that the
// 	domain consists of all URIs on the responding server.

// 	nonce
// 	A server-specified data string which may be uniquely generated each
// 	time a 401 response is made.  It is recommended that this string be
// 	base64 or hexadecimal data.  Specifically, since the string is
// 	passed in the header lines as a quoted string, the double-quote
// 	character is not allowed.

// 	The contents of the nonce are implementation dependent.  The
// 	quality of the implementation depends on a good choice.  A
// 	recommended nonce would include

// 			H(client-IP ":" time-stamp ":" private-key )

// 	Where client-IP is the dotted quad IP address of the client making
// 	the request, time-stamp is a server-generated time value,  private-
// 	key is data known only to the server.  With a nonce of this form a
// 	server would normally recalculate the nonce after receiving the
// 	client authentication header and reject the request if it did not
// 	match the nonce from that header. In this way the server can limit
// 	the reuse of a nonce to the IP address to which it was issued and
// 	limit the time of the nonce's validity.  Further discussion of the
// 	rationale for nonce construction is in section 3.2 below.

// 	An implementation might choose not to accept a previously used
// 	nonce or a previously used digest to protect against a replay
// 	attack.  Or, an implementation might choose to use one-time nonces
// 	or digests for POST or PUT requests and a time-stamp for GET
// 	requests.  For more details on the issues involved see section 3.
// 	of this document.

// 	The nonce is opaque to the client.

// 	opaque
// 	A string of data, specified by the server, which should be
// 	returned by the client unchanged.  It is recommended that this
// 	string be base64 or hexadecimal data.  This field is a
// 	"quoted-string" as specified in section 2.2 of the HTTP/1.1
// 	specification [2].

// 	stale
// 	A flag, indicating that the previous request from the client was
// 	rejected because the nonce value was stale.  If stale is TRUE (in
// 	upper or lower case), the client may wish to simply retry the
// 	request with a new encrypted response, without reprompting the
// 	user for a new username and password.  The server should only set
// 	stale to true if it receives a request for which the nonce is
// 	invalid but with a valid digest for that nonce (indicating that
// 	the client knows the correct username/password).

// 	 algorithm
//      A string indicating a pair of algorithms used to produce the
//      digest and a checksum.  If this not present it is assumed to be
//      "MD5". In this document the string obtained by applying the
//      digest algorithm to the data "data" with secret "secret" will be
//      denoted by KD(secret, data), and the string obtained by applying
//      the checksum algorithm to the data "data" will be denoted
//      H(data).

//      For the "MD5" algorithm

//         H(data) = MD5(data)

//      and

//         KD(secret, data) = H(concat(secret, ":", data))

//      i.e., the digest is the MD5 of the secret concatenated with a colon
//      concatenated with the data.

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

	//  "WWW-Authenticate"
	if len(strings.TrimSpace(wa.field)) == 0 {
		wa.field = "WWW-Authenticate"
		result.WriteString(fmt.Sprintf("%s:", strings.Title(wa.field)))
	} else {
		result.WriteString(fmt.Sprintf("%s:", wa.field))
	}
	// auth-schema: Basic / Digest
	if len(strings.TrimSpace(wa.authSchema)) == 0 {
		wa.authSchema = "Digest"
		result.WriteString(fmt.Sprintf(" %s", strings.Title(wa.authSchema)))
	} else {
		result.WriteString(fmt.Sprintf(" %s", wa.authSchema))
	}
	if wa.isOrder {
		wa.isOrder = false
		for orders := range wa.order {
			if regexp.MustCompile(`((?i)(realm))( )*=`).MatchString(orders) {
				// realm = "realm" EQUAL realm-value,realm-value = quoted-string
				if len(strings.TrimSpace(wa.realm)) > 0 {
					result.WriteString(fmt.Sprintf(" realm=\"%s\",", wa.realm))
					continue
				}
			}
			if regexp.MustCompile(`((?i)(domain))( )*=`).MatchString(orders) {
				// domain =  "domain" EQUAL LDQUOT URI,*( 1*SP URI ) RDQUOT, URI =  absoluteURI / abs-path
				if len(strings.TrimSpace(wa.domain)) > 0 {
					result.WriteString(fmt.Sprintf(" domain=\"%s\",", wa.domain))
					continue
				}
			}
			if regexp.MustCompile(`((?i)(nonce))( )*=`).MatchString(orders) {
				// nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
				if len(strings.TrimSpace(wa.nonce)) > 0 {
					result.WriteString(fmt.Sprintf(" nonce=\"%s\",", wa.nonce))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(opaque))( )*=`).MatchString(orders) {
				// opaque =  "opaque" EQUAL quoted-string
				if len(strings.TrimSpace(wa.opaque)) > 0 {
					result.WriteString(fmt.Sprintf(" opaque=\"%s\",", wa.opaque))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(stale))( )*=`).MatchString(orders) {
				// stale =  "stale" EQUAL ( "true" / "false" )
				if wa.stale {
					result.WriteString(fmt.Sprintf(" stale=\"%v\",", wa.stale))
				}
				continue
			}

			if regexp.MustCompile(`((?i)(algorithm))( )*=`).MatchString(orders) {
				// algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
				if len(strings.TrimSpace(wa.algorithm)) > 0 {
					result.WriteString(fmt.Sprintf(" algorithm=%s,", wa.algorithm))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(qop))( )*=`).MatchString(orders) {
				// qop-options =  "qop" EQUAL LDQUOT qop-value,*("," qop-value) RDQUOT,qop-value =  "auth" / "auth-int" / token
				if len(strings.TrimSpace(wa.qop)) > 0 {
					result.WriteString(fmt.Sprintf(" qop=\"%s\",", wa.qop))
				}
				continue
			}
			ordersSlice := strings.Split(orders, "=")
			if len(ordersSlice) == 1 {
				if val, ok := wa.authParam.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf("  %v=\"%v\",", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(" %v,", ordersSlice[0]))
					}
				} else {
					result.WriteString(fmt.Sprintf(" %v,", ordersSlice[0]))
				}
			} else {
				if val, ok := wa.authParam.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf(" %v=\"%v\",", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(" %v,", ordersSlice[0]))
					}

				} else {
					if len(strings.TrimSpace(fmt.Sprintf("%v", ordersSlice[1]))) > 0 {
						result.WriteString(fmt.Sprintf("  %v=\"%v\",", ordersSlice[0], ordersSlice[1]))
					} else {
						result.WriteString(fmt.Sprintf(" %v,", ordersSlice[0]))
					}
				}
			}
		}

	} else {
		// realm = "realm" EQUAL realm-value,realm-value = quoted-string
		if len(strings.TrimSpace(wa.realm)) > 0 {
			result.WriteString(fmt.Sprintf(" realm=\"%s\",", wa.realm))
		}
		// domain =  "domain" EQUAL LDQUOT URI,*( 1*SP URI ) RDQUOT, URI =  absoluteURI / abs-path
		if len(strings.TrimSpace(wa.domain)) > 0 {
			result.WriteString(fmt.Sprintf(" domain=\"%s\",", wa.domain))
		}
		// nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
		if len(strings.TrimSpace(wa.nonce)) > 0 {
			result.WriteString(fmt.Sprintf(" nonce=\"%s\",", wa.nonce))
		}
		// opaque =  "opaque" EQUAL quoted-string
		if len(strings.TrimSpace(wa.opaque)) > 0 {
			result.WriteString(fmt.Sprintf(" opaque=\"%s\",", wa.opaque))
		}
		// stale =  "stale" EQUAL ( "true" / "false" )
		if wa.stale {
			result.WriteString(fmt.Sprintf(" stale=\"%v\",", wa.stale))
		}
		// algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
		if len(strings.TrimSpace(wa.algorithm)) > 0 {
			result.WriteString(fmt.Sprintf(" algorithm=%s,", wa.algorithm))
		}
		// qop-options =  "qop" EQUAL LDQUOT qop-value,*("," qop-value) RDQUOT,qop-value =  "auth" / "auth-int" / token
		if len(strings.TrimSpace(wa.qop)) > 0 {
			result.WriteString(fmt.Sprintf(" qop=\"%s\",", wa.qop))
		}
	}
	// auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
	wa.authParam.Range(func(key, value interface{}) bool {
		if reflect.ValueOf(value).IsValid() {
			if reflect.ValueOf(value).IsZero() {
				result.WriteString(fmt.Sprintf(" %v,", key))
				return true
			}
			result.WriteString(fmt.Sprintf(" %v=\"%v\",", key, value))
			return true
		}
		result.WriteString(fmt.Sprintf(" %v,", key))
		return true
	})
	temp := result.String()
	temp = strings.TrimSuffix(temp, ",")
	result.Reset()
	result.WriteString(temp)
	result.WriteString("\r\n")
	return
}
func (wa *WWWAuthenticate) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`((?i)(?:^www-authenticate))( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	wa.source = raw
	wa.authParam = sync.Map{}
	wa.algorithm = "MD5"
	field := fieldRegexp.FindString(raw)
	field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	wa.field = field
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// auth-schema regexp
	authschemaRegexp := regexp.MustCompile(`(?i)(basic|digest)`)
	if authschemaRegexp.MatchString(raw) {
		authschema := authschemaRegexp.FindString(raw)
		wa.authSchema = authschema
		raw = authschemaRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ",")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}

	// www-authenticate order
	wa.wwwAuthenticateOrder(raw)

	// realm regexp : realm = "realm" EQUAL realm-value,realm-value = quoted-string
	realmRegexp := regexp.MustCompile(`((?i)(?:^realm))( )*=`)
	// domain regexp ：domain =  "domain" EQUAL LDQUOT URI,*( 1*SP URI ) RDQUOT, URI =  absoluteURI / abs-path
	domainRegexp := regexp.MustCompile(`((?i)(?:^domain))( )*=`)
	// nonce regexp : nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
	nonceRegexp := regexp.MustCompile(`((?i)(?:^nonce))( )*=`)
	// opaque     regexp opaque =  "opaque" EQUAL quoted-string
	opaqueRegexp := regexp.MustCompile(`((?i)(?:^opaque))( )*=`)
	// stale =  "stale" EQUAL ( "true" / "false" )
	staleRegexp := regexp.MustCompile(`((?i)(?:^stale))( )*=`)
	// algorithm regexp algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
	algorithmRegexp := regexp.MustCompile(`((?i)(?:^algorithm))( )*=`)
	// qop        regexp message-qop = "qop" EQUAL qop-value,qop-value = "auth" / "auth-int" / token
	qopRegexp := regexp.MustCompile(`((?i)(?:^qop))( )*=`)

	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		raws = stringTrimPrefixAndTrimSuffix(raws, " ")
		switch {
		case realmRegexp.MatchString(raws):
			realm := realmRegexp.ReplaceAllString(raws, "")
			realm = regexp.MustCompile(`"`).ReplaceAllString(realm, "")
			wa.realm = realm
		case domainRegexp.MatchString(raws):
			domain := domainRegexp.ReplaceAllString(raw, "")
			domain = regexp.MustCompile(`"`).ReplaceAllString(domain, "")
			wa.domain = domain
		case nonceRegexp.MatchString(raws):
			nonce := nonceRegexp.ReplaceAllString(raws, "")
			nonce = regexp.MustCompile(`"`).ReplaceAllString(nonce, "")
			wa.nonce = nonce
		case opaqueRegexp.MatchString(raws):
			opaque := opaqueRegexp.ReplaceAllString(raws, "")
			opaque = regexp.MustCompile(`"`).ReplaceAllString(opaque, "")
			wa.opaque = opaque
		case staleRegexp.MatchString(raws):
			stale := opaqueRegexp.ReplaceAllString(raws, "")
			stale = regexp.MustCompile(`"`).ReplaceAllString(stale, "")
			if regexp.MustCompile(`(?i)(true)`).MatchString(stale) {
				wa.stale = true
			}
		case algorithmRegexp.MatchString(raws):
			algorithm := algorithmRegexp.ReplaceAllString(raws, "")
			algorithm = regexp.MustCompile(`"`).ReplaceAllString(algorithm, "")
			wa.algorithm = algorithm
		case qopRegexp.MatchString(raws):
			qop := qopRegexp.ReplaceAllString(raws, "")
			qop = regexp.MustCompile(`"`).ReplaceAllString(qop, "")
			wa.qop = qop
		default:
			// authParam  sync.Map    // auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
			kvs := strings.Split(raws, "=")
			if len(kvs) == 1 {
				wa.authParam.Store(kvs[0], "")
			} else {
				wa.authParam.Store(kvs[0], kvs[1])
			}
		}
	}

}
func (wa *WWWAuthenticate) wwwAuthenticateOrder(raw string) {
	wa.isOrder = true
	wa.order = make(chan string, 1024)
	defer close(wa.order)
	raw = stringTrimPrefixAndTrimSuffix(raw, ",")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		wa.order <- raws
	}
}
