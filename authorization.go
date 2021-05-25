package sip

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// https://www.rfc-editor.org/rfc/rfc3261#section-20.7

// 20.7 Authorization
// The Authorization header field contains authentication credentials of
// a UA.  Section 22.2 overviews the use of the Authorization header
// field, and Section 22.4 describes the syntax and semantics when used
// with HTTP authentication.

// This header field, along with Proxy-Authorization, breaks the general
// rules about multiple header field values.  Although not a comma-
// separated list, this header field name may be present multiple times,
// and MUST NOT be combined into a single header line using the usual
// rules described in Section 7.3.
// In the example below, there are no quotes around the Digest
// parameter:

// 	Authorization: Digest username="Alice", realm="atlanta.com",
// 	nonce="84a4cc6f3082121f32b42a2187831a9e",
// 	response="7587245234b3434cc3412213e5f113a5432"

// https://www.rfc-editor.org/rfc/rfc3261#section-25.1
//
// Authorization     =  "Authorization" HCOLON credentials
// credentials       =  ("Digest" LWS digest-response)
//                      / other-response
// digest-response   =  dig-resp *(COMMA dig-resp)
// dig-resp          =  username / realm / nonce / digest-uri
//                       / dresponse / algorithm / cnonce
//                       / opaque / message-qop
//                       / nonce-count / auth-param
// username          =  "username" EQUAL username-value
// username-value    =  quoted-string
// digest-uri        =  "uri" EQUAL LDQUOT digest-uri-value RDQUOT
// digest-uri-value  =  rquest-uri ; Equal to request-uri as specified
//                      by HTTP/1.1
// message-qop       =  "qop" EQUAL qop-value
// cnonce            =  "cnonce" EQUAL cnonce-value
// cnonce-value      =  nonce-value
// nonce-count       =  "nc" EQUAL nc-value
// nc-value          =  8LHEX
// dresponse         =  "response" EQUAL request-digest
// request-digest    =  LDQUOT 32LHEX RDQUOT
// auth-param        =  auth-param-name EQUAL
//                      ( token / quoted-string )
// auth-param-name   =  token
// other-response    =  auth-scheme LWS auth-param
//                      *(COMMA auth-param)
// auth-scheme       =  token

// An example of the Authorization header field is:

// Authorization: Digest username="bob",
// 		realm="biloxi.com",
// 		nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093",
// 		uri="sip:bob@biloxi.com",
// 		qop=auth,
// 		nc=00000001,
// 		cnonce="0a4f113b",
// 		response="6629fae49393a05397450978507c4ef1",
// 		opaque="5ccc069c403ebaf9f0171e9517f40e41"

// https://www.rfc-editor.org/rfc/rfc3261#section-22.4
//

// 22.4 The Digest Authentication Scheme
//
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

// 1.  The URI included in the challenge has the following BNF:

// 	URI  =  SIP-URI / SIPS-URI

// 2.  The BNF in RFC 2617 has an error in that the 'uri' parameter
// 	of the Authorization header field for HTTP Digest
// 	authentication is not enclosed in quotation marks.  (The
// 	example in Section 3.5 of RFC 2617 is correct.)  For SIP, the
// 	'uri' MUST be enclosed in quotation marks.

// 3.  The BNF for digest-uri-value is:

// 	digest-uri-value  =  Request-URI ; as defined in Section 25

// 4.  The example procedure for choosing a nonce based on Etag does
// 	not work for SIP.

// 5.  The text in RFC 2617 [17] regarding cache operation does not
// 	apply to SIP.

// 6.  RFC 2617 [17] requires that a server check that the URI in the
// 	request line and the URI included in the Authorization header
// 	field point to the same resource.  In a SIP context, these two
// 	URIs may refer to different users, due to forwarding at some
// 	proxy.  Therefore, in SIP, a server MAY check that the
// 	Request-URI in the Authorization header field value
// 	corresponds to a user for whom the server is willing to accept
// 	forwarded or direct requests, but it is not necessarily a
// 	failure if the two fields are not equivalent.

