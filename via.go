package sip

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-8.1.1.7
//
// 8.1.1.7 Via
// The Via header field indicates the transport used for the transaction
// and identifies the location where the response is to be sent.  A Via
// header field value is added only after the transport that will be
// used to reach the next hop has been selected (which may involve the
// usage of the procedures in [4]).

// When the UAC creates a request, it MUST insert a Via into that
// request.  The protocol name and protocol version in the header field
// MUST be SIP and 2.0, respectively.  The Via header field value MUST
// contain a branch parameter.  This parameter is used to identify the
// transaction created by that request.  This parameter is used by both
// the client and the server.

// The branch parameter value MUST be unique across space and time for
// all requests sent by the UA.  The exceptions to this rule are CANCEL
// and ACK for non-2xx responses.  As discussed below, a CANCEL request
// will have the same value of the branch parameter as the request it
// cancels.  As discussed in Section 17.1.1.3, an ACK for a non-2xx
// response will also have the same branch ID as the INVITE whose
// response it acknowledges.

// 	The uniqueness property of the branch ID parameter, to facilitate
// 	its use as a transaction ID, was not part of RFC 2543.

// The branch ID inserted by an element compliant with this
// specification MUST always begin with the characters "z9hG4bK".  These
// 7 characters are used as a magic cookie (7 is deemed sufficient to
// ensure that an older RFC 2543 implementation would not pick such a
// value), so that servers receiving the request can determine that the
// branch ID was constructed in the fashion described by this
// specification (that is, globally unique).  Beyond this requirement,
// the precise format of the branch token is implementation-defined.

// The Via header maddr, ttl, and sent-by components will be set when
// the request is processed by the transport layer (Section 18).

// Via processing for proxies is described in Section 16.6 Item 8 and
// Section 16.7 Item 3.

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.42
//
//20.42 Via
// The Via header field indicates the path taken by the request so far
// and indicates the path that should be followed in routing responses.
// The branch ID parameter in the Via header field values serves as a
// transaction identifier, and is used by proxies to detect loops.

// A Via header field value contains the transport protocol used to send
// the message, the client's host name or network address, and possibly
// the port number at which it wishes to receive responses.  A Via
// header field value can also contain parameters such as "maddr",
// "ttl", "received", and "branch", whose meaning and use are described
// in other sections.  For implementations compliant to this
// specification, the value of the branch parameter MUST start with the
// magic cookie "z9hG4bK", as discussed in Section 8.1.1.7.

// Transport protocols defined here are "UDP", "TCP", "TLS", and "SCTP".
// "TLS" means TLS over TCP.  When a request is sent to a SIPS URI, the
// protocol still indicates "SIP", and the transport protocol is TLS.

// Via: SIP/2.0/UDP erlang.bell-telephone.com:5060;branch=z9hG4bK87asdks7
// Via: SIP/2.0/UDP 192.0.2.1:5060 ;received=192.0.2.207
//   ;branch=z9hG4bK77asjd

// The compact form of the Via header field is v.

// In this example, the message originated from a multi-homed host with
// two addresses, 192.0.2.1 and 192.0.2.207.  The sender guessed wrong
// as to which network interface would be used.  Erlang.bell-
// telephone.com noticed the mismatch and added a parameter to the
// previous hop's Via header field value, containing the address that
// the packet actually came from.

// The host or network address and port number are not required to
// follow the SIP URI syntax.  Specifically, LWS on either side of the
// ":" or "/" is allowed, as shown here:

//    Via: SIP / 2.0 / UDP first.example.com: 4000;ttl=16
// 	 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1

// Even though this specification mandates that the branch parameter be
// present in all requests, the BNF for the header field indicates that
// it is optional.  This allows interoperation with RFC 2543 elements,
// which did not have to insert the branch parameter.

