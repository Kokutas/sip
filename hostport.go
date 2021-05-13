package sip

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

type HostPort struct {
	*Host
	port uint16
}

func (hostPort *HostPort) Port() uint16 {
	return hostPort.port
}

func (hostPort *HostPort) SetPort(port uint16) {
	hostPort.port = port
}

func NewHostPort(host *Host, port uint16) *HostPort {
	return &HostPort{
		Host: host,
		port: port,
	}
}
func (hostPort *HostPort) Raw() (string, error) {
	result := ""
	if err := hostPort.Validator(); err != nil {
		return result, err
	}
	if hostPort.Host != nil {
		res, err := hostPort.Host.Raw()
		if err != nil {
			return result, err
		}
		result += res
	}
	if hostPort.port > 0 {
		result += fmt.Sprintf(":%v", hostPort.port)
	}
	return result, nil
}
func (hostPort *HostPort) String() string {
	result := ""
	if hostPort.Host != nil {
		result += fmt.Sprintf("%s,", hostPort.Host.String())
	}
	if hostPort.port > 0 {
		result += fmt.Sprintf("port: %v,", hostPort.port)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (hostPort *HostPort) Parser(raw string) error {
	if hostPort == nil {
		return errors.New("hostPort caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw," ")
	raw = strings.TrimPrefix(raw," ")
	raw = strings.TrimSuffix(raw," ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	if regexp.MustCompile(`:\d+$`).MatchString(raw) {
		portRaw := regexp.MustCompile(`:\d+$`).FindString(raw)
		portStr := strings.TrimPrefix(portRaw, ":")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		hostPort.port = uint16(port)
		raw = strings.Replace(raw, portRaw, "", 1)
	}
	hostPort.Host = new(Host)
	return hostPort.Host.Parser(raw)
}
func (hostPort *HostPort) Validator() error {
	if hostPort == nil {
		return errors.New("hostPort caller is not allowed to be nil")
	}
	return hostPort.Host.Validator()
}
