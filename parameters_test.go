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

func Test_Parameters_Raw(t *testing.T) {
	ps := NewParameters("udp", "34020000001320000001", "register", 0, "www.baidu.com", true, map[string]interface{}{"hello": "world"})
	fmt.Println(ps.Raw())
}

func Test_Parameters_JsonString(t *testing.T) {
	ps := NewParameters("udp", "34020000001320000001", "register", 0, "", true, nil)
	if res := ps.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func Test_Parameters_Parser(t *testing.T) {
	raw := ";transport=udp;user=34020000001320000001;method=register;maddr=www.baidu.com;lr;hello=world"
	fmt.Println(raw)
	ps := CreateParameters()
	if err := ps.Parser(raw); err != nil {
		log.Fatal(err)
	}
	if res := ps.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func Test_Parameters_Validator(t *testing.T) {

}
