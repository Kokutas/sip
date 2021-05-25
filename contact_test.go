package sip

import (
	"fmt"
	"sync"
	"testing"
)

func TestContact_Raw(t *testing.T) {
	ms := []*Contact{
		NewContact("", "<", "sip", "34020000001320000001", "192.168.0.1", 5060, "0.7", 3600, sync.Map{}),
		NewContact("", "", "tel", "34020000001320000001", "192.168.0.1", 5060, "0.7", 3600, sync.Map{}),
		NewContact("display name", "", "sips", "34020000001320000001", "192.168.0.1", 5060, "", 3600, sync.Map{}),
	}
	for _, m := range ms {
		result := m.Raw()
		fmt.Print(result.String())

	}
}

func TestContact_Parse(t *testing.T) {
	raws := []string{
		"Contact: <sip:34020000001320000001@192.168.0.1:5060>;expires=3600;q=0.7",
		"Contact: tel:34020000001320000001@192.168.0.1:5060;q=0.7;expires=3600",
		"m  : \"display name\" sips:34020000001320000001@192.168.0.1:5060;",
	}
	for _, raw := range raws {
		m := new(Contact)
		m.Parse(raw)
		if len(m.GetSource()) > 0 {
			result := m.Raw()
			fmt.Print(result.String())
			fmt.Println(m.GetField(), m.GetName(), m.GetUser(), m.GetHost(), m.GetPort(), m.GetExpires(), m.GetQ())
		}
	}
}
