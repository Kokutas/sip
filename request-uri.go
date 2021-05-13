package sip

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type RequestUri struct {
	*SipUri      // SIP-URI / SIPS-URI
	*AbsoluteUri // absoluteURI
}

func NewRequestUri(sipUri *SipUri, absoluteUri *AbsoluteUri) *RequestUri {
	return &RequestUri{SipUri: sipUri, AbsoluteUri: absoluteUri}
}
func (ru *RequestUri) Raw() (string, error) {
	result := ""
	if err := ru.Validator(); err != nil {
		return result, err
	}
	switch {
	case ru.SipUri != nil:
		res, err := ru.SipUri.Raw()
		if err != nil {
			return "", err
		}
		result += res
	case ru.AbsoluteUri != nil:
		res, err := ru.AbsoluteUri.Raw()
		if err != nil {
			return "", err
		}
		result += res
	}
	return result, nil
}
func (ru *RequestUri) String() string {
	result := ""
	if ru.SipUri != nil {
		result += fmt.Sprintf("%s,", ru.SipUri.String())
	}
	if ru.AbsoluteUri != nil {
		result += fmt.Sprintf("%s,", ru.AbsoluteUri.String())
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (ru *RequestUri) Parser(raw string) error {
	if ru == nil {
		return errors.New("requestUri caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw," ")
	raw = strings.TrimPrefix(raw," ")
	raw = strings.TrimSuffix(raw," ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}

	// absolute-uri regexp
	absoluteUriRegexp := regexp.MustCompile(`.*/{1,}`)
	if absoluteUriRegexp.MatchString(raw) {
		ru.AbsoluteUri = new(AbsoluteUri)
		if err := ru.AbsoluteUri.Parser(raw); err != nil {
			return err
		}
	} else {
		// sip/sips-uri regexp
		if len(strings.TrimSpace(raw)) > 0 {
			ru.SipUri = new(SipUri)
			if err := ru.SipUri.Parser(raw); err != nil {
				return err
			}
		}
	}
	return nil
}
func (ru *RequestUri) Validator() error {
	if ru == nil {
		return errors.New("requestUri caller is not allowed to be nil")
	}
	if ru.SipUri == nil && ru.AbsoluteUri == nil {
		return errors.New("sipUri or absoluteUri must has one")
	}
	if ru.SipUri != nil {
		if err := ru.SipUri.Validator(); err != nil {
			return err
		}
	}
	if ru.AbsoluteUri != nil {
		if err := ru.AbsoluteUri.Validator(); err != nil {
			return err
		}
	}
	return nil
}
