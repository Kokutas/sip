package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestUriParameters_Raw(t *testing.T) {
	uriParameters := NewUriParameters("udp", "34020000001320000001", "REGISTER", 5, "192.168.0.1", true, sync.Map{})
	fmt.Println(uriParameters.Raw())
}

func TestUriParameters_Parse(t *testing.T) {
	raw := ";transport=udp;user=34020000001320000001;method=REGISTER;maddr=192.168.0.1;lr;ttl=5;;"
	uriParameters := new(UriParameters)
	uriParameters.Parse(raw)
	fmt.Println(raw)
	fmt.Println(uriParameters.Raw())
}

func TestHostPort_Raw(t *testing.T) {
	hostport := NewHostPort("www.baidu.com", nil, nil, 5060)
	fmt.Println(hostport.Raw())
}

func TestHostPort_Parse(t *testing.T) {
	raws := []string{
		"192.168.0.5",
		"192.168.0.1",
		"192.168.0.8:5060",
		"www.baidu.com:80",
		"[fe80::d133:ad17:2520:9421]:8060",
	}
	for _, raw := range raws {
		hostport := new(HostPort)
		hostport.Parse(raw)
		if len(hostport.source) > 0 {
			fmt.Println("hostname:", hostport.hostname, ",ipv4:", hostport.ipv4Address.String(), ",ipv6:", hostport.ipv6Reference, ",port:", hostport.port)
			fmt.Println(hostport.Raw())
		}
	}
}

func TestUserInfo_Raw(t *testing.T) {
	userInfo := NewUserInfo("34020000001320000001", "", "Ali12345")
	fmt.Println(userInfo.Raw())
}

func TestUserInfo_Parse(t *testing.T) {
	raws := []string{
		"010-12345678",
		"010-1234567",
		"+12125551212",
		"+12125551212@phone2net.com",
		"+1-212-555-1212:1234@gateway.com;user=phone",
		"+086-13755969903",
		"17521500865:5060",
		"+17521500865:5060",
		"+086-17521500865:5060",
		"86-17521500865:5060",
		"13755969903:abcd@qq.com",
		"+13755969903:xyz",
		"+86-010-40020021",
		"86-010-40020021",
		"010-40020020",
		"+86-13523458056",
		"10-13523458056",
		"34020000001320000001",
		"34020000001320000001:Ali123",
		"+086-0559-6959003:kokutas@163.com",
		"sipabc:Ali123",
	}
	for _, raw := range raws {
		userinfo := new(UserInfo)
		userinfo.Parse(raw)
		if len(userinfo.source) > 0 {
			fmt.Println("user:", userinfo.user, ",telephone:", userinfo.telephoneSubscriber, ",password:", userinfo.password)
		}
	}
}

func TestSipUri_Raw(t *testing.T) {
	var headers sync.Map
	headers.Store("token", "xyz")
	headers.Store("expires", 3600)
	headers.Store("xxxxxxx", nil)
	sipUri := NewSipUri(NewUserInfo("34020000001320000001", "+086-17621400864", "Ali12345"),
		NewHostPort("www.baidu.com", net.IPv4(192, 168, 0, 1), nil, 5060),
		NewUriParameters("udp", "kokutas", "register", 5, "192.168.0.26", true, sync.Map{}), headers)
	fmt.Println(sipUri.Raw())
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
		if len(sipUri.source) > 0 {
			fmt.Println(index, "schema:", sipUri.schema)
			fmt.Println(index, "userinfo-user:", sipUri.userinfo.user)
			fmt.Println(index, "userinfo-telephone-subscriber:", sipUri.userinfo.telephoneSubscriber)
			fmt.Println(index, "userinfo-password:", sipUri.userinfo.password)
			fmt.Println(index, "hostport-hostname:", sipUri.hostport.hostname)
			fmt.Println(index, "hostport-ipv4:", sipUri.hostport.ipv4Address.String())
			fmt.Println(index, "hostport-ipv6:", sipUri.hostport.ipv6Reference.String())
			fmt.Println(index, "hostport-port:", sipUri.hostport.port)
			fmt.Println(index, "uriparameters-transport:", sipUri.uriparameters.transport)
			fmt.Println(index, "uriparameters-user:", sipUri.uriparameters.user)
			fmt.Println(index, "uriparameters-method:", sipUri.uriparameters.method)
			fmt.Println(index, "uriparameters-ttl:", sipUri.uriparameters.ttl)
			fmt.Println(index, "uriparameters-maddr:", sipUri.uriparameters.maddr)
			fmt.Println(index, "uriparameters-lr:", sipUri.uriparameters.lr)
			sipUri.headers.Range(func(key, value interface{}) bool {
				fmt.Println(index, key, value)
				return true
			})
		}
	}
}
