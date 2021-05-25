package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSipUri_Raw(t *testing.T) {
	var headers sync.Map
	headers.Store("token", "xyz")
	headers.Store("expires", 3600)
	headers.Store("xxxxxxx", nil)
	sipUri := NewSipUri(NewUserInfo("34020000001320000001", "+086-17621400864", "Ali12345"),
		NewHostPort("www.baidu.com", net.IPv4(192, 168, 0, 1), nil, 5060),
		NewParameters("udp", "kokutas", "register", 5, "192.168.0.26", true, sync.Map{}), headers)
	result := sipUri.Raw()
	fmt.Println(result.String())
}

func TestSipUri_Parse(t *testing.T) {
	raws := []string{
		"sip:34020000001320000001:Ali12345@www.baidu.com:5060;transport=udp;user=kokutas;method=register;ttl=5;maddr=192.168.0.26;lr?token=xyz&expires=3600&xxxxxxx",
		"sip:34020000001320000001@192.168.0.6:5060;transport=udp;user=kokutas;method=register;ttl=5;maddr=192.168.0.26;lr?token=xyz&expires=3600&xxxxxxx",
		"sip:34020000001320000001@192.168.0.6:5060",
	}
	for index, raw := range raws {
		sipUri := new(SipUri)
		sipUri.Parse(raw)
		if len(sipUri.GetSource()) > 0 {
			fmt.Println(index, "schema:", sipUri.schema)
			fmt.Println(index, "userinfo-user:", sipUri.userinfo.user)
			fmt.Println(index, "userinfo-telephone-subscriber:", sipUri.userinfo.telephoneSubscriber)
			fmt.Println(index, "userinfo-password:", sipUri.userinfo.password)
			fmt.Println(index, "hostport-hostname:", sipUri.hostport.GetName())
			fmt.Println(index, "hostport-ipv4:", sipUri.hostport.GetIPv4().String())
			fmt.Println(index, "hostport-ipv6:", sipUri.hostport.GetIPv6().String())
			fmt.Println(index, "hostport-port:", sipUri.hostport.port)
			fmt.Println(index, "uri-parameters-transport:", sipUri.parameters.transport)
			fmt.Println(index, "uri-parameters-user:", sipUri.parameters.user)
			fmt.Println(index, "uri-parameters-method:", sipUri.parameters.method)
			fmt.Println(index, "uri-parameters-ttl:", sipUri.parameters.ttl)
			fmt.Println(index, "uri-parameters-maddr:", sipUri.parameters.maddr)
			fmt.Println(index, "uri-parameters-lr:", sipUri.parameters.lr)
			sipUri.headers.Range(func(key, value interface{}) bool {
				fmt.Println(index, key, value)
				return true
			})
			result := sipUri.Raw()
			fmt.Println(index, result.String())
		}
	}
}
