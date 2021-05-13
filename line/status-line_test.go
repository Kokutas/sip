package line

import (
	"fmt"
	"log"
	"github.com/kokutas/sip"
	"testing"
)

func TestNewStatusLine(t *testing.T) {
	sl := NewStatusLine(sip.NewSipVersion(sip.SIP, 2.0), 604, sip.GlobalFailure[604])
	fmt.Printf("%s\r\n", sl)
}

func TestStatusLine_Parser(t *testing.T) {
	raw := "SIP/2.0 604 Does not exist anywhere"
	sl := new(StatusLine)
	if err := sl.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Println(sl.String())
}

func TestStatusLine_Raw(t *testing.T) {
	sl := NewStatusLine(sip.NewSipVersion(sip.SIP, 2.0), 604, sip.GlobalFailure[604])
	fmt.Print(sl.Raw())
}

func TestStatusLine_ReasonPhrase(t *testing.T) {
}

func TestStatusLine_SetReasonPhrase(t *testing.T) {
}

func TestStatusLine_SetStatusCode(t *testing.T) {
}

func TestStatusLine_StatusCode(t *testing.T) {
}

func TestStatusLine_String(t *testing.T) {
	sl := NewStatusLine(sip.NewSipVersion(sip.SIP, 2.0), 604, sip.GlobalFailure[604])
	fmt.Print(sl.String())
}

func TestStatusLine_Validator(t *testing.T) {
	sl := new(StatusLine)
	if err := sl.Validator(); err != nil {
		log.Fatal(err)
	}
}
