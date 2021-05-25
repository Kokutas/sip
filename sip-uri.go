package sip

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
// SIP-URI          =  "sip:" [ userinfo ] hostport
// 						uri-parameters [ headers ]
// userinfo         =  ( user / telephone-subscriber ) [ ":" password ] "@"
// user             =  1*( unreserved / escaped / user-unreserved )
// user-unreserved  =  "&" / "=" / "+" / "$" / "," / ";" / "?" / "/"
// password         =  *( unreserved / escaped /
// 					"&" / "=" / "+" / "$" / "," )
// hostport         =  host [ ":" port ]
// host             =  hostname / IPv4address / IPv6reference
// hostname         =  *( domainlabel "." ) toplabel [ "." ]
// domainlabel      =  alphanum
// 					/ alphanum *( alphanum / "-" ) alphanum
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

// SIP-URI          =  "sip:" [ userinfo ] hostport
// 						uri-parameters [ headers ]
type SipUri struct {
	schema     string // sip/sips
	userinfo   *UserInfo
	hostport   *HostPort
	parameters *Parameters
	headers    sync.Map
	isOrder    bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order      chan string // It is convenient to record the order of the original parameter fields when parsing
	source     string      // sip-uri/sips-uri source string
}

func (su *SipUri) SetSchema(schema string) {
	if regexp.MustCompile(`^(?i)(sip|sips)$`).MatchString(schema) {
		su.schema = strings.ToLower(schema)
	} else {
		su.schema = "sip"
	}
}
func (su *SipUri) GetSchema() string {
	return su.schema
}
func (su *SipUri) SetUserInfo(userinfo *UserInfo) {
	su.userinfo = userinfo
}
func (su *SipUri) GetUserInfo() *UserInfo {
	return su.userinfo
}
func (su *SipUri) SetHostPort(hostport *HostPort) {
	su.hostport = hostport
}
func (su *SipUri) GetHostPort() *HostPort {
	return su.hostport
}
func (su *SipUri) SetParameters(parameters *Parameters) {
	su.parameters = parameters
}
func (su *SipUri) GetParameters() *Parameters {
	return su.parameters
}

func (su *SipUri) SetHeaders(headers sync.Map) {
	su.headers = headers
}
func (su *SipUri) GetHeaders() sync.Map {
	return su.headers
}

