package sip

import (
	"encoding/json"
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
	Port uint16 `json:"port,omitempty"`
}

func CreateHostPort() Sip {
	return &HostPort{}
}

func NewHostPort(host *Host, port uint16) Sip {
	return &HostPort{
		Host: host,
		Port: port,
	}
}
func (hp *HostPort) Raw() string {
	result := ""
	if hp == nil {
		return result
	}
	if hp.Host != nil {
		result += hp.Host.Raw()
	}
	if hp.Port > 0 {
		result += fmt.Sprintf(":%v", hp.Port)
	}
	return result
}
func (hp *HostPort) JsonString() string {
	result := ""
	if hp == nil {
		return result
	}
	data, err := json.Marshal(hp)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (hp *HostPort) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if hp == nil {
		return errors.New("HostPort caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	if regexp.MustCompile(`\:\d+$`).MatchString(raw) {
		portRaw := regexp.MustCompile(`\:\d+$`).FindString(raw)
		portStr := strings.TrimPrefix(portRaw, ":")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		hp.Port = uint16(port)
		raw = strings.Replace(raw, portRaw, "", 1)
	}
	hp.Host = CreateHost().(*Host)
	return hp.Host.Parser(raw)
}
func (hp *HostPort) Validator() error {
	if hp == nil {
		return errors.New("HostPort caller is not allowed to be nil")
	}
	return hp.Host.Validator()
}
