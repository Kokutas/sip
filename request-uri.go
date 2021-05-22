package sip

// https://www.rfc-editor.org/rfc/rfc3261.html#section-8.1.1.1
//
// 8.1.1.1 Request-URI
// The initial Request-URI of the message SHOULD be set to the value of
// the URI in the To field.  One notable exception is the REGISTER
// method; behavior for setting the Request-URI of REGISTER is given in
// Section 10.  It may also be undesirable for privacy reasons or
// convenience to set these fields to the same value (especially if the
// originating UA expects that the Request-URI will be changed during
// transit).

// In some special circumstances, the presence of a pre-existing route
// set can affect the Request-URI of the message.  A pre-existing route
// set is an ordered set of URIs that identify a chain of servers, to
// which a UAC will send outgoing requests that are outside of a dialog.
// Commonly, they are configured on the UA by a user or service provider
// manually, or through some other non-SIP mechanism.  When a provider
// wishes to configure a UA with an outbound proxy, it is RECOMMENDED
// that this be done by providing it with a pre-existing route set with
// a single URI, that of the outbound proxy.

// When a pre-existing route set is present, the procedures for
// populating the Request-URI and Route header field detailed in Section
// 12.2.1.1 MUST be followed (even though there is no dialog), using the
// desired Request-URI as the remote target URI.

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
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
//
// SIP-URI          =  "sip:" [ userinfo ] hostport
//                     uri-parameters [ headers ]
// SIPS-URI         =  "sips:" [ userinfo ] hostport
//                     uri-parameters [ headers ]
// userinfo         =  ( user / telephone-subscriber ) [ ":" password ] "@"
// user             =  1*( unreserved / escaped / user-unreserved )
// user-unreserved  =  "&" / "=" / "+" / "$" / "," / ";" / "?" / "/"
// password         =  *( unreserved / escaped /
//                     "&" / "=" / "+" / "$" / "," )
// hostport         =  host [ ":" port ]
// host             =  hostname / IPv4address / IPv6reference
// hostname         =  *( domainlabel "." ) toplabel [ "." ]
// domainlabel      =  alphanum
//                     / alphanum *( alphanum / "-" ) alphanum
// toplabel         =  ALPHA / ALPHA *( alphanum / "-" ) alphanum
// IPv4address    =  1*3DIGIT "." 1*3DIGIT "." 1*3DIGIT "." 1*3DIGIT
// IPv6reference  =  "[" IPv6address "]"
// IPv6address    =  hexpart [ ":" IPv4address ]
// hexpart        =  hexseq / hexseq "::" [ hexseq ] / "::" [ hexseq ]
// hexseq         =  hex4 *( ":" hex4)
// hex4           =  1*4HEXDIG
// port           =  1*DIGIT

//    The BNF for telephone-subscriber can be found in RFC 2806 [9].  Note,
//    however, that any characters allowed there that are not allowed in
//    the user part of the SIP URI MUST be escaped.

// uri-parameters    =  *( ";" uri-parameter)
// uri-parameter     =  transport-param / user-param / method-param
//                      / ttl-param / maddr-param / lr-param / other-param
// transport-param   =  "transport="
//                      ( "udp" / "tcp" / "sctp" / "tls"
//                      / other-transport)
// other-transport   =  token
// user-param        =  "user=" ( "phone" / "ip" / other-user)
// other-user        =  token
// method-param      =  "method=" Method
// ttl-param         =  "ttl=" ttl
// maddr-param       =  "maddr=" host
// lr-param          =  "lr"
// other-param       =  pname [ "=" pvalue ]
// pname             =  1*paramchar
// pvalue            =  1*paramchar
// paramchar         =  param-unreserved / unreserved / escaped
// param-unreserved  =  "[" / "]" / "/" / ":" / "&" / "+" / "$"

// headers         =  "?" header *( "&" header )
// header          =  hname "=" hvalue
// hname           =  1*( hnv-unreserved / unreserved / escaped )
// hvalue          =  *( hnv-unreserved / unreserved / escaped )
// hnv-unreserved  =  "[" / "]" / "/" / "?" / ":" / "+" / "$"

// Request-URI    =  SIP-URI / SIPS-URI / absoluteURI
type RequestUri struct {
	// sipUri
	// sipsUri
	// aba
}
