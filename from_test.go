package sip

import (
	"fmt"
	"sync"
	"testing"
)

func TestFrom_Raw(t *testing.T) {
	p := sync.Map{}
	p.Store("hei", "hei")
	p.Store("ha", nil)
	fs := []*From{
		NewFrom("", "<", "sip", "34020000001320000001", "192.168.0.1", 5060, "tag123", sync.Map{}),
		NewFrom("34020000001320000001", "'", "sip", "34020000001320000001", "www.baidu.com", 0, "tag123", sync.Map{}),
		NewFrom("tom", "\"", "sip", "34020000001320000001", "www.baidu.com", 0, "tag123", sync.Map{}),
		NewFrom("alisa", "", "sip", "34020000001320000001", "www.baidu.com", 0, "tag123", p),
	}
	for _, f := range fs {
		result := f.Raw()
		fmt.Print(result.String())
	}
}

func TestFrom_Parse(t *testing.T) {
	raws := []string{
		"From: <sip:34020000001320000001@192.168.0.1:5060>;tag=tag123",
		"From: \"\"34020000001320000001\"\" <sip:34020000001320000001@www.baidu.com>;tag=tag123",
		"From   : \"Tom\" <sip:34020000001320000001@www.baidu.com>;tag=tag123",
		"f: \"Bob\"<sip:34020000001320000001@www.baidu.com> ;tag=tag123;hai=hai;w;token=123",
		"f          : \"alias\"sips:34020000001320000001@www.baidu.com ;tag=tag123;he=he;wa;token=123",
		"From: \"Alisa\"tel:34020000001320000001@www.baidu.com ;tag=tag123;hei=hei;wb;token=123",
		"From: hack template 3 sip:34020000001320000001@www.baidu.com ;tag=tag123;ws=ws;af;token=123",
	}
	for _, raw := range raws {
		f := new(From)
		f.Parse(raw)
		if len(f.source) > 0 {
			f.SetUser("2345678")
			f.SetTag("xxxxxxxxxxxxxxx")
			if f.GetPort() > 0 {
				f.SetPort(8080)
			}
			// parameter := f.GetParameter()
			// parameter.Range(func(key, value interface{}) bool {
			// 	fmt.Println(key, value)
			// 	return true
			// })
			result := f.Raw()
			fmt.Print("--------------------", result.String())
		}
	}
}