// Two Via header fields are equal if their sent-protocol and sent-by
// fields are equal, both have the same set of parameters, and the
// values of all parameters are equal.

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Via               =  ( "Via" / "v" ) HCOLON via-parm *(COMMA via-parm)
// via-parm          =  sent-protocol LWS sent-by *( SEMI via-params )
// via-params        =  via-ttl / via-maddr
//                      / via-received / via-branch
//                      / via-extension
// via-ttl           =  "ttl" EQUAL ttl
// via-maddr         =  "maddr" EQUAL host
// via-received      =  "received" EQUAL (IPv4address / IPv6address)
// via-branch        =  "branch" EQUAL token
// via-extension     =  generic-param
// sent-protocol     =  protocol-name SLASH protocol-version
//                      SLASH transport
// protocol-name     =  "SIP" / token
// protocol-version  =  token
// transport         =  "UDP" / "TCP" / "TLS" / "SCTP"
//                      / other-transport
// sent-by           =  host [ COLON port ]
// ttl               =  1*3DIGIT ; 0 to 255

// https://www.rfc-editor.org/rfc/rfc3581.html
//
// 3.  Client Behavior
// The client behavior specified here affects the transport processing
// defined in Section 18.1 of SIP (RFC 3261) [1].

// A client, compliant to this specification (clients include UACs and
// proxies), MAY include an "rport" parameter in the top Via header
// field value of requests it generates.  This parameter MUST have no
// value; it serves as a flag to indicate to the server that this
// extension is supported and requested for the transaction.

// When the client sends the request, if the request is sent using UDP,
// the client MUST be prepared to receive the response on the same IP
// address and port it used to populate the source IP address and source
// port of the request.  For backwards compatibility, the client MUST
// still be prepared to receive a response on the port indicated in the
// sent-by field of the topmost Via header field value, as specified in
// Section 18.1.1 of SIP [1].

// When there is a NAT between the client and server, the request will
// create (or refresh) a binding in the NAT.  This binding must remain
// in existence for the duration of the transaction in order for the
// client to receive the response.  Most UDP NAT bindings appear to have
// a timeout of about one minute.  This exceeds the duration of non-
// INVITE transactions.  Therefore, responses to a non-INVITE request
// will be received while the binding is still in existence.  INVITE
// transactions can take an arbitrarily long amount of time to complete.
// As a result, the binding may expire before a final response is
// received.  To keep the binding fresh, the client SHOULD retransmit
// its INVITE every 20 seconds or so.  These retransmissions will need
// to take place even after receiving a provisional response.
// A UA MAY execute the binding lifetime discovery algorithm in Section
// 10.2 of RFC 3489 [4] to determine the actual binding lifetime in the
// NAT.  If it is longer than 1 minute, the client SHOULD increase the
// interval for request retransmissions up to half of the discovered
// lifetime.  If it is shorter than one minute, it SHOULD decrease the
// interval for request retransmissions to half of the discovered
// lifetime.  Note that discovery of binding lifetimes can be
// unreliable.  See Section 14.3 of RFC 3489 [4].
//
// 4.  Server Behavior
// The server behavior specified here affects the transport processing
// defined in Section 18.2 of SIP [1].

// When a server compliant to this specification (which can be a proxy
// or UAS) receives a request, it examines the topmost Via header field
// value.  If this Via header field value contains an "rport" parameter
// with no value, it MUST set the value of the parameter to the source
// port of the request.  This is analogous to the way in which a server
// will insert the "received" parameter into the topmost Via header
// field value.  In fact, the server MUST insert a "received" parameter
// containing the source IP address that the request came from, even if
// it is identical to the value of the "sent-by" component.  Note that
// this processing takes place independent of the transport protocol.

// When a server attempts to send a response, it examines the topmost
// Via header field value of that response.  If the "sent-protocol"
// component indicates an unreliable unicast transport protocol, such as
// UDP, and there is no "maddr" parameter, but there is both a
// "received" parameter and an "rport" parameter, the response MUST be
// sent to the IP address listed in the "received" parameter, and the
// port in the "rport" parameter.  The response MUST be sent from the
// same address and port that the corresponding request was received on.
// This effectively adds a new processing step between bullets two and
// three in Section 18.2.2 of SIP [1].

