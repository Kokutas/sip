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