// 7.  As a clarification to the calculation of the A2 value for
// 	message integrity assurance in the Digest authentication
// 	scheme, implementers should assume, when the entity-body is
// 	empty (that is, when SIP messages have no body) that the hash
// 	of the entity-body resolves to the MD5 hash of an empty
// 	string, or:

// 		H(entity-body) = MD5("") =
// 	"d41d8cd98f00b204e9800998ecf8427e"

// 8.  RFC 2617 notes that a cnonce value MUST NOT be sent in an
// 	Authorization (and by extension Proxy-Authorization) header
// 	field if no qop directive has been sent.  Therefore, any
// 	algorithms that have a dependency on the cnonce (including
// 	"MD5-Sess") require that the qop directive be sent.  Use of
// 	the "qop" parameter is optional in RFC 2617 for the purposes
// 	of backwards compatibility with RFC 2069; since RFC 2543 was
// 	based on RFC 2069, the "qop" parameter must unfortunately
// 	remain optional for clients and servers to receive.  However,
// 	servers MUST always send a "qop" parameter in WWW-Authenticate
// 	and Proxy-Authenticate header field values.  If a client
// 	receives a "qop" parameter in a challenge header field, it
// 	MUST send the "qop" parameter in any resulting authorization
// 	header field.
// RFC 2543 did not allow usage of the Authentication-Info header field
// (it effectively used RFC 2069).  However, we now allow usage of this
// header field, since it provides integrity checks over the bodies and
// provides mutual authentication.  RFC 2617 [17] defines mechanisms for
// backwards compatibility using the qop attribute in the request.
// These mechanisms MUST be used by a server to determine if the client
// supports the new mechanisms in RFC 2617 that were not specified in
// RFC 2069.

//
type Authorization struct {
	field      string      // "Authorization"
	authSchema string      // auth-schema: Basic / Digest
	username   string      // username = "username" EQUAL username-value,username-value = quoted-string
	realm      string      // realm = "realm" EQUAL realm-value,realm-value = quoted-string
	nonce      string      // nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
	uri        *RequestUri // digest-uri = "uri" EQUAL LDQUOT digest-uri-value RDQUOT,digest-uri-value = rquest-uri ; Equal to request-uri as specified by HTTP/1.1
	response   string      // dresponse = "response" EQUAL request-digest, request-digest = LDQUOT 32LHEX RDQUOT
	algorithm  string      // algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
	cnonce     string      // cnonce = "cnonce" EQUAL cnonce-value,cnonce-value = nonce-value
	opaque     string      // opaque =  "opaque" EQUAL quoted-string
	qop        string      // message-qop = "qop" EQUAL qop-value,qop-value = "auth" / "auth-int" / token
	nc         string      // nonce-count = "nc" EQUAL nc-value,nc-value = 8LHEX
	authParam  sync.Map    // auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
	isOrder    bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order      chan string // It is convenient to record the order of the original parameter fields when parsing
	source     string      // source string
}

// "Authorization"
func (authorization *Authorization) SetField(field string) {
	if regexp.MustCompile(`^(?i)(authorization)$`).MatchString(field) {
		authorization.field = strings.Title(field)
	} else {
		authorization.field = strings.Title("Authorization")
	}
}
func (authorization *Authorization) GetField() string {
	return authorization.field
}

// auth-schema: Basic / Digest
func (authorization *Authorization) SetAuthSchema(authSchema string) {
	if regexp.MustCompile(`(?i)(basic|digest)`).MatchString(authSchema) {
		authorization.authSchema = strings.Title(authSchema)
	}
	authorization.authSchema = "Digest"
}
func (authorization *Authorization) GetAuthSchema() string {
	return authorization.authSchema
}

// username = "username" EQUAL username-value,username-value = quoted-string
func (authorization *Authorization) SetUsername(username string) {
	authorization.username = username
}
func (authorization *Authorization) GetUsername() string {
	return authorization.username
}

// realm = "realm" EQUAL realm-value,realm-value = quoted-string
func (authorization *Authorization) SetRealm(realm string) {
	authorization.realm = realm
}
func (authorization *Authorization) GetRealm() string {
	return authorization.realm
}

// nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
func (authorization *Authorization) SetNonce(nonce string) {
	authorization.nonce = nonce
}
func (authorization *Authorization) GetNonce() string {
	return authorization.nonce
}

// digest-uri = "uri" EQUAL LDQUOT digest-uri-value RDQUOT,digest-uri-value = rquest-uri ; Equal to request-uri as specified by HTTP/1.1
func (authorization *Authorization) SetUri(uri *RequestUri) {
	authorization.uri = uri
}
func (authorization *Authorization) GetUri() *RequestUri {
	return authorization.uri
}

// dresponse = "response" EQUAL request-digest, request-digest = LDQUOT 32LHEX RDQUOT
func (authorization *Authorization) SetResponse(response string) {
	authorization.response = response
}
func (authorization *Authorization) GetResponse() string {
	return authorization.response
}

// algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
func (authorization *Authorization) SetAlgorithm(algorithm string) {
	authorization.algorithm = algorithm
}
func (authorization *Authorization) GetAlgorithm() string {
	return authorization.algorithm
}

// cnonce = "cnonce" EQUAL cnonce-value,cnonce-value = nonce-value
func (authorization *Authorization) SetCNonce(cnonce string) {
	authorization.cnonce = cnonce
}
func (authorization *Authorization) GetCNonce() string {
	return authorization.cnonce
}

// opaque =  "opaque" EQUAL quoted-string
func (authorization *Authorization) SetOpaque(opaque string) {
	authorization.opaque = opaque
}
func (authorization *Authorization) GetOpaque() string {
	return authorization.opaque
}

// message-qop = "qop" EQUAL qop-value,qop-value = "auth" / "auth-int" / token
func (authorization *Authorization) SetQop(qop string) {
	authorization.qop = qop
}
func (authorization *Authorization) GetQop() string {
	return authorization.qop
}

// nonce-count = "nc" EQUAL nc-value,nc-value = 8LHEX
func (authorization *Authorization) SetNc(nc string) {
	authorization.nc = nc
}
func (authorization *Authorization) GetNc() string {
	return authorization.nc
}

// auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
func (authorization *Authorization) SetAuthParam(authParam sync.Map) {
	authorization.authParam = authParam
}
func (authorization *Authorization) GetAuthParam() sync.Map {
	return authorization.authParam
}