// The response must be sent from the same address and port that the
// request was received on in order to traverse symmetric NATs.  When a
// server is listening for requests on multiple ports or interfaces, it
// will need to remember the one on which the request was received.  For
// a stateful proxy, storing this information for the duration of the
// transaction is not an issue.  However, a stateless proxy does not
// store state between a request and its response, and therefore cannot
// remember the address and port on which a request was received.  To
// properly implement this specification, a stateless proxy can encode
// the destination address and port of a request into the Via header
// field value that it inserts.  When the response arrives, it can
// extract this information and use it to forward the response.

// 5.  Syntax
// The syntax for the "rport" parameter is:

// response-port = "rport" [EQUAL 1*DIGIT]

// This extends the existing definition of the Via header field
// parameters, so that its BNF now looks like:

// via-params        =  via-ttl / via-maddr
// 					 / via-received / via-branch
// 					 / response-port / via-extension

// via-extension     =  generic-param
// generic-param  =  token [ EQUAL gen-value ]
// gen-value      =  token / host / quoted-string

type Via struct {
	field     string      // "Via" / "v"
	schema    string      // sip,sips,tel etc.
	version   float64     // 2.0
	transport string      // "UDP" / "TCP" / "TLS" / "SCTP"/ other-transport
	host      string      // host part,sent-by =  host [ COLON port ]
	port      uint16      // port part,sent-by =  host [ COLON port ]
	ttl       uint8       // via-ttl  =  "ttl" EQUAL ttl,ttl =  1*3DIGIT ; 0 to 255
	maddr     string      // via-maddr =  "maddr" EQUAL host
	received  net.IP      // via-received =  "received" EQUAL (IPv4address / IPv6address)
	branch    string      // via-branch =  "branch" EQUAL token
	rport     uint16      // response port -- RFC3581
	trans     string      // parameter transport,transport-param = "transport="( "udp" / "tcp" / "sctp" / "tls"/ other-transport),other-transport   =  token
	generic   sync.Map    // key:string,value:string, via-extension = generic-param,generic-param = token [ EQUAL gen-value ], gen-value = token / host / quoted-string
	isOrder   bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order     chan string // It is convenient to record the order of the original parameter fields when parsing
	source    string      // via header line source string
}

func (v *Via) SetField(field string) {
	if regexp.MustCompile(`^(?i)(via|v)$`).MatchString(field) {
		v.field = strings.Title(field)
	}
}
func (v *Via) GetField() string {
	return v.field
}
func (v *Via) SetSchema(schema string) {
	v.schema = schema
}
func (v *Via) GetSchema() string {
	return v.schema
}
func (v *Via) SetVersion(version float64) {
	v.version = version
}
func (v *Via) GetVersion() float64 {
	return v.version
}
func (v *Via) SetTransport(transport string) {
	v.transport = transport
}
func (v *Via) GetTransport() string {
	return v.transport
}
func (v *Via) SetHost(host string) {
	v.host = host
}
func (v *Via) GetHost() string {
	return v.host
}
func (v *Via) SetPort(port uint16) {
	v.port = port
}
func (v *Via) GetPort() uint16 {
	return v.port
}
func (v *Via) SetTtl(ttl uint8) {
	v.ttl = ttl
}
func (v *Via) GetTtl() uint8 {
	return v.ttl
}
func (v *Via) SetMaddr(maddr string) {
	v.maddr = maddr
}
func (v *Via) GetMaddr() string {
	return v.maddr
}
func (v *Via) SetReceived(received net.IP) {
	v.received = received
}
func (v *Via) GetReceived() net.IP {
	return v.received
}
func (v *Via) SetBranch(branch string) {
	v.branch = branch
}
func (v *Via) GetBranch() string {
	return v.branch
}
func (v *Via) SetRport(rport uint16) {
	v.rport = rport
}
func (v *Via) GetRport() uint16 {
	return v.rport
}
func (v *Via) SetTrans(transport string) {
	v.trans = transport
}
func (v *Via) GetTrans() string {
	return v.trans
}

