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
	field   string //"Content-Length" / "l"
	length  uint
	isOrder bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order   chan string // It is convenient to record the order of the original parameter fields when parsing
	source  string      // source string
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
		field:   "Content-Length",
		length:  length,
		isOrder: false,
	}
}
func (l *ContentLength) Raw() string {
	result := ""
	if l.isOrder {
		for data := range l.order {
			result += data
		}
		l.isOrder = false
		result += "\r\n"
		return result
	}

	// "Content-Length"
	if len(strings.TrimSpace(l.field)) == 0 {
		l.field = "Content-Length"
	}
	result += fmt.Sprintf("%s:", l.field)
	result += fmt.Sprintf(" %d", l.length)
	result += "\r\n"
	return result
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
	// content-length order
	l.contentlengthOrder(raw)
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
func (l *ContentLength) contentlengthOrder(raw string) {
	l.order = make(chan string, 1024)
	l.isOrder = true
	defer close(l.order)
	l.order <- raw
}
