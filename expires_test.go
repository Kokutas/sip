package sip

import (
	"fmt"
	"testing"
)

func TestExpires_Raw(t *testing.T) {
	expires := NewExpires(3600)
	result := expires.Raw()
	fmt.Print(result.String())
}

func TestExpires_Parse(t *testing.T) {
	raw := "expires:3600"
	expires := new(Expires)
	expires.Parse(raw)
	if len(expires.GetSource()) > 0 {
		result := expires.Raw()
		fmt.Print(result.String())
	}
}
