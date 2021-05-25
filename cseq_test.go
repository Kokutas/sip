package sip

import (
	"fmt"
	"testing"
)

func TestCSeq_Raw(t *testing.T) {
	cseq := NewCSeq(0, "register")
	result := cseq.Raw()
	fmt.Print(result.String())
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
		if len(cseq.GetSource()) > 0 {
			fmt.Println(cseq.GetField(), cseq.GetNumber(), cseq.GetMethod())
			num := cseq.GetNumber()
			num++
			cseq.SetNumber(num)
			result := cseq.Raw()
			fmt.Print(result.String())
		}
	}

}
