package sip

import (
	"fmt"
	"testing"
)

func TestStatusLine_Raw(t *testing.T) {
	statusLine := NewStatusLine("sip", 2.0, 200, "OK")
	fmt.Print(statusLine.Raw())
}

func TestStatusLine_Parse(t *testing.T) {
	raws := []string{
		"sip/2.0 200 ok\r\n",
	}
	for _, raw := range raws {
		statusLine := new(StatusLine)
		statusLine.Parse(raw)
		fmt.Println(statusLine.GetSchema(), statusLine.GetVersion(), statusLine.GetStatusCode(), statusLine.GetReasonPhrase())
	}
}
