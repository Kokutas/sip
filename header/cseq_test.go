package header

import (
	"fmt"
	"log"
	"github.com/kokutas/sip"
	"testing"
)

func TestCSeq_Field(t *testing.T) {
}

func TestCSeq_Method(t *testing.T) {

}

func TestCSeq_Parser(t *testing.T) {
	raw := "CSeq: 1 REGISTER\r\n"
	cseq := new(CSeq)
	if err := cseq.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(cseq.Raw())
	fmt.Println(cseq.String())
}

func TestCSeq_Raw(t *testing.T) {
	cseq := NewCSeq(1, sip.REGISTER)
	fmt.Print(cseq.Raw())
}

func TestCSeq_Sequence(t *testing.T) {

}

func TestCSeq_SetField(t *testing.T) {

}

func TestCSeq_SetMethod(t *testing.T) {

}

func TestCSeq_SetSequence(t *testing.T) {

}

func TestCSeq_String(t *testing.T) {
	cseq := NewCSeq(1, sip.REGISTER)
	fmt.Println(cseq.String())
}

func TestCSeq_Validator(t *testing.T) {
	// cseq := NewCSeq(1, sip.REGISTER)
	cseq := new(CSeq)
	fmt.Println(cseq.Validator())
}

func TestNewCSeq(t *testing.T) {

}
