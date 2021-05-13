package msg

import (
	"github.com/kokutas/sip"
	"github.com/kokutas/sip/body"
	"github.com/kokutas/sip/header"
	"github.com/kokutas/sip/line"
	"github.com/kokutas/sip/util"
	"regexp"
	"strings"
)

type SipMsg struct {
	*line.RequestLine
	*line.StatusLine
	*header.Header
	*body.Body
}

func NewSip(requestLine *line.RequestLine, statusLine *line.StatusLine, header *header.Header, body *body.Body) *SipMsg {
	return &SipMsg{RequestLine: requestLine, StatusLine: statusLine, Header: header, Body: body}
}

func (sm *SipMsg) Raw() (string, error) {
	return "", nil
}
func (sm *SipMsg) String() string {
	return ""
}
func (sm *SipMsg) Parser(raw string) error {
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// reqeust-line regexp
	// methods regexp
	methodsRegexpStr := `^(?i)(`
	for _, v := range sip.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ") .*?\n$"
	requestLineRegexp := regexp.MustCompile(methodsRegexpStr)
	if requestLineRegexp.MatchString(raw) {
		sm.RequestLine = new(line.RequestLine)
		sm.RequestLine.Parser(requestLineRegexp.FindString(raw))
		raw = requestLineRegexp.ReplaceAllString(raw, "")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// status-line regexp
	// sip-schema regexp
	sipSchemaRegexpStr := `^(?i)(`
	for _, v := range sip.Schemas {
		sipSchemaRegexpStr += v + "|"
	}
	sipSchemaRegexpStr = strings.TrimSuffix(sipSchemaRegexpStr, "|")
	sipSchemaRegexpStr += ")"
	statusLineRegexp := regexp.MustCompile(sipSchemaRegexpStr + `/\d+\.\d+ \d+ .*?\n$`)
	if statusLineRegexp.MatchString(raw) {
		sm.StatusLine = new(line.StatusLine)
		sm.StatusLine.Parser(statusLineRegexp.FindString(raw))
		raw = statusLineRegexp.ReplaceAllString(raw, "")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// header-line regexp
	// body-line regexp
	// content-length

	return nil
}

func (sm *SipMsg) Validator() error {
	return nil
}
