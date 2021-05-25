package sip

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

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
	name   string
	ipv4   net.IP
	ipv6   net.IP
	port   uint16
	source string // source string
}

func (hp *HostPort) SetName(name string) {
	hp.name = name
}
func (hp *HostPort) GetName() string {
	return hp.name
}

func (hp *HostPort) SetIPv4(ipv4 net.IP) {
	hp.ipv4 = ipv4
}
func (hp *HostPort) GetIPv4() net.IP {
	return hp.ipv4
}
func (hp *HostPort) SetIPv6(ipv6 net.IP) {
	hp.ipv6 = ipv6
}
func (hp *HostPort) GetIPv6() net.IP {
	return hp.ipv6
}
func (hp *HostPort) SetPort(port uint16) {
	hp.port = port
}
func (hp *HostPort) GetPort() uint16 {
	return hp.port
}
func (hp *HostPort) GetSource() string {
	return hp.source
}
func NewHostPort(name string, ipv4 net.IP, ipv6 net.IP, port uint16) *HostPort {
	return &HostPort{
		name: name,
		ipv4: ipv4,
		ipv6: ipv6,
		port: port,
	}
}

func (hp *HostPort) Raw() (result strings.Builder) {
	switch {
	case len(strings.TrimSpace(hp.name)) > 0:
		result.WriteString(hp.name)
	case hp.ipv4 != nil:
		result.WriteString(hp.ipv4.String())
	case hp.ipv6 != nil:
		result.WriteString(fmt.Sprintf("["+"%s"+"]", hp.ipv6.String()))
	}
	if hp.port > 0 {
		result.WriteString(fmt.Sprintf(":%d", hp.port))
	}
	return
}
func (hp *HostPort) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	hp.source = raw
	// ipv4 address regexp
	ipv4AddressRegexp := regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
	// host name regexp
	hostnameRegexp := regexp.MustCompile(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`)
	switch {
	case ipv4AddressRegexp.MatchString(raw):
		hp.ipv4 = net.ParseIP(ipv4AddressRegexp.FindString(raw))
		raw = regexp.MustCompile(`.*`+hp.ipv4.String()).ReplaceAllString(raw, "")
	case hostnameRegexp.MatchString(raw):
		hp.name = hostnameRegexp.FindString(raw)
		raw = regexp.MustCompile(`.*`+hp.name).ReplaceAllString(raw, "")
	default:
		ipAddr, err := net.ResolveIPAddr("ip", raw)
		if err == nil {
			if ipAddr.IP.To16() != nil {
				hp.ipv6 = ipAddr.IP.To16()
				raw = regexp.MustCompile(`.*`+hp.ipv6.String()).ReplaceAllString(raw, "")
			} else if ipAddr.IP.To4() != nil {
				hp.ipv4 = ipAddr.IP.To4()
				raw = regexp.MustCompile(`.*`+hp.ipv4.String()).ReplaceAllString(raw, "")
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
				hp.port = uint16(port)
			}
		}
	}
}