func (v *Via) SetGeneric(generic sync.Map) {
	v.generic = generic
}
func (v *Via) GetGeneric() sync.Map {
	return v.generic
}
func (v *Via) GetOrder() []string {
	result := make([]string, 0)
	if reflect.DeepEqual(nil, v.order) {
		return result
	}
	for data := range v.order {
		result = append(result, data)
	}
	return result
}

func (v *Via) SetSource(source string) {
	v.source = source
}
func (v *Via) GetSource() string {
	return v.source
}
func NewVia(schema string, version float64, transport string, host string, port uint16, ttl uint8, maddr string, received net.IP, branch string, rport uint16, trans string, generic sync.Map) *Via {
	return &Via{
		schema:    schema,
		version:   version,
		transport: transport,
		host:      host,
		port:      port,
		ttl:       ttl,
		maddr:     maddr,
		received:  received,
		branch:    branch,
		rport:     rport,
		trans:     trans,
		generic:   generic,
		isOrder:   false,
		order:     make(chan string, 1024),
	}
}

func (v *Via) Raw() string {
	result := ""
	if len(strings.TrimSpace(v.field)) > 0 {
		result += fmt.Sprintf("%s:", v.field)
	} else {
		result += fmt.Sprintf("%s:", strings.Title("Via"))
	}
	if len(strings.TrimSpace(v.schema)) > 0 {
		result += fmt.Sprintf(" %s", strings.ToUpper(v.schema))
	}
	if v.version > 0 {
		result += fmt.Sprintf("/%1.1f", v.version)
	}
	if len(strings.TrimSpace(v.transport)) > 0 {
		result += fmt.Sprintf("/%s", strings.ToUpper(v.transport))
	}
	uri := ""
	if len(strings.TrimSpace(v.host)) > 0 {
		uri += fmt.Sprintf(" %s", v.host)
	}
	if v.port > 0 {
		if len(uri) > 0 {
			uri += fmt.Sprintf(":%d", v.port)
		} else {
			uri += fmt.Sprintf("%d", v.port)
		}
	}
	if v.isOrder {
		for data := range v.order {
			uri += fmt.Sprintf(";%s", data)
		}

		// if len(v.order) > 0 {
		// 	var orders []int
		// 	for index := range v.order {
		// 		orders = append(orders, index)
		// 	}
		// 	sort.Ints(orders)
		// 	for _, index := range orders {
		// 		switch strings.ToLower(v.order[index]) {
		// 		case "ttl":
		// 			if v.ttl > 0 {
		// 				uri += fmt.Sprintf(";ttl=%d", v.ttl)
		// 			}
		// 		case "maddr":
		// 			if len(strings.TrimSpace(v.maddr)) > 0 {
		// 				uri += fmt.Sprintf(";maddr=%s", v.maddr)
		// 			}
		// 		case "received":
		// 			if !reflect.DeepEqual(nil, v.received) {
		// 				uri += fmt.Sprintf(";received=%s", v.received.String())
		// 			}
		// 		case "branch":
		// 			if len(strings.TrimSpace(v.branch)) > 0 {
		// 				uri += fmt.Sprintf(";branch=%s", v.branch)
		// 			}
		// 		case "rport":
		// 			if v.rport == 1 {
		// 				uri += fmt.Sprintf(";%s", "rport")
		// 			} else if v.rport > 1 {
		// 				uri += fmt.Sprintf(";rport=%d", v.rport)
		// 			}
		// 		// case "generic":
		// 		default:
		// 			// if len(v.generic) > 0 {
		// 			// 	var ogs []int
		// 			// 	for i := range v.generic {
		// 			// 		ogs = append(ogs, i)
		// 			// 	}
		// 			// 	sort.Ints(ogs)
		// 			// 	for p, pv := range v.generic {
		// 			// 		uri += fmt.Sprintf(";%s", v.generic)
		// 			// 	}

		// 			// }
		// 		}
		// 	}

	} else {
		if v.rport == 1 {
			uri += fmt.Sprintf(";%s", "rport")
		} else if v.rport > 1 {
			uri += fmt.Sprintf(";rport=%d", v.rport)
		}
		if len(strings.TrimSpace(v.trans)) > 0 {
			uri += fmt.Sprintf(";transport=%s", v.trans)
		}
		if v.ttl > 0 {
			uri += fmt.Sprintf(";ttl=%d", v.ttl)
		}
		if len(strings.TrimSpace(v.maddr)) > 0 {
			uri += fmt.Sprintf(";maddr=%s", v.maddr)
		}
		if len(strings.TrimSpace(v.branch)) > 0 {
			uri += fmt.Sprintf(";branch=%s", v.branch)
		}
		if !reflect.DeepEqual(nil, v.received) {
			uri += fmt.Sprintf(";received=%s", v.received.String())
		}
		v.generic.Range(func(key, value interface{}) bool {
			switch value.(type) {
			// 这里需要判断——------------------------------------------------------------------------------------
			case interface{}:
				fmt.Println(value)
				uri += fmt.Sprintf(";%v=%v", key, value)
			default:
				uri += fmt.Sprintf(";%v", key)
			}

			// if len(strings.TrimSpace(fmt.Sprintf("%v", value))) == 0 {
			// 	uri += fmt.Sprintf(";%v", key)
			// } else {
			// 	uri += fmt.Sprintf(";%v=%v", key, value)
			// }
			return true
		})
		// if len(strings.TrimSpace(v.generic)) > 0 {
		// 	uri += fmt.Sprintf(";%s", v.generic)
		// }
	}
	v.isOrder = false
	if len(uri) > 0 {
		result += fmt.Sprintf(" %s", uri)
	}
	result += "\r\n"
	return result
}
func (v *Via) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// via field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(via|v)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	v.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	v.source = raw
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")

	// schema regexp
	schemasRegexpStr := `(?i)(`
	for _, v := range schemas {
		schemasRegexpStr += v + "|"
	}
	schemasRegexpStr = strings.TrimSuffix(schemasRegexpStr, "|")
	schemasRegexpStr += ")( )?"
	schemaRegexp := regexp.MustCompile(schemasRegexpStr)
	if schemaRegexp.MatchString(raw) {
		schema := schemaRegexp.FindString(raw)
		schema = stringTrimPrefixAndTrimSuffix(schema, " ")
		v.schema = schema
		raw = regexp.MustCompile(`.*`+schema).ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// version regexp
	versionRegexp := regexp.MustCompile(`/( )*\d+\.\d+`)
	if versionRegexp.MatchString(raw) {
		versions := versionRegexp.FindString(raw)
		versions = regexp.MustCompile(`.*/`).ReplaceAllString(versions, "")
		versions = stringTrimPrefixAndTrimSuffix(versions, " ")
		if len(versions) > 0 {
			version, _ := strconv.ParseFloat(versions, 64)
			v.version = version
			raw = regexp.MustCompile(`.*`+versions).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// transport regexp : "UDP" / "TCP" / "TLS" / "SCTP"/ other-transport
	transportRegexp := regexp.MustCompile(`/(?i)(udp|tcp|tls|sctp)`)
	if transportRegexp.MatchString(raw) {
		transport := transportRegexp.FindString(raw)
		transport = regexp.MustCompile(`.*/`).ReplaceAllString(transport, "")
		transport = stringTrimPrefixAndTrimSuffix(transport, " ")
		if len(transport) > 0 {
			v.transport = transport
			raw = regexp.MustCompile(`.*`+transport).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// host regexp
	hostRegexp := regexp.MustCompile(`(\w+\.\w+.*)|(\d+\.\d+\.\d+\.\d+)`)
	if hostRegexp.MatchString(raw) {
		host := hostRegexp.FindString(raw)
		host = stringTrimPrefixAndTrimSuffix(host, " ")
		if len(strings.TrimSpace(host)) > 0 {
			v.host = host
			raw = regexp.MustCompile(`.*`+host).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// port regexp
	portRegexp := regexp.MustCompile(`:( )*\d+`)
	if portRegexp.MatchString(raw) {
		ports := portRegexp.FindString(raw)
		ports = regexp.MustCompile(`\d+`).FindString(ports)
		if len(strings.TrimSpace(ports)) > 0 {
			port, _ := strconv.Atoi(ports)
			if port > 0 {
				v.port = uint16(port)
				raw = regexp.MustCompile(`.*`+ports).ReplaceAllString(raw, "")
			}
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// parameters
	// parameters order
	v.parametersOrder(raw)
	// ttl parameter regexp
	ttlRegexp := regexp.MustCompile(`(?i)(ttl).*`)
	// maddr parameter regexp
	maddrRegexp := regexp.MustCompile(`(?i)(maddr).*`)
	// received parameter regexp
	receivedRegexp := regexp.MustCompile(`(?i)(received).*`)
	// branch parameter regexp
	branchRegexp := regexp.MustCompile(`(?i)(branch).*`)
	// rport parameter regexp
	rportRegexp := regexp.MustCompile(`(?i)(rport).*`)
	// tranport parameter regexp
	transRegexp := regexp.MustCompile(`(?i)(transport).*`)
	rawSlice := strings.Split(raw, ";")
	for _, raws := range rawSlice {
		switch {
		case ttlRegexp.MatchString(raws):
			ttls := regexp.MustCompile(`(?i)(maddr)`).ReplaceAllString(raws, "")
			ttls = regexp.MustCompile(`.*=`).ReplaceAllString(ttls, "")
			ttls = stringTrimPrefixAndTrimSuffix(ttls, " ")
			if len(ttls) > 0 {
				ttl, _ := strconv.Atoi(ttls)
				if ttl > 0 {
					v.ttl = uint8(ttl)
				}
			}
		case maddrRegexp.MatchString(raws):

		case receivedRegexp.MatchString(raws):
		case branchRegexp.MatchString(raws):
		case rportRegexp.MatchString(raws):
		case transRegexp.MatchString(raws):
		default:
			if len(strings.TrimSpace(raws)) > 0 {
				if strings.Contains(raws, "=") {
					gs := strings.Split(raws, "=")
					if len(gs) > 1 {
						v.generic.Store(gs[0], gs[1])
					} else {
						v.generic.Store(gs[0], "")
					}

				} else {
					v.generic.Store(raws, "")
				}
			}
		}

	}

}

func (v *Via) parametersOrder(parameter string) {
	if reflect.DeepEqual(nil, v.order) {
		v.order = make(chan string, 1024)
	}
	// ttlRegexp := regexp.MustCompile(`(?i)(ttl)`)
	// maddrRegexp := regexp.MustCompile(`(?i)(maddr)`)
	// receivedRegexp := regexp.MustCompile(`(?i)(received)`)
	// branchRegexp := regexp.MustCompile(`(?i)(branch)`)
	// rportRegexp := regexp.MustCompile(`(?i)(rport)`)
	// transportRegexp := regexp.MustCompile(`(?i)(transport)`)
	v.isOrder = true
	parameters := strings.Split(parameter, ";")
	for _, data := range parameters {
		v.order <- data
		// switch {
		// case ttlRegexp.MatchString(p):
		// 	// v.order = append(v.order, "ttl")
		// 	v.order <- "ttl"
		// case maddrRegexp.MatchString(p):
		// 	// v.order = append(v.order, "maddr")
		// 	v.order <- "maddr"
		// case receivedRegexp.MatchString(p):
		// 	// v.order = append(v.order, "received")
		// 	v.order <- "received"
		// case branchRegexp.MatchString(p):
		// 	// v.order = append(v.order, "branch")
		// 	v.order <- "branch"
		// case rportRegexp.MatchString(p):
		// 	// v.order = append(v.order, "rport")
		// 	v.order <- "rport"
		// case transportRegexp.MatchString(p):
		// 	// v.order = append(v.order, "transport")
		// 	v.order <- "transport"
		// default:
		// 	// if len(strings.TrimSpace(p)) > 0 {
		// 	// 	v.order = append(v.order, "generic")
		// 	// }
		// }
	}
	close(v.order)
}
