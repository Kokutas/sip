package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestRoute_Raw(t *testing.T) {
	p := sync.Map{}
	p.Store("lr", nil)
	routers := []*Route{
		NewRoute(
			NewNameAddr("sip", NewHostPort("www.baidu.com", nil, nil, 5060), p),
			NewNameAddr("sip", NewHostPort("www.163.com", nil, nil, 0), sync.Map{}),
			NewNameAddr("sip", NewHostPort("", net.IPv4(192, 168, 0, 26), nil, 5060), sync.Map{}),
		),
	}
	for _, route := range routers {
		result := route.Raw()
		fmt.Print(result.String())

	}
}

func TestRoute_Parse(t *testing.T) {
	raws := []string{
		"Route: <sip:www.baidu.com:5060;lr>, <sip:www.163.com>, <sip:192.168.0.26:5060>",
		"Route: <sip:www.163.com>, <sip:192.168.0.26:5060>",
		"Route: <sip:www.baidu.com:5060;lr>",
		"Route: <sip:www.baidu.com:5060;lr>,,",
	}
	for index, raw := range raws {
		route := new(Route)
		route.Parse(raw)
		if len(route.GetSource()) > 0 {
			result := route.Raw()
			fmt.Print(index, "-----", result.String())
		}
	}
}
