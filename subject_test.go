package sip

import (
	"fmt"
	"testing"
)

func TestSubject_Raw(t *testing.T) {
	s := NewSubject("hello")
	result := s.Raw()
	fmt.Print(result.String())
}

func TestSubject_Parse(t *testing.T) {
	raws := []string{
		"Subject: hello\r\n",
		"s: world\r\n",
	}
	for _, raw := range raws {
		s := new(Subject)
		s.Parse(raw)
		if len(s.GetSource()) > 0 {
			fmt.Println(s.GetField(), s.GetText())
			result := s.Raw()
			fmt.Print(result.String())
		}

	}
}
