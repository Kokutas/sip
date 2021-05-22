package sip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-7.1
//
// 7.1 Requests
// SIP requests are distinguished by having a Request-Line for a start-
// line.  A Request-Line contains a method name, a Request-URI, and the
// protocol version separated by a single space (SP) character.

// The Request-Line ends with CRLF.  No CR or LF are allowed except in
// the end-of-line CRLF sequence.  No linear whitespace (LWS) is allowed
// in any of the elements.

// 		Request-Line  =  Method SP Request-URI SP SIP-Version CRLF

// 	Method: This specification defines six methods: REGISTER for
// 		registering contact information, INVITE, ACK, and CANCEL for
// 		setting up sessions, BYE for terminating sessions, and
// 		OPTIONS for querying servers about their capabilities.  SIP
// 		extensions, documented in standards track RFCs, may define
// 		additional methods.
// 		Request-URI: The Request-URI is a SIP or SIPS URI as described in
// 		Section 19.1 or a general URI (RFC 2396 [5]).  It indicates
// 		the user or service to which this request is being addressed.
// 		The Request-URI MUST NOT contain unescaped spaces or control
// 		characters and MUST NOT be enclosed in "<>".

// 		SIP elements MAY support Request-URIs with schemes other than
// 		"sip" and "sips", for example the "tel" URI scheme of RFC
// 		2806 [9].  SIP elements MAY translate non-SIP URIs using any
// 		mechanism at their disposal, resulting in SIP URI, SIPS URI,
// 		or some other scheme.

// 	SIP-Version: Both request and response messages include the
// 		version of SIP in use, and follow [H3.1] (with HTTP replaced
// 		by SIP, and HTTP/1.1 replaced by SIP/2.0) regarding version
// 		ordering, compliance requirements, and upgrading of version
// 		numbers.  To be compliant with this specification,
// 		applications sending SIP messages MUST include a SIP-Version
// 		of "SIP/2.0".  The SIP-Version string is case-insensitive,
// 		but implementations MUST send upper-case.

// 		Unlike HTTP/1.1, SIP treats the version number as a literal
// 		string.  In practice, this should make no difference.

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Request-Line   =  Method SP Request-URI SP SIP-Version CRLF
// Request-URI    =  SIP-URI / SIPS-URI / absoluteURI
// absoluteURI    =  scheme ":" ( hier-part / opaque-part )
// hier-part      =  ( net-path / abs-path ) [ "?" query ]
// net-path       =  "//" authority [ abs-path ]
// abs-path       =  "/" path-segments
// opaque-part    =  uric-no-slash *uric
// uric           =  reserved / unreserved / escaped
// uric-no-slash  =  unreserved / escaped / ";" / "?" / ":" / "@"
//                   / "&" / "=" / "+" / "$" / ","
// path-segments  =  segment *( "/" segment )
// segment        =  *pchar *( ";" param )
// param          =  *pchar
// pchar          =  unreserved / escaped /
//                   ":" / "@" / "&" / "=" / "+" / "$" / ","
// scheme         =  ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
// authority      =  srvr / reg-name
// srvr           =  [ [ userinfo "@" ] hostport ]
// reg-name       =  1*( unreserved / escaped / "$" / ","
//                   / ";" / ":" / "@" / "&" / "=" / "+" )
// query          =  *uric
// SIP-Version    =  "SIP" "/" 1*DIGIT "." 1*DIGIT

// https://www.rfc-editor.org/rfc/rfc3261.html#section-27.4
//
// 27.4 Method and Response Codes
//    This specification establishes the Method and Response-Code sub-
//    registries under http://www.iana.org/assignments/sip-parameters and
//    initiates their population as follows.  The initial Methods table is:
//    INVITE                   [RFC3261]
//    ACK                      [RFC3261]
//    BYE                      [RFC3261]
//    CANCEL                   [RFC3261]
//    REGISTER                 [RFC3261]
//    OPTIONS                  [RFC3261]
//    INFO                     [RFC2976]

// The response code table is initially populated from Section 21, the
// portions labeled Informational, Success, Redirection, Client-Error,
// Server-Error, and Global-Failure.  The table has the following
// format:

// Type (e.g., Informational)
// 	  Number    Default Reason Phrase         [RFC3261]

// The following information needs to be provided in an RFC publication
// in order to register a new response code or method:

// o  The RFC number in which the method or response code is
//    registered;

// o  the number of the response code or name of the method being
//    registered;

// o  the default reason phrase for that response code, if
//    applicable;
type RequestLine struct {
	method     string      // method:INVITE, ACK,BYE,CANCEL,REGISTER,OPTIONS,INFO etc.
	requestUri *RequestUri // SIP-URI/SIPS-URI
	schema     string      // sip,sips,tel etc.
	version    float64     // 2.0
	isOrder    bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order      chan string // It is convenient to record the order of the original parameter fields when parsing
	source     string      // userinfo source string
}

