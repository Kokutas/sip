package sip

import (
	"fmt"
	"net"
	"testing"
)

func TestHostPort_Raw(t *testing.T) {
	hostports := []*HostPort{
		NewHostPort("www.baidu.com", nil, nil, 5060),
		NewHostPort("www.baidu.com", net.IPv4(192, 168, 0, 26), nil, 5060),
		NewHostPort("", net.IPv4(192, 168, 0, 26), nil, 5060),
		NewHostPort("", net.IPv4(192, 168, 0, 26), nil, 0),
	}
	for _, hostport := range hostports {
		result := hostport.Raw()
		fmt.Println(result.String())
	}

}

func TestHostPort_Parse(t *testing.T) {
	raws := []string{
		"www.baidu.com:5060",
		"www.baidu.com:5060",
		"192.168.0.26:5060",
		"192.168.0.26",
		"192.168.0.5",
		"192.168.0.1",
		"192.168.0.8:5060",
		"www.baidu.com:80",
		"[fe80::d133:ad17:2520:9421]:8060",
	}
	for index, raw := range raws {
		hostport := new(HostPort)
		hostport.Parse(raw)
		if len(hostport.GetSource()) > 0 {
			fmt.Println("index: ", index, ",host-name: ", hostport.GetName(), ",ipv4: ", hostport.GetIPv4(), ",ipv6: ", hostport.GetIPv6(), ",port :", hostport.GetPort())
			result := hostport.Raw()
			fmt.Println("index: ", index, ",raw: ", result.String())
		}
	}
}
