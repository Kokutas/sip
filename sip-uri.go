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
	schema        string // sip/sips
	userinfo      *UserInfo
	hostport      *HostPort
	uriparameters *UriParameters
	headers       sync.Map
	isOrder       bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order         chan string // It is convenient to record the order of the original parameter fields when parsing
	source        string      // userinfo source string
}

func (sipUri *SipUri) SetSchema(schema string) {
	if regexp.MustCompile(`^(?i)(sip|sips)$`).MatchString(schema) {
		sipUri.schema = strings.ToLower(schema)
	} else {
		sipUri.schema = "sip"
	}
}
func (sipUri *SipUri) GetSchema() string {
	return sipUri.schema
}
func (sipUri *SipUri) SetUserInfo(userinfo *UserInfo) {
	sipUri.userinfo = userinfo
}
func (sipUri *SipUri) GetUserInfo() *UserInfo {
	return sipUri.userinfo
}
func (sipUri *SipUri) SetHostPort(hostport *HostPort) {
	sipUri.hostport = hostport
}
func (sipUri *SipUri) GetHostPort() *HostPort {
	return sipUri.hostport
}
func (sipUri *SipUri) SetUriParameters(uriParameters *UriParameters) {
	sipUri.uriparameters = uriParameters
}
func (sipUri *SipUri) GetUriParameters() *UriParameters {
	return sipUri.uriparameters
}