// source string
func (authorization *Authorization) GetSource() string {
	return authorization.source
}
func NewAuthorization(username string, realm string, nonce string, uri *RequestUri, response string, algorithm string, cnonce string, opaque string, qop string, nc string, authParam sync.Map) *Authorization {
	return &Authorization{
		field:      "Authorization",
		authSchema: "Digest",
		username:   username,
		realm:      realm,
		nonce:      nonce,
		uri:        uri,
		response:   response,
		algorithm:  algorithm,
		cnonce:     cnonce,
		opaque:     opaque,
		qop:        qop,
		nc:         nc,
		authParam:  authParam,
		isOrder:    false,
	}
}
func (authorization *Authorization) Raw() (result strings.Builder) {

	// "Authorization"
	if len(strings.TrimSpace(authorization.field)) == 0 {
		authorization.field = "Authorization"
		result.WriteString(fmt.Sprintf("%s:", strings.Title(authorization.field)))
	}
	result.WriteString(fmt.Sprintf("%s:", authorization.field))
	// auth-schema: Basic / Digest
	if len(strings.TrimSpace(authorization.authSchema)) == 0 {
		authorization.authSchema = "Digest"
		result.WriteString(fmt.Sprintf(" %s", strings.Title(authorization.authSchema)))
	}
	result.WriteString(fmt.Sprintf(" %s", authorization.authSchema))

	if authorization.isOrder {
		authorization.isOrder = false
		for orders := range authorization.order {
			if regexp.MustCompile(`((?i)(username))( )*=`).MatchString(orders) {
				// username = "username" EQUAL username-value,username-value = quoted-string
				if len(strings.TrimSpace(authorization.username)) > 0 {
					result.WriteString(fmt.Sprintf(" username=\"%s\",", authorization.username))
					continue
				}
			}
			if regexp.MustCompile(`((?i)(realm))( )*=`).MatchString(orders) {
				// realm = "realm" EQUAL realm-value,realm-value = quoted-string
				if len(strings.TrimSpace(authorization.realm)) > 0 {
					result.WriteString(fmt.Sprintf(" realm=\"%s\",", authorization.realm))
					continue
				}
			}
			if regexp.MustCompile(`((?i)(nonce))( )*=`).MatchString(orders) {
				// nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
				if len(strings.TrimSpace(authorization.nonce)) > 0 {
					result.WriteString(fmt.Sprintf(" nonce=\"%s\",", authorization.nonce))
				}
				continue
			}

			if regexp.MustCompile(`((?i)(uri))( )*=`).MatchString(orders) {
				// digest-uri = "uri" EQUAL LDQUOT digest-uri-value RDQUOT,digest-uri-value = rquest-uri ; Equal to request-uri as specified by HTTP/1.1
				if authorization.uri != nil {
					uri := authorization.uri.Raw()
					result.WriteString(fmt.Sprintf(" uri=\"%s\",", uri.String()))
				}
				continue
			}

			if regexp.MustCompile(`((?i)(response))( )*=`).MatchString(orders) {
				// dresponse = "response" EQUAL request-digest, request-digest = LDQUOT 32LHEX RDQUOT
				if len(strings.TrimSpace(authorization.response)) > 0 {
					result.WriteString(fmt.Sprintf(" response=\"%s\",", authorization.response))
				}
				continue
			}

			if regexp.MustCompile(`((?i)(algorithm))( )*=`).MatchString(orders) {
				// algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
				if len(strings.TrimSpace(authorization.algorithm)) > 0 {
					result.WriteString(fmt.Sprintf(" algorithm=%s,", authorization.algorithm))
				}
				continue
			}

			if regexp.MustCompile(`((?i)(cnonce))( )*=`).MatchString(orders) {
				// cnonce = "cnonce" EQUAL cnonce-value,cnonce-value = nonce-value
				if len(strings.TrimSpace(authorization.cnonce)) > 0 {
					result.WriteString(fmt.Sprintf(" cnonce=\"%s\",", authorization.cnonce))
				}
				continue
			}

			if regexp.MustCompile(`((?i)(opaque))( )*=`).MatchString(orders) {
				// opaque =  "opaque" EQUAL quoted-string
				if len(strings.TrimSpace(authorization.opaque)) > 0 {
					result.WriteString(fmt.Sprintf(" opaque=\"%s\",", authorization.opaque))
				}
				continue
			}

			if regexp.MustCompile(`((?i)(qop))( )*=`).MatchString(orders) {
				// message-qop = "qop" EQUAL qop-value,qop-value = "auth" / "auth-int" / token
				if len(strings.TrimSpace(authorization.qop)) > 0 {
					result.WriteString(fmt.Sprintf(" qop=%s,", authorization.qop))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(nc))( )*=`).MatchString(orders) {
				// nonce-count = "nc" EQUAL nc-value,nc-value = 8LHEX
				if len(strings.TrimSpace(authorization.nc)) > 0 {
					result.WriteString(fmt.Sprintf(" nc=%s,", authorization.nc))
				}
				continue
			}
			ordersSlice := strings.Split(orders, "=")
			if len(ordersSlice) == 1 {
				if val, ok := authorization.authParam.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf("  %v=\"%v\",", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(" %v,", ordersSlice[0]))
					}
				} else {
					result.WriteString(fmt.Sprintf(" %v,", ordersSlice[0]))
				}
			} else {
				if val, ok := authorization.authParam.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf(" %v=\"%v\",", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(" %v,", ordersSlice[0]))
					}

				} else {
					result.WriteString(fmt.Sprintf(" %v=\"%v\",", ordersSlice[0], ordersSlice[1]))
				}
			}
		}
	} else {
		// username = "username" EQUAL username-value,username-value = quoted-string
		if len(strings.TrimSpace(authorization.username)) > 0 {
			result.WriteString(fmt.Sprintf(" username=\"%s\",", authorization.username))
		}
		// realm = "realm" EQUAL realm-value,realm-value = quoted-string
		if len(strings.TrimSpace(authorization.realm)) > 0 {
			result.WriteString(fmt.Sprintf(" realm=\"%s\",", authorization.realm))
		}
		// nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
		if len(strings.TrimSpace(authorization.nonce)) > 0 {
			result.WriteString(fmt.Sprintf(" nonce=\"%s\",", authorization.nonce))
		}
		// digest-uri = "uri" EQUAL LDQUOT digest-uri-value RDQUOT,digest-uri-value = rquest-uri ; Equal to request-uri as specified by HTTP/1.1
		if authorization.uri != nil {
			uri := authorization.uri.Raw()
			result.WriteString(fmt.Sprintf(" uri=\"%s\",", uri.String()))
		}
		// dresponse = "response" EQUAL request-digest, request-digest = LDQUOT 32LHEX RDQUOT
		if len(strings.TrimSpace(authorization.response)) > 0 {
			result.WriteString(fmt.Sprintf(" response=\"%s\",", authorization.response))
		}
		// algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
		if len(strings.TrimSpace(authorization.algorithm)) > 0 {
			result.WriteString(fmt.Sprintf(" algorithm=%s,", authorization.algorithm))
		}
		// cnonce = "cnonce" EQUAL cnonce-value,cnonce-value = nonce-value
		if len(strings.TrimSpace(authorization.cnonce)) > 0 {
			result.WriteString(fmt.Sprintf(" cnonce=\"%s\",", authorization.cnonce))
		}
		// opaque =  "opaque" EQUAL quoted-string
		if len(strings.TrimSpace(authorization.opaque)) > 0 {
			result.WriteString(fmt.Sprintf(" opaque=\"%s\",", authorization.opaque))
		}
		// message-qop = "qop" EQUAL qop-value,qop-value = "auth" / "auth-int" / token
		if len(strings.TrimSpace(authorization.qop)) > 0 {
			result.WriteString(fmt.Sprintf(" qop=%s,", authorization.qop))
		}
		// nonce-count = "nc" EQUAL nc-value,nc-value = 8LHEX
		if len(strings.TrimSpace(authorization.nc)) > 0 {
			result.WriteString(fmt.Sprintf(" nc=%s,", authorization.nc))
		}
	}

	// auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
	authorization.authParam.Range(func(key, value interface{}) bool {
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
func (authorization *Authorization) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`((?i)(?:^authorization))( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	authorization.source = raw
	authorization.uri = new(RequestUri)
	authorization.authParam = sync.Map{}

	field := fieldRegexp.FindString(raw)
	field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	authorization.field = field
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// auth-schema regexp
	authschemaRegexp := regexp.MustCompile(`(?i)(basic|digest)`)
	if authschemaRegexp.MatchString(raw) {
		authschema := authschemaRegexp.FindString(raw)
		authorization.authSchema = authschema
		raw = authschemaRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ",")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}

	// authorization order
	authorization.authorizationOrder(raw)

	// username regexp : username = "username" EQUAL username-value,username-value = quoted-string
	usernameRegexp := regexp.MustCompile(`((?i)(?:^username))( )*=`)
	// realm regexp : realm = "realm" EQUAL realm-value,realm-value = quoted-string
	realmRegexp := regexp.MustCompile(`((?i)(?:^realm))( )*=`)
	// nonce regexp : nonce = "nonce" EQUAL nonce-value,nonce-value = quoted-string
	nonceRegexp := regexp.MustCompile(`((?i)(?:^nonce))( )*=`)
	// uri regexp : digest-uri = "uri" EQUAL LDQUOT digest-uri-value RDQUOT,digest-uri-value = rquest-uri ; Equal to request-uri as specified by HTTP/1.1
	uriRegexp := regexp.MustCompile(`((?i)(?:^uri))( )*=`)
	// response regexp : dresponse = "response" EQUAL request-digest, request-digest = LDQUOT 32LHEX RDQUOT
	responseRegexp := regexp.MustCompile(`((?i)(?:^response))( )*=`)
	// algorithm regexp algorithm = "algorithm" EQUAL ( "MD5" / "MD5-sess"/ token )
	algorithmRegexp := regexp.MustCompile(`((?i)(?:^algorithm))( )*=`)
	// cnonce regexp : cnonce = "cnonce" EQUAL cnonce-value,cnonce-value = nonce-value
	cnonceRegexp := regexp.MustCompile(`((?i)(?:^cnonce))( )*=`)
	// opaque     regexp opaque =  "opaque" EQUAL quoted-string
	opaqueRegexp := regexp.MustCompile(`((?i)(?:^opaque))( )*=`)
	// qop        regexp message-qop = "qop" EQUAL qop-value,qop-value = "auth" / "auth-int" / token
	qopRegexp := regexp.MustCompile(`((?i)(?:^qop))( )*=`)
	// nc         regexp nonce-count = "nc" EQUAL nc-value,nc-value = 8LHEX
	ncRegexp := regexp.MustCompile(`((?i)(?:^nc))( )*=`)
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		raws = stringTrimPrefixAndTrimSuffix(raws, " ")
		switch {
		case usernameRegexp.MatchString(raws):
			username := usernameRegexp.ReplaceAllString(raws, "")
			username = regexp.MustCompile(`"`).ReplaceAllString(username, "")
			authorization.username = username
		case realmRegexp.MatchString(raws):
			realm := realmRegexp.ReplaceAllString(raws, "")
			realm = regexp.MustCompile(`"`).ReplaceAllString(realm, "")
			authorization.realm = realm
		case nonceRegexp.MatchString(raws):
			nonce := nonceRegexp.ReplaceAllString(raws, "")
			nonce = regexp.MustCompile(`"`).ReplaceAllString(nonce, "")
			authorization.nonce = nonce
		case uriRegexp.MatchString(raws):
			uri := uriRegexp.ReplaceAllString(raws, "")
			uri = regexp.MustCompile(`"`).ReplaceAllString(uri, "")
			authorization.uri.Parse(uri)
		case responseRegexp.MatchString(raws):
			response := responseRegexp.ReplaceAllString(raws, "")
			response = regexp.MustCompile(`"`).ReplaceAllString(response, "")
			authorization.response = response
		case algorithmRegexp.MatchString(raws):
			algorithm := algorithmRegexp.ReplaceAllString(raws, "")
			algorithm = regexp.MustCompile(`"`).ReplaceAllString(algorithm, "")
			authorization.algorithm = algorithm
		case cnonceRegexp.MatchString(raws):
			cnonce := cnonceRegexp.ReplaceAllString(raws, "")
			cnonce = regexp.MustCompile(`"`).ReplaceAllString(cnonce, "")
			authorization.cnonce = cnonce
		case opaqueRegexp.MatchString(raws):
			opaque := opaqueRegexp.ReplaceAllString(raws, "")
			opaque = regexp.MustCompile(`"`).ReplaceAllString(opaque, "")
			authorization.opaque = opaque
		case qopRegexp.MatchString(raws):
			qop := qopRegexp.ReplaceAllString(raws, "")
			qop = regexp.MustCompile(`"`).ReplaceAllString(qop, "")
			authorization.qop = qop
		case ncRegexp.MatchString(raws):
			nc := ncRegexp.ReplaceAllString(raws, "")
			nc = regexp.MustCompile(`"`).ReplaceAllString(nc, "")
			authorization.nc = nc
		default:
			// authParam  sync.Map    // auth-param = auth-param-name EQUAL ( token / quoted-string ),auth-param-name = token
			kvs := strings.Split(raws, "=")
			if len(kvs) == 1 {
				authorization.authParam.Store(kvs[0], "")
			} else {
				authorization.authParam.Store(kvs[0], kvs[1])
			}
		}
	}

}
func (authorization *Authorization) authorizationOrder(raw string) {
	authorization.isOrder = true
	authorization.order = make(chan string, 1024)
	defer close(authorization.order)
	raw = stringTrimPrefixAndTrimSuffix(raw, ",")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		authorization.order <- raws
	}
}
