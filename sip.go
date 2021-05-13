package sip

import (
	"regexp"
	"sip/body"
	"sip/header"
	"sip/line"
	"sip/util"
	"strings"
)

type Sip struct {
	*line.RequestLine
	*line.StatusLine
	*header.Header
	*body.Body
}

func NewSip(requestLine *line.RequestLine, statusLine *line.StatusLine, header *header.Header, body *body.Body) *Sip {
	return &Sip{RequestLine: requestLine, StatusLine: statusLine, Header: header, Body: body}
}

func (sip *Sip)Raw()(string,error){
	return "",nil
}
func (sip *Sip)String()string{
	return ""
}
func (sip *Sip)Parser(raw string) error {
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// reqeust-line regexp
	// methods regexp
	methodsRegexpStr := `^(?i)(`
	for _, v := range Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ") .*?\n$"
	requestLineRegexp := regexp.MustCompile(methodsRegexpStr)
	if requestLineRegexp.MatchString(raw) {
		sip.RequestLine = new(line.RequestLine)
		sip.RequestLine.Parser(requestLineRegexp.FindString(raw))
		raw = requestLineRegexp.ReplaceAllString(raw, "")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// status-line regexp
	// sip-schema regexp
	sipSchemaRegexpStr := `^(?i)(`
	for _, v := range Schemas {
		sipSchemaRegexpStr += v + "|"
	}
	sipSchemaRegexpStr = strings.TrimSuffix(sipSchemaRegexpStr, "|")
	sipSchemaRegexpStr += ")"
	statusLineRegexp := regexp.MustCompile(sipSchemaRegexpStr + `/\d+\.\d+ \d+ .*?\n$`)
	if statusLineRegexp.MatchString(raw) {
		sip.StatusLine = new(line.StatusLine)
		sip.StatusLine.Parser(statusLineRegexp.FindString(raw))
		raw = statusLineRegexp.ReplaceAllString(raw, "")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// header-line regexp
	// body-line regexp
	// content-length

	return nil
}

func (sip *Sip)Validator()error{
	return nil
}