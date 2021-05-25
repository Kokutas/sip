package sip

import (
	"fmt"
	"sync"
	"testing"
)

func TestTo_Raw(t *testing.T) {
	p := sync.Map{}
	p.Store("hei", "hei")
	p.Store("ha", nil)
	ts := []*To{
		NewTo("", "<", "sip", "34020000001320000001", "192.168.0.1", 5060, "tag123", sync.Map{}),
		NewTo("34020000001320000001", "'", "sip", "34020000001320000001", "www.baidu.com", 0, "tag123", sync.Map{}),
		NewTo("tom", "\"", "sip", "34020000001320000001", "www.baidu.com", 0, "tag123", sync.Map{}),
		NewTo("alisa", "", "sip", "34020000001320000001", "www.baidu.com", 0, "tag123", p),
	}
	for _, t := range ts {
		result := t.Raw()
		fmt.Print(result.String())
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
			t.SetUser("2345678")
			t.SetTag("xxxxxxxxxxxxxxx")
			if t.GetPort() > 0 {
				t.SetPort(8080)
			}
			result := t.Raw()
			fmt.Print("||||||||||||||||||||", result.String())
		}

	}
}
