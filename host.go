package sip

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type Host struct {
	hostName      string
	ipv4Address   net.IP
	ipv6Reference net.IP
}

func (host *Host) HostName() string {
	return host.hostName
}

func (host *Host) SetHostName(hostName string) {
	host.hostName = hostName
}

func (host *Host) Ipv4Address() net.IP {
	return host.ipv4Address
}

func (host *Host) SetIpv4Address(ipv4Address net.IP) {
	host.ipv4Address = ipv4Address
}

func (host *Host) Ipv6Reference() net.IP {
	return host.ipv6Reference
}

func (host *Host) SetIpv6Reference(ipv6Reference net.IP) {
	host.ipv6Reference = ipv6Reference
}

func NewHost(hostName string, ipv4, ipv6 net.IP) *Host {
	return &Host{
		hostName:      hostName,
		ipv4Address:   ipv4,
		ipv6Reference: ipv6,
	}
}

func (host *Host) Raw() (string, error) {
	result := ""
	if err := host.Validator(); err != nil {
		return result, err
	}
	switch {
	case len(strings.TrimSpace(host.hostName)) > 0:
		result += host.hostName
	case host.ipv4Address != nil:
		result += host.ipv4Address.String()
	case host.ipv6Reference != nil:
		result += host.ipv6Reference.String()
	}
	return result, nil
}
func (host *Host) String() string {
	result:=""
	if len(strings.TrimSpace(host.hostName))>0{
		result+=fmt.Sprintf("hostname: %s,",host.hostName)
	}
	if host.ipv4Address!=nil{
		result+=fmt.Sprintf("ipv4-address: %s,",host.ipv4Address.String())
	}
	if host.ipv6Reference!=nil{
		result+=fmt.Sprintf("ipv6-reference: %s,",host.ipv6Reference.String())
	}
	result=strings.TrimSuffix(result,",")
	return result
}

func (host *Host) Parser(raw string) error {
	if host == nil {
		return errors.New("host caller is not allowed to be nil")
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
	ip := net.ParseIP(raw)
	if ip == nil {
		host.hostName = raw
	} else {
		if ipv4 := ip.To4(); ipv4 != nil {
			host.ipv4Address = ipv4
		} else if ipv6 := ip.To16(); ipv6 != nil {
			host.ipv6Reference = ipv6
		} else {
			host.hostName = raw
		}
	}
	return nil
}
func (host *Host) Validator() error {
	if host == nil {
		return errors.New("host caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(host.hostName)) == 0 && host.ipv4Address == nil && host.ipv6Reference == nil {
		return errors.New("hostName or ipv4Address or ipv6Reference must has one")
	}
	return nil
}
