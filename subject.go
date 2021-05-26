package sip

import (
	"fmt"
	"regexp"
	"strings"
)

//https://www.rfc-editor.org/rfc/rfc3261.html#section-20.36

// 20.36 Subject

//    The Subject header field provides a summary or indicates the nature
//    of the call, allowing call filtering without having to parse the
//    session description.  The session description does not have to use
//    the same subject indication as the invitation.

//    The compact form of the Subject header field is s.
// Example:

//    Subject: Need more boxes
//    s: Tech Support
//
//https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
//Subject  =  ( "Subject" / "s" ) HCOLON [TEXT-UTF8-TRIM]
type Subject struct {
	field  string //"Subject" / "s"
	text   string
	source string // source string
}

func (s *Subject) SetField(field string) {
	if regexp.MustCompile(`^(?i)(subject|s)$`).MatchString(field) {
		s.field = field
	} else {
		s.field = "Subject"
	}
}
func (s *Subject) GetField() string {
	return s.field
}
func (s *Subject) SetText(text string) {
	s.text = text
}
func (s *Subject) GetText() string {
	return s.text
}
func (s *Subject) GetSource() string {
	return s.source
}
func NewSubject(text string) *Subject {
	return &Subject{
		field: "Subject",
		text:  text,
	}
}
func (s *Subject) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(s.field)) == 0 {
		s.field = "Subject"
	}
	result.WriteString(fmt.Sprintf("%s:", s.field))
	if len(strings.TrimSpace(s.text)) > 0 {
		result.WriteString(fmt.Sprintf(" %s", s.text))
	}
	result.WriteString("\r\n")
	return
}
func (s *Subject) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(subject|s)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	field := regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	s.field = field
	s.source = raw
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		s.text = raw
	}

}
