package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sip"
	"sip/util"
	"strconv"
	"strings"
)

// The Content-Length header field indicates the size of the message-
//  body, in decimal number of octets, sent to the recipient.
//  Applications SHOULD use this field to indicate the size of the
//  message-body to be transferred, regardless of the media type of the
//  entity. If a stream-based protocol (such as TCP) is used as
//  transport, the header field MUST be used.
//  The size of the message-body does not include the CRLF separating
//  header fields and body. Any Content-Length greater than or equal to
//  zero is a valid value. If no body is present in a message, then the
//  Content-Length header field value MUST be set to zero.
//  The ability to omit Content-Length simplifies the creation of
//  cgi-like scripts that dynamically generate responses.
//  The compact form of the header field is l.
//  Examples:
//  	Content-Length: 349
//  	l: 173
// Content-Length = ( "Content-Length" / "l" ) HCOLON 1*DIGIT
type ContentLength struct {
	Field  string `json:"field"`
	Length uint   `json:"length"`
}

func CreateContentLength() sip.Sip {
	return &ContentLength{}
}
func NewContentLength(length uint) sip.Sip {
	return &ContentLength{
		Field:  "Content-Length",
		Length: length,
	}
}
func (cl *ContentLength) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, cl) {
		return result
	}
	result += fmt.Sprintf("%v:", cl.Field)
	result += fmt.Sprintf(" %v", cl.Length)
	result += "\r\n"
	return result
}
func (cl *ContentLength) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, cl) {
		return result
	}
	data, err := json.Marshal(cl)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (cl *ContentLength) Parser(raw string) error {
	if cl == nil {
		return errors.New("content-length caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(content-length).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		cl.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
	} else {
		return errors.New("filed is not match")
	}
	lengthRegexp := regexp.MustCompile(`\d+`)
	if lengthRegexp.MatchString(raw) {
		length, err := strconv.Atoi(lengthRegexp.FindString(raw))
		if err != nil {
			return err
		}
		cl.Length = uint(length)
	} else {
		return errors.New("length is not match")
	}

	return nil
}
func (cl *ContentLength) Validator() error {
	if cl == nil {
		return errors.New("content-length caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(cl.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(content-length)`).Match([]byte(cl.Field)) {
		return errors.New("field is not match")
	}
	return nil
}
