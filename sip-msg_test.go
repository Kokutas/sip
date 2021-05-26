package sip

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSipMsg_Raw(t *testing.T) {
	var (
		uasId            = "34020000002000000001"
		uasIp            = net.IPv4(129, 168, 0, 108)
		uasPort   uint16 = 5060
		uacId            = "34020000001320000001"
		uacIp            = net.IPv4(129, 168, 0, 26)
		uacPort   uint16 = 5060
		expire    uint32 = 3600
		schema           = "sip"
		version          = 2.0
		transport        = "udp"
		method           = "register"
		algorithm        = "MD5"
	)
	reqUri := NewRequestUri(
		NewSipUri(
			NewUserInfo(uasId, "", ""),
			NewHostPort("", uasIp, nil, uasPort),
			nil,
			sync.Map{}))
	reqLine := NewRequestLine(method, reqUri, schema, version)
	from := NewFrom("", "<", schema, uacId, uacIp.String(), uacPort, "123", sync.Map{})
	to := NewTo("", "<", schema, uacId, uacIp.String(), uacPort, "", sync.Map{})
	contact := NewContact("", "<", schema, uasId, uasIp.String(), uasPort, "", -1, sync.Map{})
	callId := NewCallID("abcdefg", uacIp.String())
	via := NewVia(schema, 2.0, transport, uasIp.String(), uasPort, 0, "", "", "xxxxx", 0, "", sync.Map{})
	expires := NewExpires(expire)
	maxForwards := NewMaxForwards(70)
	contentLength := NewContentLength(0)
	sm := new(SipMsg)
	sm.SetRequestLine(reqLine)
	sm.SetFrom(from)
	sm.SetTo(to)
	sm.SetContact(contact)
	sm.SetCallID(callId)
	sm.SetVia(via)
	sm.SetExpires(expires)
	sm.SetMaxForwards(maxForwards)
	sm.SetContentLength(contentLength)
	result := sm.Raw()
	fmt.Println("----------------------------REGISTER REQUEST----------------------------")
	fmt.Print(result.String())
	// 100
	fmt.Println("----------------------------100 Trying RESPONSE----------------------------")
	sm.SetRequestLine(nil)
	statusLine := NewStatusLine(schema, version, 100, Informational[100])
	sm.SetStatusLine(statusLine)
	result = sm.Raw()
	fmt.Print(result.String())

	// 401 challenge
	fmt.Println("----------------------------401 Unauthorized RESPONSE----------------------------")
	sm.SetRequestLine(nil)
	statusLine = NewStatusLine(schema, version, 401, ClientError[401])
	sm.GetTo().SetTag("456")
	nonce := GenNonce(uacIp.String(), callId.GetSource())
	wwwAuthenticate := NewWWWAuthenticate(uasId[:10], "", nonce, "", false, algorithm, "", sync.Map{})
	sm.SetStatusLine(statusLine)
	sm.SetWWWAuthenticate(wwwAuthenticate)
	result = sm.Raw()
	fmt.Print(result.String())
	// authorization
	fmt.Println("----------------------------REGISTER AUTHORIZATION REQUEST----------------------------")
	sm.SetStatusLine(nil)
	sm.SetWWWAuthenticate(nil)
	reqUriRaw := reqUri.Raw()
	dp := &DigestParams{
		Algorithm: algorithm,
		Method:    method,
		URI:       reqUriRaw.String(),
		Nonce:     nonce,
	}
	response := GenDigestResponse(dp)
	authorization := NewAuthorization(uacId, uasId[:10], nonce, reqUri, response, algorithm, "", "", "", "", sync.Map{})
	sm.SetRequestLine(reqLine)
	sm.SetAuthorization(authorization)
	result = sm.Raw()
	fmt.Print(result.String())

}
