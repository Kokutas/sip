package header

import (
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestCallID_CallId(t *testing.T) {

}

func TestCallID_Field(t *testing.T) {

}

func TestCallID_Parser(t *testing.T) {
	raws := []string{
		"Call-ID: ms1214-322164710-681262131542511620107-0\r\n",
		"Call-ID: ms1214-322164710-681262131542511620107-0@3402000000\r\n",
		"Call-ID: ms1214-322164710-681262131542511620107-0@192.168.0.26\r\n",
	}
	cid := new(CallID)
	for _, raw := range raws {
		if err := cid.Parser(raw); err != nil {
			log.Fatal(err)
		}
		fmt.Println(cid.Raw())
	}
}

func TestCallID_Raw(t *testing.T) {
	cids := []*CallID{
		NewCallID("ms1214-322164710-681262131542511620107-0", nil),
		NewCallID("ms1214-322164710-681262131542511620107-0", sip.NewHost("3402000000", nil, nil)),
		NewCallID("ms1214-322164710-681262131542511620107-0", sip.NewHost("", net.IPv4(192, 168, 0, 26), nil)),
	}
	for _, cid := range cids {
		fmt.Print(cid.Raw())
	}
}

func TestCallID_SetCallId(t *testing.T) {

}

func TestCallID_SetField(t *testing.T) {

}

func TestCallID_String(t *testing.T) {
	cid := NewCallID("ms1214-322164710-681262131542511620107-0", sip.NewHost("", net.IPv4(192, 168, 0, 12), nil))
	fmt.Println(cid.String())
}

func TestCallID_Validator(t *testing.T) {
	cid := NewCallID("ms1214-322164710-681262131542511620107-0", nil)
	fmt.Println(cid.Validator())
}

func TestNewCallID(t *testing.T) {
	cid := NewCallID("ms1214-322164710-681262131542511620107-0", nil)
	fmt.Printf("%s\r\n", cid)
}
