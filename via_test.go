package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestVia_Raw(t *testing.T) {
	var generic sync.Map
	generic.Store("hello", "")
	generic.Store("zz", "xx")
	generic.Store("hi", nil)
	generic.Store("heihei", 123)
	v := NewVia("sip", 2.0, "udp", "192.168.0.1", 5060, 5, "192.168.0.108", net.IPv4(192, 168, 0, 26), "z9hG4bK-branch", 1, "udp", generic)
	fmt.Print(v.Raw())
}

func TestVia_Parse(t *testing.T) {
	raws := []string{
		"Via: SIP/2.0/UDP 192.168.0.1:5060;rport;transport=udp;ttl=5;maddr=192.168.0.108;branch=z9hG4bK-branch;received=192.168.0.26\r\n",
		"Via: SIP/2.0/UDP baidu.com:5060;rport;transport=udp;ttl=5;maddr=192.168.0.108;branch=z9hG4bK-branch;received=192.168.0.26\r\n",
		"Via: SIP/2.0/UDP baidu.com;rport;transport=udp;ttl=5;maddr=192.168.0.108;branch=z9hG4bK-branch;received=192.168.0.26\r\n",
		"Via: SIP/2.0/UDP 192.168.0.1:5060;rport;transport=udp;branch=z9hG4bK-branch;received=192.168.0.26;ttl=5;maddr=192.168.0.108;\r\n",
		"Via: SIP/2.0/udp 192.168.0.1;rport;transport=udp;branch=z9hG4bK-branch;received=www.baidu.com;ttl=5;maddr=192.168.0.108;\r\n",
	}
	for _, raw := range raws {
		v := new(Via)
		v.Parse(raw)
		if len(v.source) > 0 {
			fmt.Println(v.version)
			fmt.Print(v.Raw())
		}

	}
}
