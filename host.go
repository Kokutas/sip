package sip

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Host struct {
	Hostname      string `json:"hostname"`
	IPv4address   net.IP `json:"ipv4address"`
	IPv6reference net.IP `json:"ipv6reference"`
}

func CreateHost() Sip {
	return &Host{}
}
func NewHost(hostname string, ipv4, ipv6 net.IP) Sip {
	return &Host{
		Hostname:      hostname,
		IPv4address:   ipv4,
		IPv6reference: ipv6,
	}
}
func (h *Host) Raw() string {
	result := ""
	if h == nil {
		return result
	}
	switch {
	case len(strings.TrimSpace(h.Hostname)) > 0:
		result += h.Hostname
	case h.IPv4address != nil:
		result += h.IPv4address.String()
	case h.IPv6reference != nil:
		result += h.IPv6reference.String()
	}
	return result
}
func (h *Host) JsonString() string {
	result := ""
	if h == nil {
		return result
	}
	data, err := json.Marshal(h)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (h *Host) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if h == nil {
		return errors.New("Host caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	ip := net.ParseIP(raw)
	if ip == nil {
		h.Hostname = raw
	} else {
		if ipv4 := ip.To4(); ipv4 != nil {
			h.IPv4address = ipv4
		} else if ipv6 := ip.To16(); ipv6 != nil {
			h.IPv6reference = ipv6
		} else {
			h.Hostname = raw
		}
	}
	return nil
}
func (h *Host) Validator() error {
	if h == nil {
		return errors.New("Host caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(h.Hostname)) == 0 && h.IPv4address == nil && h.IPv6reference == nil {
		return errors.New("hostname or IPv4address or IPv6reference must has one")
	}
	return nil
}
