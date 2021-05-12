package header

import (
	"fmt"
	"log"
	"sip"
	"testing"
)

func TestCSeq_Raw(t *testing.T) {
	cseq := NewCSeq(1, sip.REGISTER)
	fmt.Print(cseq.Raw())
}

func TestCSeq_JsonString(t *testing.T) {
	cseq := NewCSeq(1, sip.REGISTER)
	fmt.Println(cseq.JsonString())
}

func TestCSeq_Parser(t *testing.T) {
	raw := "CSeq: 1 REGISTER\r\n"
	cseq := CreateCSeq()
	if err := cseq.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(cseq.Raw())
	fmt.Println(cseq.JsonString())
}

func TestCSeq_Validator(t *testing.T) {
	// cseq := NewCSeq(1, sip.REGISTER)
	cseq := CreateCSeq()
	fmt.Println(cseq.Validator())
}
