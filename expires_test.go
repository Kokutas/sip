package sip

import (
	"fmt"
	"testing"
)

func TestExpires_Raw(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Print(expires.Raw())
}

func TestExpires_Parse(t *testing.T) {
	raw := "expires:3600"
	expires := new(Expires)
	expires.Parse(raw)
	if len(expires.source) > 0 {
		fmt.Print(expires.Raw())
	}
}
