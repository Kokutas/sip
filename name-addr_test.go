package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestNameAddr_Raw(t *testing.T) {
	p := sync.Map{}
	p.Store("lr", "")
	p.Store("hello", nil)
	nas := []*NameAddr{
		NewNameAddr("sip", NewHostPort("", net.IPv4(192, 168, 0, 26), nil, 0), p),
		NewNameAddr("sip", NewHostPort("www.baidu.com", net.IPv4(192, 168, 0, 26), nil, 0), p),
	}
	for _, na := range nas {
		result := na.Raw()
		fmt.Println(result.String())

	}
}

func TestNameAddr_Parse(t *testing.T) {
	raws := []string{
		"sip:192.168.0.26;lr",
		"sip:www.baidu.com;lr",
		"sip:www.baidu.com;lr=3;hello=w;ls",
	}
	for index, raw := range raws {
		na := new(NameAddr)
		na.Parse(raw)
		if len(na.GetSource()) > 0 {
			result := na.Raw()
			fmt.Println(index, result.String())
		}
	}
}
