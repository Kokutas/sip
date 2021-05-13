package sip

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewParameters(t *testing.T) {
	ps := NewParameters("", "", "", 0, "", false, nil)
	data, err := json.Marshal(ps)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestParameters_Lr(t *testing.T) {}

func TestParameters_Maddr(t *testing.T) {}

func TestParameters_Method(t *testing.T) {}

func TestParameters_Other(t *testing.T) {}

func TestParameters_Parser(t *testing.T) {
	raw := ";transport=udp;user=34020000001320000001;method=register;maddr=www.baidu.com;lr;hello=world"
	fmt.Println(raw)
	ps := new(Parameters)
	if err := ps.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Println(ps.String())
}

func TestParameters_Raw(t *testing.T) {
	ps := NewParameters("udp", "34020000001320000001", "register", 0, "www.baidu.com", true, map[string]interface{}{"hello": "world"})
	fmt.Println(ps.Raw())
}

func TestParameters_SetLr(t *testing.T) {}

func TestParameters_SetMaddr(t *testing.T) {}

func TestParameters_SetMethod(t *testing.T) {}

func TestParameters_SetOther(t *testing.T) {
}

func TestParameters_SetTransport(t *testing.T) {
}

func TestParameters_SetTtl(t *testing.T) {
}

func TestParameters_SetUser(t *testing.T) {}

func TestParameters_String(t *testing.T) {
	ps := NewParameters("udp", "34020000001320000001", "register", 0, "", true, nil)
	fmt.Println(ps.String())
}

func TestParameters_Transport(t *testing.T) {}

func TestParameters_Ttl(t *testing.T) {}

func TestParameters_User(t *testing.T) {}

func TestParameters_Validator(t *testing.T) {}
