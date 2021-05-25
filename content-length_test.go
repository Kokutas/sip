package sip

import (
	"fmt"
	"testing"
)

func TestContentLength_Raw(t *testing.T) {
	l := NewContentLength(0)
	result := l.Raw()
	fmt.Println(result.String())
}

func TestContentLength_Parse(t *testing.T) {
	raws := []string{
		"l: 5060",
		"content-length: 0",
		"content-length:60",
	}
	for _, raw := range raws {
		l := new(ContentLength)
		l.Parse(raw)
		if len(l.GetSource()) > 0 {
			result := l.Raw()
			fmt.Println(result.String())
			fmt.Println(l.GetField(), l.GetLength())
		}
	}
}