func (su *SipUri) GetSource() string {
	return su.source
}
func NewSipUri(userinfo *UserInfo, hostport *HostPort, parameters *Parameters, headers sync.Map) *SipUri {
	return &SipUri{
		schema:     "sip",
		userinfo:   userinfo,
		hostport:   hostport,
		parameters: parameters,
		headers:    headers,
		isOrder:    false,
	}
}
func (su *SipUri) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(su.schema)) == 0 {
		su.schema = "sip"
	}
	result.WriteString(fmt.Sprintf("%s:", strings.ToLower(su.schema)))
	if su.userinfo != nil {
		userinfo := su.userinfo.Raw()
		result.WriteString(userinfo.String())
	}
	if su.hostport != nil {
		hostport := su.hostport.Raw()
		result.WriteString(fmt.Sprintf("@%s", hostport.String()))
	}
	if su.parameters != nil {
		parameters := su.parameters.Raw()
		result.WriteString(parameters.String())
	}
	var headers strings.Builder
	if su.isOrder {
		for orders := range su.order {
			ordersSlice := strings.Split(orders, "=")
			if len(ordersSlice) == 1 {
				if val, ok := su.headers.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						headers.WriteString(fmt.Sprintf("&%v=%v", ordersSlice[0], val))
					} else {
						headers.WriteString(fmt.Sprintf("&%v", ordersSlice[0]))
					}
				} else {
					headers.WriteString(fmt.Sprintf("&%v", ordersSlice[0]))
				}
			} else {
				if val, ok := su.headers.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						headers.WriteString(fmt.Sprintf("&%v=%v", ordersSlice[0], val))
					} else {
						headers.WriteString(fmt.Sprintf("&%v", ordersSlice[0]))
					}
				} else {
					if len(strings.TrimSpace(fmt.Sprintf("%v", ordersSlice[1]))) > 0 {
						headers.WriteString(fmt.Sprintf("&%v=%v", ordersSlice[0], ordersSlice[1]))
					} else {
						headers.WriteString(fmt.Sprintf("&%v", ordersSlice[0]))
					}
				}
			}
		}
	}

	su.headers.Range(func(key, value interface{}) bool {
		if reflect.ValueOf(value).IsValid() {
			if reflect.ValueOf(value).IsZero() {
				headers.WriteString(fmt.Sprintf("&%v", key))
				return true
			}
			headers.WriteString(fmt.Sprintf("&%v=%v", key, value))
			return true
		}
		headers.WriteString(fmt.Sprintf("&%v", key))
		return true
	})
	if len(headers.String()) > 0 {
		result.WriteString(fmt.Sprintf("?%s", headers.String()))
	}
	return
}
func (su *SipUri) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// schema regexp
	schemaRegexp := regexp.MustCompile(`^((?i)(sip|sips)( )?:)`)
	if !schemaRegexp.MatchString(raw) {
		return
	}
	schema := schemaRegexp.FindString(raw)
	schema = regexp.MustCompile(`:`).ReplaceAllString(schema, "")
	schema = stringTrimPrefixAndTrimSuffix(schema, " ")
	su.schema = schema
	su.source = raw

	su.userinfo = new(UserInfo)
	su.hostport = new(HostPort)
	su.parameters = new(Parameters)
	su.headers = sync.Map{}

	raw = schemaRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// headers regexp
	headersRegexp := regexp.MustCompile(`\?.*`)
	if headersRegexp.MatchString(raw) {

		headers := headersRegexp.FindString(raw)
		raw = headersRegexp.ReplaceAllString(raw, "")
		headers = regexp.MustCompile(`\?`).ReplaceAllString(headers, "")
		headers = stringTrimPrefixAndTrimSuffix(headers, "&")
		headers = stringTrimPrefixAndTrimSuffix(headers, " ")
		// sip-uri/sips-uri header order
		su.headersOrder(headers)
		if len(strings.TrimSpace(headers)) > 0 {
			headersSlice := strings.Split(headers, "&")
			if len(headersSlice) == 1 {
				kvs := strings.Split(headersSlice[0], "=")
				if len(kvs) == 1 {
					su.headers.Store(kvs[0], "")
				} else {
					su.headers.Store(kvs[0], kvs[1])
				}
			} else {
				for _, hs := range headersSlice {
					kvs := strings.Split(hs, "=")
					if len(kvs) == 1 {
						su.headers.Store(kvs[0], "")
					} else {
						su.headers.Store(kvs[0], kvs[1])
					}
				}
			}
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// uri-parameters regexp
	uriparametersRegexp := regexp.MustCompile(`;.*`)
	if uriparametersRegexp.MatchString(raw) {
		parameters := uriparametersRegexp.FindString(raw)
		parameters = stringTrimPrefixAndTrimSuffix(parameters, ";")
		su.parameters.Parse(parameters)
		raw = uriparametersRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// host port regexp
	hostportRegexp := regexp.MustCompile(`@.*`)
	if hostportRegexp.MatchString(raw) {
		hostport := hostportRegexp.FindString(raw)
		hostport = regexp.MustCompile(`@`).ReplaceAllString(hostport, "")
		hostport = stringTrimPrefixAndTrimSuffix(hostport, " ")
		su.hostport.Parse(hostport)
		raw = hostportRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		su.userinfo.Parse(raw)
	}
}
func (su *SipUri) headersOrder(raw string) {
	su.isOrder = true
	su.order = make(chan string, 1024)
	defer close(su.order)
	raw = stringTrimPrefixAndTrimSuffix(raw, "&")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	rawSlice := strings.Split(raw, "&")
	for _, raws := range rawSlice {
		su.order <- raws
	}
}