func (requestLine *RequestLine) SetMethod(method string) {
	requestLine.method = method
}
func (requestLine *RequestLine) GetMethod() string {
	return requestLine.method
}
func (requestLine *RequestLine) SetRequestUri(requestUri *RequestUri) {
	requestLine.requestUri = requestUri
}
func (requestLine *RequestLine) GetRequestUri() *RequestUri {
	return requestLine.requestUri
}

func (requestLine *RequestLine) SetSchema(schema string) {
	requestLine.schema = schema
}
func (requestLine *RequestLine) GetSchema() string {
	return requestLine.schema
}
func (requestLine *RequestLine) SetVersion(version float64) {
	requestLine.version = version
}
func (requestLine *RequestLine) GetVersion() float64 {
	return requestLine.version
}
func (requestLine *RequestLine) GetIsOrder() bool {
	return requestLine.isOrder
}
func (requestLine *RequestLine) GetOrder() []string {
	result := make([]string, 0)
	if requestLine.order == nil {
		return result
	}
	for data := range requestLine.order {
		result = append(result, data)
	}
	return result
}
func (requestLine *RequestLine) SetSource(source string) {
	requestLine.source = source
}
func (requestLine *RequestLine) GetSource() string {
	return requestLine.source
}

func NewRequestLine(method string, requestUri *RequestUri, schema string, version float64) *RequestLine {
	return &RequestLine{
		method:     method,
		requestUri: requestUri,
		schema:     schema,
		version:    version,
		isOrder:    false,
		order:      make(chan string, 1024),
	}
}
func (requestLine *RequestLine) Raw() string {
	result := ""
	if requestLine.isOrder {
		for data := range requestLine.order {
			result += data
		}
		requestLine.isOrder = false
		result += "\r\n"
		return result
	}

	// method:INVITE, ACK,BYE,CANCEL,REGISTER,OPTIONS,INFO etc.
	if len(strings.TrimSpace(requestLine.method)) > 0 {
		result += strings.ToUpper(requestLine.method)
	}
	// SIP-URI/SIPS-URI
	if requestLine.requestUri != nil {
		result += fmt.Sprintf(" %s", requestLine.requestUri.Raw())
	}
	// schema: sip,sips,tel etc.
	if len(strings.TrimSpace(requestLine.schema)) > 0 {
		result += fmt.Sprintf(" %s", strings.ToUpper(requestLine.schema))
	}
	// version: 2.0
	result += fmt.Sprintf("/%1.1f", requestLine.version)
	result += "\r\n"
	return result
}
func (requestLine *RequestLine) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// method regexp string
	methodsRegexpStr := `^(?i)(`
	for _, v := range methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")( )?"
	// method regexp
	methodRegexp := regexp.MustCompile(methodsRegexpStr)

	// schema regexp string
	schemasRegexpStr := `( )?(?i)(`
	for _, v := range schemas {
		schemasRegexpStr += v + "|"
	}
	schemasRegexpStr = strings.TrimSuffix(schemasRegexpStr, "|")
	schemasRegexpStr += ")( )?"
	// schema and version regexp
	schemaAndVersionRegexp := regexp.MustCompile(schemasRegexpStr + `/( )?\d\.\d`)

	if !methodRegexp.MatchString(raw) && !schemaAndVersionRegexp.MatchString(raw) {
		return
	}
	requestLine.source = raw
	requestLine.requestUri = new(RequestUri)
	// request-line order
	requestLine.requestlineOrder(raw)
	if methodRegexp.MatchString(raw) {
		method := methodRegexp.FindString(raw)
		method = stringTrimPrefixAndTrimSuffix(method, " ")
		requestLine.method = method
		raw = strings.TrimPrefix(raw, method)
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// schema regexp
	schemaRegexp := regexp.MustCompile(schemasRegexpStr)
	// version regexp
	versionRegexp := regexp.MustCompile(`\d\.[0-9]{1}`)
	if schemaAndVersionRegexp.MatchString(raw) {
		schemaAndVersion := schemaAndVersionRegexp.FindString(raw)
		schemaAndVersion = stringTrimPrefixAndTrimSuffix(schemaAndVersion, " ")
		if schemaRegexp.MatchString(schemaAndVersion) {
			schema := schemaRegexp.FindString(schemaAndVersion)
			schema = stringTrimPrefixAndTrimSuffix(schema, " ")
			requestLine.schema = schema
		}
		if versionRegexp.MatchString(schemaAndVersion) {
			versions := versionRegexp.Find([]byte(schemaAndVersion))
			version, _ := strconv.ParseFloat(string(versions), 64)
			requestLine.version = version
		}
		raw = strings.ReplaceAll(raw, schemaAndVersion, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {

		requestLine.requestUri.Parse(raw)
	}
}
func (requestLine *RequestLine) requestlineOrder(raw string) {
	if requestLine.order == nil {
		requestLine.order = make(chan string, 1024)
	}
	requestLine.isOrder = true
	defer close(requestLine.order)
	requestLine.order <- raw
}
