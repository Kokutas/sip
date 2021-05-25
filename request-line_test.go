package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestRequestLine_Raw(t *testing.T) {
	reqLines := []*RequestLine{
		NewRequestLine("register",
			NewRequestUri(
				NewSipUri(NewUserInfo("34020000001320000001", "17631300986", ""), NewHostPort("", net.IPv4(192, 168, 0, 26), nil, 5060), NewParameters("udp", "", "", 0, "", false, sync.Map{}), sync.Map{})), "sip", 2.0),
		NewRequestLine("INvite",
			NewRequestUri(
				NewSipUri(NewUserInfo("34020000001320000001", "17631300989", "xxYYzz123"), NewHostPort("", net.IPv4(192, 168, 0, 26), nil, 5060), NewParameters("udp", "kokutas", "invite", 0, "192.168.0.1", false, sync.Map{}), sync.Map{})), "sip", 2.0),
	}
	for _, reqLine := range reqLines {
		result := reqLine.Raw()
		fmt.Print(result.String())

	}
}

func TestRequestLine_Parse(t *testing.T) {
	raws := []string{
		"REGISTER sip:34020000001320000001@192.168.0.26:5060;transport=udp;hello=world?token=xxyyzz&abc SIP/2.0\r\n",
		"REGISTER sip:17621690968@www.baidu.com:5060;transport=udp;hello=world?token=xxyyzz&abc SIP/2.0\r\n",
		"REGISTER sip:+086-17621690968@192.168.0.26:5060;transport=udp;hello=world?token=xxyyzz&abc SIP/2.0\r\n",
		"INVITE sip:34020000001320000001:xxYYzz123@192.168.0.26:5060;transport=udp;user=kokutas;method=invite;maddr=192.168.0.1 SIP/2.0\r\n",
	}
	for index, raw := range raws {
		reqLine := new(RequestLine)
		reqLine.Parse(raw)
		if len(reqLine.source) > 0 {
			fmt.Println(index, "method:", reqLine.method)
			fmt.Println(index, "request-uri->sip/uri->schema:", reqLine.uri.sipUri.schema)
			fmt.Println(index, "request-uri->sip/uri->userinfo->user:", reqLine.uri.sipUri.userinfo.user)
			fmt.Println(index, "request-uri->sip/uri->userinfo->telephone-subscriber:", reqLine.uri.sipUri.userinfo.telephoneSubscriber)
			fmt.Println(index, "request-uri->sip/uri->userinfo->password:", reqLine.GetUri().GetSipUri().GetUserInfo().GetPassword())
			fmt.Println(index, "request-uri->sip/uri->hostport->hostname:", reqLine.GetUri().GetSipUri().GetHostPort().GetName())
			fmt.Println(index, "request-uri->sip/uri->hostport->ipv4:", reqLine.uri.sipUri.hostport.GetIPv4().String())
			fmt.Println(index, "request-uri->sip/uri->hostport->ipv6:", reqLine.uri.sipUri.hostport.GetIPv6().String())
			fmt.Println(index, "request-uri->sip/uri->hostport->port:", reqLine.uri.sipUri.hostport.port)
			fmt.Println(index, "request-uri->sip/uri->uri-parameters->transport:", reqLine.uri.sipUri.parameters.transport)
			fmt.Println(index, "request-uri->sip/uri->uri-parameters->user:", reqLine.uri.sipUri.parameters.user)
			fmt.Println(index, "request-uri->sip/uri->uri-parameters->method:", reqLine.uri.sipUri.parameters.method)
			fmt.Println(index, "request-uri->sip/uri->uri-parameters->ttl:", reqLine.uri.sipUri.parameters.ttl)
			fmt.Println(index, "request-uri->sip/uri->uri-parameters->maddr:", reqLine.uri.sipUri.parameters.maddr)
			fmt.Println(index, "request-uri->sip/uri->uri-parameters->lr:", reqLine.uri.sipUri.parameters.lr)
			reqLine.uri.sipUri.parameters.other.Range(func(key, value interface{}) bool {
				fmt.Println(index, "request-uri->sip/uri->uri-parameters->other:", key, value)
				return true
			})
			reqLine.uri.sipUri.headers.Range(func(key, value interface{}) bool {
				fmt.Println(index, "request-uri->sip/uri->headers:", key, value)
				return true
			})
			fmt.Println(index, "schema:", reqLine.schema)
			fmt.Println(index, "version:", reqLine.version)
			result := reqLine.Raw()
			fmt.Print(index, "-------", result.String())
		}

	}
}
