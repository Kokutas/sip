package line

import (
	"fmt"
	"log"
	"sip"
	"testing"
)

func TestNewStatusLine(t *testing.T) {
	sl := NewStatusLine(sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion), 604, sip.GlobalFailure[604])
	fmt.Println(sl.Raw())
}

func TestStatusLine_JsonString(t *testing.T) {
	sl := NewStatusLine(sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion), 604, sip.GlobalFailure[604])
	fmt.Print(sl.JsonString())
}

func TestStatusLine_Parser(t *testing.T) {
	raw := "SIP/2.0 604 Does not exist anywhere"
	sl := CreateRequestLine()
	if err := sl.Parser(raw); err != nil {
		log.Fatal(err)
	}
	if res := sl.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func TestStatusLine_Raw(t *testing.T) {
	sl := NewStatusLine(sip.NewSipVersion(sip.SIP, 2.0).(*sip.SipVersion), 604, sip.GlobalFailure[604])
	fmt.Print(sl.Raw())
}

func TestStatusLine_Validator(t *testing.T) {
	sl := CreateStatusLine()
	if err := sl.Validator(); err != nil {
		log.Fatal(err)
	}
}
