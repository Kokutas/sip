package sip

import (
	"fmt"
	"testing"
)

func TestWarning_Raw(t *testing.T) {
	w := NewWarning(307, "uas.com", "Session parameter 'foo' not understood")
	result := w.Raw()
	fmt.Print(result.String())
}

func TestWarning_Parse(t *testing.T) {
	raws := []string{
		`Warning: 307 uas.com "Session parameter 'foo' not understood"`,
		`Warning: 307 192.168.0.1 "Session parameter 'foo' not understood"`,
		`Warning: 307 192.168.0.1:5060 "Session parameter 'foo' not understood"`,
		`Warning: 307  "Session parameter 'foo' not understood"`,
		`Warning:  192.168.0.1 "Session parameter 'foo' not understood"`,
		`Warning:  "Session parameter 'foo' not understood"`,
		`Warning:  "Session parameter 10 hm"`,
	}
	for _, raw := range raws {
		w := new(Warning)
		w.Parse(raw)
		if len(w.GetSource()) > 0 {
			fmt.Println(w.GetField(), w.GetWarnCode(), w.GetWarnAgent(), w.GetWarnText())
			result := w.Raw()
			fmt.Print(result.String())
		}

	}
}