func (sipUri *SipUri) SetHeaders(headers sync.Map) {
	sipUri.headers = headers
}
func (sipUri *SipUri) GetHeaders() sync.Map {
	return sipUri.headers
}
func (sipUri *SipUri) GetIsOrder() bool {
	return sipUri.isOrder
}
func (sipUri *SipUri) GetOrder() []string {
	result := make([]string, 0)
	if sipUri.order == nil {
		return result
	}
	for data := range sipUri.order {
		result = append(result, data)
	}
	return result
}
func (sipUri *SipUri) SetSource(source string) {
	sipUri.source = source
}
func (sipUri *SipUri) GetSource() string {
	return sipUri.source
}
func NewSipUri(userinfo *UserInfo, hostport *HostPort, uriparameters *UriParameters, headers sync.Map) *SipUri {
	return &SipUri{
		schema:        "sip",
		userinfo:      userinfo,
		hostport:      hostport,
		uriparameters: uriparameters,
		headers:       headers,
		isOrder:       false,
		order:         make(chan string, 1024),
	}
}
func (sipUri *SipUri) Raw() string {
	result := ""
	if sipUri.isOrder {
		for data := range sipUri.order {
			result += data
		}
	} else {
		if len(strings.TrimSpace(sipUri.schema)) == 0 {
			sipUri.schema = "sip"
		}
		result += fmt.Sprintf("%s:", strings.ToLower(sipUri.schema))
		if sipUri.userinfo != nil {
			result += sipUri.userinfo.Raw()
		}
		if sipUri.hostport != nil {
			result += fmt.Sprintf("@%s", sipUri.hostport.Raw())
		}
		if sipUri.uriparameters != nil {
			result += sipUri.uriparameters.Raw()
		}
		headers := ""
		sipUri.headers.Range(func(key, value interface{}) bool {
			if reflect.ValueOf(value).IsValid() {
				if reflect.ValueOf(value).IsZero() {
					headers += fmt.Sprintf("&%v", key)
					return true
				}
				headers += fmt.Sprintf("&%v=%v", key, value)
				return true
			}
			headers += fmt.Sprintf("&%v", key)
			return true
		})
		headers = stringTrimPrefixAndTrimSuffix(headers, "&")
		if len(headers) > 0 {

			result += fmt.Sprintf("?%s", headers)
		}
	}
	return result
}
func (sipUri *SipUri) Parse(raw string) {
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
	sipUri.schema = schema
	sipUri.source = raw

	sipUri.userinfo = new(UserInfo)
	sipUri.hostport = new(HostPort)
	sipUri.uriparameters = new(UriParameters)
	sipUri.headers = sync.Map{}

	// sip-uri/sips-uri order
	sipUri.sipUriOrder(raw)
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
		if len(strings.TrimSpace(headers)) > 0 {
			headersSlice := strings.Split(headers, "&")
			if len(headersSlice) == 1 {
				kvs := strings.Split(headersSlice[0], "=")
				if len(kvs) == 1 {
					sipUri.headers.Store(kvs[0], "")
				} else {
					sipUri.headers.Store(kvs[0], kvs[1])
				}
			} else {
				for _, hs := range headersSlice {
					kvs := strings.Split(hs, "=")
					if len(kvs) == 1 {
						sipUri.headers.Store(kvs[0], "")
					} else {
						sipUri.headers.Store(kvs[0], kvs[1])
					}
				}
			}
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// uri-parameters regexp
	uriparametersRegexp := regexp.MustCompile(`;.*`)
	if uriparametersRegexp.MatchString(raw) {
		uriparameters := uriparametersRegexp.FindString(raw)
		uriparameters = stringTrimPrefixAndTrimSuffix(uriparameters, ";")
		sipUri.uriparameters.Parse(uriparameters)
		raw = uriparametersRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// host port regexp
	hostportRegexp := regexp.MustCompile(`@.*`)
	if hostportRegexp.MatchString(raw) {
		hostport := hostportRegexp.FindString(raw)
		hostport = regexp.MustCompile(`@`).ReplaceAllString(hostport, "")
		hostport = stringTrimPrefixAndTrimSuffix(hostport, " ")
		sipUri.hostport.Parse(hostport)
		raw = hostportRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		sipUri.userinfo.Parse(raw)
	}
}
func (sipUri *SipUri) sipUriOrder(raw string) {
	if sipUri.order == nil {
		sipUri.order = make(chan string, 1024)
	}
	sipUri.isOrder = true
	defer close(sipUri.order)
	sipUri.order <- raw
}

// userinfo         =  ( user / telephone-subscriber ) [ ":" password ] "@"
// user             =  1*( unreserved / escaped / user-unreserved )
// user-unreserved  =  "&" / "=" / "+" / "$" / "," / ";" / "?" / "/"
// password         =  *( unreserved / escaped /
// 					"&" / "=" / "+" / "$" / "," )
//
// https://www.rfc-editor.org/rfc/rfc2806
//
// telephone-subscriber  = global-phone-number / local-phone-number
// global-phone-number   = "+" base-phone-number [isdn-subaddress]
//                         [post-dial] *(area-specifier /
//                         service-provider / future-extension)
// base-phone-number     = 1*phonedigit
// local-phone-number    = 1*(phonedigit / dtmf-digit /
//                         pause-character) [isdn-subaddress]
//                         [post-dial] area-specifier
//                         *(area-specifier / service-provider /
//                         future-extension)
// isdn-subaddress       = ";isub=" 1*phonedigit
// post-dial             = ";postd=" 1*(phonedigit /
//                         dtmf-digit / pause-character)
// area-specifier        = ";" phone-context-tag "=" phone-context-ident
// phone-context-tag     = "phone-context"
// phone-context-ident   = network-prefix / private-prefix
// network-prefix        = global-network-prefix / local-network-prefix
// global-network-prefix = "+" 1*phonedigit
// local-network-prefix  = 1*(phonedigit / dtmf-digit / pause-character)
// private-prefix        = (%x21-22 / %x24-27 / %x2C / %x2F / %x3A /
//                         %x3C-40 / %x45-4F / %x51-56 / %x58-60 /
//                         %x65-6F / %x71-76 / %x78-7E)
//                         *(%x21-3A / %x3C-7E)
//                         ; Characters in URLs must follow escaping rules
//                         ; as explained in [RFC2396]
type UserInfo struct {
	user                string
	telephoneSubscriber string
	password            string
	isOrder             bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order               chan string // It is convenient to record the order of the original parameter fields when parsing
	source              string      // userinfo source string
}

func (userInfo *UserInfo) SetUser(user string) {
	userInfo.user = user
}
func (userInfo *UserInfo) GetUser() string {
	return userInfo.user
}
func (userInfo *UserInfo) SetTelephoneSubscriber(telephoneSubscriber string) {
	userInfo.telephoneSubscriber = telephoneSubscriber
}
func (userInfo *UserInfo) GetTelephoneSubscriber() string {
	return userInfo.telephoneSubscriber
}
func (userInfo *UserInfo) SetPassword(password string) {
	userInfo.password = password
}
func (userInfo *UserInfo) GetPassword() string {
	return userInfo.password
}
func (userInfo *UserInfo) GetIsOrder() bool {
	return userInfo.isOrder
}
func (userInfo *UserInfo) GetOrder() []string {
	result := make([]string, 0)
	if userInfo.order == nil {
		return result
	}
	for data := range userInfo.order {
		result = append(result, data)
	}
	return result
}
func (userInfo *UserInfo) SetSource(source string) {
	userInfo.source = source
}
func (userInfo *UserInfo) GetSource() string {
	return userInfo.source
}

func NewUserInfo(user string, telephoneSubscriber string, password string) *UserInfo {
	return &UserInfo{
		user:                user,
		telephoneSubscriber: telephoneSubscriber,
		password:            password,
		isOrder:             false,
		order:               make(chan string, 1024),
	}
}
func (userInfo *UserInfo) Raw() string {
	result := ""
	if userInfo.isOrder {
		for data := range userInfo.order {
			result += data
		}
	} else {
		switch {
		case len(strings.TrimSpace(userInfo.user)) > 0:
			result += userInfo.user
		case len(strings.TrimSpace(userInfo.telephoneSubscriber)) > 0:
			result += userInfo.telephoneSubscriber
		}
		if len(strings.TrimSpace(userInfo.password)) > 0 {
			result += fmt.Sprintf(":%s", userInfo.password)
		}
	}
	return result
}
func (userInfo *UserInfo) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	userInfo.source = raw
	// userinfo order
	userInfo.userinfoOrder(raw)
	// password regexp
	passwordRegexp := regexp.MustCompile(`:.*(@)*?(;)*?(\?)*?`)
	if passwordRegexp.MatchString(raw) {
		password := regexp.MustCompile(`:`).ReplaceAllString(passwordRegexp.FindString(raw), "")
		password = regexp.MustCompile(`@.*`).ReplaceAllString(password, "")
		password = regexp.MustCompile(`;.*`).ReplaceAllString(password, "")
		password = regexp.MustCompile(`\?.*`).ReplaceAllString(password, "")
		userInfo.password = password
		raw = passwordRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// telephone-subscriber regexp
	// 1.global-phone-number
	// 2.local-phone-number
	telephoneSubscribeRegexp := regexp.MustCompile(`(^(\+)?(\d{1,3}\-)?\d+\-\d+(\-\d+)?)$|(^\+\d+)|(^\d{11}$)`)
	if len(strings.TrimSpace(raw)) > 0 {
		if telephoneSubscribeRegexp.MatchString(raw) {
			userInfo.telephoneSubscriber = telephoneSubscribeRegexp.FindString(raw)
		} else {
			userInfo.user = raw
		}
	}
}
func (userInfo *UserInfo) userinfoOrder(raw string) {
	if userInfo.order == nil {
		userInfo.order = make(chan string, 1024)
	}
	userInfo.isOrder = true
	defer close(userInfo.order)
	userInfo.order <- raw
}

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
type HostPort struct {
	hostname      string
	ipv4Address   net.IP
	ipv6Reference net.IP
	port          uint16
	isOrder       bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order         chan string // It is convenient to record the order of the original parameter fields when parsing
	source        string      // hostport source string
}

func (hostport *HostPort) SetHostname(hostname string) {
	hostport.hostname = hostname
}
func (hostport *HostPort) GetHostname() string {
	return hostport.hostname
}

func (hostport *HostPort) SetIPv4Address(ipv4Address net.IP) {
	hostport.ipv4Address = ipv4Address
}
func (hostport *HostPort) GetIPv4Address() net.IP {
	return hostport.ipv4Address
}
func (hostport *HostPort) SetIPv6Reference(ipv6Reference net.IP) {
	hostport.ipv6Reference = ipv6Reference
}
func (hostport *HostPort) GetIPv6Reference() net.IP {
	return hostport.ipv6Reference
}
func (hostport *HostPort) SetPort(port uint16) {
	hostport.port = port
}
func (hostport *HostPort) GetPort() uint16 {
	return hostport.port
}
func (hostport *HostPort) GetIsOrder() bool {
	return hostport.isOrder
}
func (hostport *HostPort) GetOrder() []string {
	result := make([]string, 0)
	if hostport.order == nil {
		return result
	}
	for data := range hostport.order {
		result = append(result, data)
	}
	return result
}
func (hostport *HostPort) SetSource(source string) {
	hostport.source = source
}
func (hostport *HostPort) GetSource() string {
	return hostport.source
}
func NewHostPort(hostname string, ipv4Address net.IP, ipv6Reference net.IP, port uint16) *HostPort {
	return &HostPort{
		hostname:      hostname,
		ipv4Address:   ipv4Address,
		ipv6Reference: ipv6Reference,
		port:          port,
		isOrder:       false,
		order:         make(chan string, 1024),
	}
}

func (hostport *HostPort) Raw() string {
	result := ""
	if hostport.isOrder {
		for data := range hostport.order {
			result += data
		}
	} else {
		switch {
		case len(strings.TrimSpace(hostport.hostname)) > 0:
			result += hostport.hostname
		case hostport.ipv4Address != nil:
			result += hostport.ipv4Address.String()
		case hostport.ipv6Reference != nil:
			result += fmt.Sprintf("["+"%s"+"]", hostport.ipv6Reference.String())
		}
		if hostport.port > 0 {
			result += fmt.Sprintf(":%d", hostport.port)
		}
	}
	return result
}
func (hostport *HostPort) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	hostport.source = raw
	// host port order
	hostport.hostportOrder(raw)
	// ipv4 address regexp
	ipv4AddressRegexp := regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
	// host name regexp
	hostnameRegexp := regexp.MustCompile(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`)
	switch {
	case ipv4AddressRegexp.MatchString(raw):
		hostport.ipv4Address = net.ParseIP(ipv4AddressRegexp.FindString(raw))
		raw = regexp.MustCompile(`.*`+hostport.ipv4Address.String()).ReplaceAllString(raw, "")
	case hostnameRegexp.MatchString(raw):
		hostport.hostname = hostnameRegexp.FindString(raw)
		raw = regexp.MustCompile(`.*`+hostport.hostname).ReplaceAllString(raw, "")
	default:
		ipAddr, err := net.ResolveIPAddr("ip", raw)
		if err == nil {
			if ipAddr.IP.To16() != nil {
				hostport.ipv6Reference = ipAddr.IP.To16()
				raw = regexp.MustCompile(`.*`+hostport.ipv6Reference.String()).ReplaceAllString(raw, "")
			} else if ipAddr.IP.To4() != nil {
				hostport.ipv4Address = ipAddr.IP.To4()
				raw = regexp.MustCompile(`.*`+hostport.ipv4Address.String()).ReplaceAllString(raw, "")
			}
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
				hostport.port = uint16(port)
			}
		}
	}
}
func (hostport *HostPort) hostportOrder(raw string) {
	if hostport.order == nil {
		hostport.order = make(chan string, 1024)
	}
	hostport.isOrder = true
	defer close(hostport.order)
	hostport.order <- raw
}

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
//
type UriParameters struct {
	transport string      // transport-param = "transport="( "udp" / "tcp" / "sctp" / "tls"/ other-transport),other-transport = token
	user      string      // user-param =  "user=" ( "phone" / "ip" / other-user), other-user = token
	method    string      // method-param =  "method=" Method
	ttl       uint8       // ttl-param =  "ttl=" ttl
	maddr     string      // maddr-param       =  "maddr=" host
	lr        bool        // lr-param          =  "lr"
	other     sync.Map    // other-param       =  pname [ "=" pvalue ]
	isOrder   bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order     chan string // It is convenient to record the order of the original parameter fields when parsing
	source    string      // uri-parameters source string
}

func (uriParameters *UriParameters) SetTransport(transport string) {
	uriParameters.transport = transport
}
func (uriParameters *UriParameters) GetTransport() string {
	return uriParameters.transport
}
func (uriParameters *UriParameters) SetUser(user string) {
	uriParameters.user = user
}
func (uriParameters *UriParameters) GetUser() string {
	return uriParameters.user
}
func (uriParameters *UriParameters) SetMethod(method string) {
	uriParameters.method = method
}
func (uriParameters *UriParameters) GetMethod() string {
	return uriParameters.method
}
func (uriParameters *UriParameters) SetTtl(ttl uint8) {
	uriParameters.ttl = ttl
}
func (uriParameters *UriParameters) GetTtl() uint8 {
	return uriParameters.ttl
}
func (uriParameters *UriParameters) SetMaddr(maddr string) {
	uriParameters.maddr = maddr
}
func (uriParameters *UriParameters) GetMaddr() string {
	return uriParameters.maddr
}
func (uriParameters *UriParameters) SetLr(lr bool) {
	uriParameters.lr = lr
}
func (uriParameters *UriParameters) GetLr() bool {
	return uriParameters.lr
}
func (uriParameters *UriParameters) SetOther(other sync.Map) {
	uriParameters.other = other
}
func (uriParameters *UriParameters) GetOther() sync.Map {
	return uriParameters.other
}
func (uriParameters *UriParameters) GetIsOrder() bool {
	return uriParameters.isOrder
}
func (uriParameters *UriParameters) GetOrder() []string {
	result := make([]string, 0)
	if uriParameters.order == nil {
		return result
	}
	for data := range uriParameters.order {
		result = append(result, data)
	}
	return result
}
func (uriParameters *UriParameters) SetSource(source string) {
	uriParameters.source = source
}
func (uriParameters *UriParameters) GetSource() string {
	return uriParameters.source
}
func NewUriParameters(transport string, user string, method string, ttl uint8, maddr string, lr bool, other sync.Map) *UriParameters {
	return &UriParameters{
		transport: transport,
		user:      user,
		method:    method,
		ttl:       ttl,
		maddr:     maddr,
		lr:        lr,
		other:     other,
		isOrder:   false,
		order:     make(chan string, 1024),
	}
}
func (uriParameters *UriParameters) Raw() string {
	result := ""
	if uriParameters.isOrder {
		for data := range uriParameters.order {
			result += fmt.Sprintf(";%s", data)
		}
	} else {
		// transport-param = "transport="( "udp" / "tcp" / "sctp" / "tls"/ other-transport),other-transport = token
		if len(strings.TrimSpace(uriParameters.transport)) > 0 {
			result += fmt.Sprintf(";transport=%s", strings.ToLower(uriParameters.transport))
		}
		// user-param =  "user=" ( "phone" / "ip" / other-user), other-user = token
		if len(strings.TrimSpace(uriParameters.user)) > 0 {
			result += fmt.Sprintf(";user=%s", uriParameters.user)
		}
		// method-param =  "method=" Method
		if len(strings.TrimSpace(uriParameters.method)) > 0 {
			result += fmt.Sprintf(";method=%s", uriParameters.method)
		}
		// ttl-param =  "ttl=" ttl
		if uriParameters.ttl > 0 {
			result += fmt.Sprintf(";ttl=%d", uriParameters.ttl)
		}
		// maddr-param       =  "maddr=" host
		if len(strings.TrimSpace(uriParameters.maddr)) > 0 {
			result += fmt.Sprintf(";maddr=%s", uriParameters.maddr)
		}
		// lr-param          =  "lr"
		if uriParameters.lr {
			result += ";lr"
		}
		// other     sync.Map
		uriParameters.other.Range(func(key, value interface{}) bool {
			if reflect.ValueOf(value).IsValid() {
				if reflect.ValueOf(value).IsZero() {
					result += fmt.Sprintf(";%v", key)
					return true
				}
				result += fmt.Sprintf(";%v=%v", key, value)
				return true
			}
			result += fmt.Sprintf(";%v", key)
			return true
		})
	}
	return result
}
func (uriParameters *UriParameters) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// parameters regexp
	parametersRegexp := regexp.MustCompile(`(?:;.*)`)
	if parametersRegexp.MatchString(raw) {
		uriParameters.source = raw
	}
	// parameters order
	uriParameters.parametersOrder(raw)
	// tranport parameter regexp
	transportRegexp := regexp.MustCompile(`(?i)(transport).*`)
	// user parameter regexp
	userRegexp := regexp.MustCompile(`(?i)(user).*`)
	// method parameter regexp
	methodRegexp := regexp.MustCompile(`(?i)(method).*`)
	// ttl parameter regexp
	ttlRegexp := regexp.MustCompile(`(?i)(ttl).*`)
	// maddr parameter regexp
	maddrRegexp := regexp.MustCompile(`(?i)(maddr).*`)
	// lr parameterregexp
	lrRegexp := regexp.MustCompile(`(?i)(lr).*`)

	rawSlice := strings.Split(raw, ";")
	for _, raws := range rawSlice {
		switch {
		case transportRegexp.MatchString(raws):
			transport := regexp.MustCompile(`(?i)(transport)`).ReplaceAllString(raws, "")
			transport = regexp.MustCompile(`.*=`).ReplaceAllString(transport, "")
			transport = stringTrimPrefixAndTrimSuffix(transport, " ")
			if len(transport) > 0 {
				uriParameters.transport = transport
			}
		case userRegexp.MatchString(raws):
			user := regexp.MustCompile(`(?i)(user)`).ReplaceAllString(raws, "")
			user = regexp.MustCompile(`.*=`).ReplaceAllString(user, "")
			user = stringTrimPrefixAndTrimSuffix(user, " ")
			if len(user) > 0 {
				uriParameters.user = user
			}
		case methodRegexp.MatchString(raws):
			method := regexp.MustCompile(`(?i)(method)`).ReplaceAllString(raws, "")
			method = regexp.MustCompile(`.*=`).ReplaceAllString(method, "")
			method = stringTrimPrefixAndTrimSuffix(method, " ")
			if len(method) > 0 {
				uriParameters.user = method
			}
		case ttlRegexp.MatchString(raws):
			ttlStr := regexp.MustCompile(`(?i)(ttl)`).ReplaceAllString(raws, "")
			ttlStr = regexp.MustCompile(`.*=`).ReplaceAllString(ttlStr, "")
			ttlStr = stringTrimPrefixAndTrimSuffix(ttlStr, " ")
			if len(ttlStr) > 0 {
				ttl, _ := strconv.Atoi(ttlStr)
				uriParameters.ttl = uint8(ttl)
			}
		case maddrRegexp.MatchString(raws):
			maddr := regexp.MustCompile(`(?i)(maddr)`).ReplaceAllString(raws, "")
			maddr = regexp.MustCompile(`.*=`).ReplaceAllString(maddr, "")
			maddr = stringTrimPrefixAndTrimSuffix(maddr, " ")
			if len(maddr) > 0 {
				uriParameters.maddr = maddr
			}
		case lrRegexp.MatchString(raws):
			uriParameters.lr = true
		default:
			if len(strings.TrimSpace(raws)) > 0 {
				if strings.Contains(raws, "=") {
					gs := strings.Split(raws, "=")
					if len(gs) > 1 {
						uriParameters.other.Store(gs[0], gs[1])
					} else {
						uriParameters.other.Store(gs[0], "")
					}
				} else {
					uriParameters.other.Store(raws, "")
				}
			}
		}
	}
}
func (uriParameters *UriParameters) parametersOrder(parameter string) {
	if uriParameters.order == nil {
		uriParameters.order = make(chan string, 1024)
	}
	uriParameters.isOrder = true
	defer close(uriParameters.order)
	parameter = stringTrimPrefixAndTrimSuffix(parameter, ";")
	parameter = stringTrimPrefixAndTrimSuffix(parameter, " ")
	parameters := strings.Split(parameter, ";")
	for _, data := range parameters {
		uriParameters.order <- data
	}
}
