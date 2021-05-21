package sip

import (
	"fmt"
	"testing"
)

func TestTo_Raw(t *testing.T) {
	ts := []*To{
		NewTo("", "<", "sip", "34020000001320000001", "192.168.0.1", "5060", "tag123", ""),
		NewTo("34020000001320000001", "'", "sip", "34020000001320000001", "www.baidu.com", "", "tag123", ""),
		NewTo("tom", "\"", "sip", "34020000001320000001", "www.baidu.com", "", "tag123", ""),
		NewTo("alisa", "", "sip", "34020000001320000001", "www.baidu.com", "", "tag123", "hei=hei;ha=ha"),
	}
	for _, t := range ts {
		fmt.Print(t.Raw())
	}
}

func TestTo_Parse(t *testing.T) {
	raws := []string{
		"to: <sip:34020000001320000001@192.168.0.1:5060>;tag=tag123",
		"to: \"\"34020000001320000001\"\" <sip:34020000001320000001@www.baidu.com>;tag=tag123",
		"t   : \"Tom\" <sip:34020000001320000001@www.baidu.com>;tag=tag123",
		"To: \"Bob\"<sip:34020000001320000001@www.baidu.com> ;tag=tag123;hai=hai;w;token=123",
		"To: \"alias\"sips:34020000001320000001@www.baidu.com ;tag=tag123;he=he;wa;token=123",
		"t: \"Alisa\"tel:34020000001320000001@www.baidu.com ;tag=tag123;hei=hei;wb;token=123",
		"To   : hack template 3 sip:34020000001320000001@www.baidu.com ;tag=tag123;ws=ws;af;token=123",
	}
	for _, raw := range raws {
		t := new(To)
		t.Parse(raw)
		if len(t.source) > 0 {
			fmt.Print("--------------------", t.Raw())
			t.SetUser("2345678")
			t.SetTag("xxxxxxxxxxxxxxx")
			if len(t.GetPort()) > 0 {
				t.SetPort("8080")
			}
			fmt.Print("||||||||||||||||||||", t.Raw())
		}

	}
}
