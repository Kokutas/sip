package sip

import "strings"

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
