package header

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sip"
	"testing"
)

func TestCallID_JsonString(t *testing.T) {
	cid := NewCallID("ms1214-322164710-681262131542511620107-0", sip.NewHost("", net.IPv4(192, 168, 0, 12), nil).(*sip.Host))
	fmt.Println(cid.JsonString())
}

func TestCallID_Parser(t *testing.T) {
	raws := []string{
		"Call-ID: ms1214-322164710-681262131542511620107-0\r\n",
		"Call-ID: ms1214-322164710-681262131542511620107-0@3402000000\r\n",
		"Call-ID: ms1214-322164710-681262131542511620107-0@192.168.0.26\r\n",
	}
	cid := CreateCallID()
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
		NewCallID("ms1214-322164710-681262131542511620107-0", sip.NewHost("3402000000", nil, nil).(*sip.Host)),
		NewCallID("ms1214-322164710-681262131542511620107-0", sip.NewHost("", net.IPv4(192, 168, 0, 26), nil).(*sip.Host)),
	}
	for _, cid := range cids {
		fmt.Print(cid.Raw())
	}
}

func TestCallID_Validator(t *testing.T) {
	cid := NewCallID("ms1214-322164710-681262131542511620107-0", nil)
	fmt.Println(cid.Validator())
}

func TestNewCallID(t *testing.T) {
	cid := NewCallID("ms1214-322164710-681262131542511620107-0", nil)
	data, err := json.Marshal(cid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}
