package gb28181

import (
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/kokutas/sip"
)

func TestServer_Response(t *testing.T) {
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
	)
	reqUri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.NewUserInfo(uasId, "", ""),
			sip.NewHostPort("", uasIp, nil, uasPort),
			nil,
			sync.Map{}))
	reqLine := sip.NewRequestLine(method, reqUri, schema, version)
	from := sip.NewFrom("", "<", schema, uacId, uacIp.String(), uacPort, "123", sync.Map{})
	to := sip.NewTo("", "<", schema, uacId, uacIp.String(), uacPort, "", sync.Map{})
	contact := sip.NewContact("", "<", schema, uasId, uasIp.String(), uasPort, "", -1, sync.Map{})
	callId := sip.NewCallID("abcdefg", uacIp.String())
	via := sip.NewVia(schema, 2.0, transport, uasIp.String(), uasPort, 0, "", "", "xxxxx", 1, "", sync.Map{})
	expires := sip.NewExpires(expire)
	maxForwards := sip.NewMaxForwards(70)
	contentLength := sip.NewContentLength(0)
	userAgent := sip.NewUserAgent("SIP", "UAC-IPC", "com.kokutas", "V1.0.0")
	sm := new(sip.SipMsg)
	sm.SetRequestLine(reqLine)
	sm.SetFrom(from)
	sm.SetTo(to)
	sm.SetContact(contact)
	sm.SetCallID(callId)
	sm.SetVia(via)
	sm.SetExpires(expires)
	sm.SetMaxForwards(maxForwards)
	sm.SetContentLength(contentLength)
	sm.SetUserAgent(userAgent)
	result := sm.Raw()
	fmt.Print(result.String())
	server := NewServer(uasId, uasId[:10], uasIp, 5060, "udp")

	// 401 challenge
	fmt.Println("----------------------------401 Unauthorized RESPONSE----------------------------")
	result = server.Response(sm)
	fmt.Print(result.String())
}
