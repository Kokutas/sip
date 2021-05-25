package sip

import (
	"fmt"
	"testing"
)

func TestCallID_Raw(t *testing.T) {
	callId := NewCallID("abc", "192.168.0.1")
	result := callId.Raw()
	fmt.Print(result.String())
}

func TestCallID_Parse(t *testing.T) {
	raws := []string{
		"i: 192.156",
		"call-id: hello",
		"call-id: hello@192.168.0.1",
	}
	for index, raw := range raws {
		callId := new(CallID)
		callId.Parse(raw)
		if len(callId.GetSource()) > 0 {
			fmt.Println(index, callId.GetLocalId(), callId.GetHost())
			result := callId.Raw()
			fmt.Print(result.String())
		}

	}
}
