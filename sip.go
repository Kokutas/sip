package sip

import (
	"strings"
)

const (
	sip  = "sip"
	sips = "sips"
	tel  = "tel"
)

var schemas = map[string]string{
	sip:  sip,
	sips: sips,
	tel:  tel,
}

type SipLayer interface {
	Raw() string
	Parse()
}

func stringTrimPrefixAndTrimSuffix(source string, sub string) string {
	for strings.HasPrefix(source, sub) || strings.HasSuffix(source, sub) {
		source = strings.TrimPrefix(source, sub)
		source = strings.TrimSuffix(source, sub)
	}
	return source
}

// type generic struct {
// 	index int // order
// 	kv    map[int]map[string]interface{}
// 	gk    sync.RWMutex
// }

// func (g *generic) store(k string, v interface{}) {
// 	g.gk.Lock()
// 	defer g.gk.Unlock()
// 	g.index++
// 	g.kv[g.index] = map[string]interface{}{k: v}
// }
