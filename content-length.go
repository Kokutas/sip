package sip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//https://www.rfc-editor.org/rfc/rfc3261.html#section-20.14
//
// 20.14 Content-Length

// The Content-Length header field indicates the size of the message-
// body, in decimal number of octets, sent to the recipient.
// Applications SHOULD use this field to indicate the size of the
// message-body to be transferred, regardless of the media type of the
// entity.  If a stream-based protocol (such as TCP) is used as
// transport, the header field MUST be used.

// The size of the message-body does not include the CRLF separating
// header fields and body.  Any Content-Length greater than or equal to
// zero is a valid value.  If no body is present in a message, then the
// Content-Length header field value MUST be set to zero.

// The ability to omit Content-Length simplifies the creation of
// cgi-like scripts that dynamically generate responses.

// The compact form of the header field is l.

// Examples:

//    Content-Length: 349
//    l: 173

//https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Content-Length  =  ( "Content-Length" / "l" ) HCOLON 1*DIGIT
//
type ContentLength struct {
	field  string //"Content-Length" / "l"
	length uint
	source string // source string
}

//"Content-Length" / "l"
func (l *ContentLength) SetField(field string) {
	if regexp.MustCompile(`^(?i)(content-length|l)$`).MatchString(field) {
		l.field = field
	} else {
		l.field = "Content-Length"
	}
}
func (l *ContentLength) GetField() string {
	return l.field
}
func (l *ContentLength) SetLength(length uint) {
	l.length = length
}
func (l *ContentLength) GetLength() uint {
	return l.length
}

// source string
func (l *ContentLength) GetSource() string {
	return l.source
}
func NewContentLength(length uint) *ContentLength {
	return &ContentLength{
		field:  "Content-Length",
		length: length,
	}
}
func (l *ContentLength) Raw() (result strings.Builder) {
	// "Content-Length"
	if len(strings.TrimSpace(l.field)) == 0 {
		l.field = "Content-Length"
	}
	result.WriteString(fmt.Sprintf("%s:", l.field))
	result.WriteString(fmt.Sprintf(" %d", l.length))
	result.WriteString("\r\n")
	return
}
func (l *ContentLength) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(content-length|l)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	l.source = raw
	field := fieldRegexp.FindString(raw)
	field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	l.field = field
	// length regexp
	lengthRegexp := regexp.MustCompile(`\d+`)
	if lengthRegexp.MatchString(raw) {
		lengths := lengthRegexp.FindString(raw)
		length, _ := strconv.Atoi(lengths)
		l.length = uint(length)
	}
}
