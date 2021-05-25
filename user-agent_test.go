package sip

import (
	"fmt"
	"testing"
)

func TestUserAgent_Raw(t *testing.T) {
	ua := NewUserAgent("Softphone", "Beta1.5")
	result := ua.Raw()
	fmt.Print(result.String())
}

func TestUserAgent_Parse(t *testing.T) {
	raws := []string{
		"User-Agent: Softphone Beta1.5\r\n",
		"User-Agent: Uas-x v1.0.0\r\n",
		"User-Agent: Uas-x\r\n",
		"User-Agent: \r\n",
	}
	for _, raw := range raws {
		ua := new(UserAgent)
		ua.Parse(raw)
		if len(ua.GetSource()) > 0 {
			result := ua.Raw()
			fmt.Print(result.String())
		}

	}
}
