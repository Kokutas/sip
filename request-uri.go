package sip

import (
	"errors"
	"regexp"
	"strings"
)

type RequestUri struct {
	*SipUri      // SIP-URI / SIPS-URI
	*AbsoluteUri // absoluteURI
}

func CreateRequestUri() Sip {
	return &RequestUri{}
}

func NewRequestUri(uri *SipUri, absolute *AbsoluteUri) Sip {
	return &RequestUri{uri, absolute}
}

func (ru *RequestUri) Raw() string {
	result := ""
	if ru == nil {
		return result
	}
	switch {
	case ru.SipUri != nil:
		result += ru.SipUri.Raw()
	case ru.AbsoluteUri != nil:
		result += ru.AbsoluteUri.Raw()
	}
	return result
}
func (ru *RequestUri) JsonString() string {
	result := ""
	return result
}
func (ru *RequestUri) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimSuffix(raw, "\r")
	raw = strings.TrimSuffix(raw, "\n")
	raw = strings.TrimPrefix(raw, "\r")
	raw = strings.TrimPrefix(raw, "\n")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	if ru == nil {
		return errors.New("RequestUri caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}

	// absolute-uri regexp
	absoluteUriRegexp := regexp.MustCompile(`.*/{1,}`)
	if absoluteUriRegexp.MatchString(raw) {
		ru.AbsoluteUri = CreateAbsoluteUri().(*AbsoluteUri)
		if err := ru.AbsoluteUri.Parser(raw); err != nil {
			return err
		}
	} else {
		// sip/sips-uri regexp
		if len(strings.TrimSpace(raw)) > 0 {
			ru.SipUri = CreateSipUri().(*SipUri)
			if err := ru.SipUri.Parser(raw); err != nil {
				return err
			}
		}
	}
	return nil
}
func (ru *RequestUri) Validator() error {
	return nil
}
