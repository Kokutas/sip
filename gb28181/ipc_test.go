package gb28181

import (
	"fmt"
	"net"
	"testing"

	"github.com/kokutas/sip"
)

func TestIPC_Request(t *testing.T) {
	ipc := NewIPC("34020000002000000001", net.IPv4(192, 168, 0, 108), 5060, "34020000001320000001", net.IPv4(192, 168, 0, 26), 5060, "udp", 3600)
	// register request
	fmt.Println("----------------------------REGISTER REQUEST----------------------------")
	result := ipc.Request("register", new(sip.SipMsg))
	fmt.Print(result.String())
	// register response 401
	fmt.Println("----------------------------REGISTER REQUEST WITH AUTHORIZATION----------------------------")
	ipc.SetNonce(sip.GenNonce(net.IPv4(192, 168, 0, 108).String(), "123"))
	result = ipc.Request("register", new(sip.SipMsg))
	fmt.Print(result.String())

}
