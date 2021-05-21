package sip

import (
	"fmt"
	"testing"
)

func TestCSeq_Raw(t *testing.T) {
	cseq := NewCSeq(0, "register")
	fmt.Print(cseq.Raw())
}

func TestCSeq_Parse(t *testing.T) {
	raws := []string{
		"cseq: 1   Register",
		"cseQ   : 1   Register",
		"CSeq: 0 REGISTER\r\n\r\n",
	}
	for _, raw := range raws {
		cseq := new(CSeq)
		cseq.Parse(raw)
		if len(cseq.source) > 0 {
			fmt.Print(cseq.Raw())
		}

	}
}
