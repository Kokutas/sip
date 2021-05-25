package sip

import (
	"fmt"
	"regexp"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.41
//
// 20.41 User-Agent

// The User-Agent header field contains information about the UAC
// originating the request.  The semantics of this header field are
// defined in [H14.43].

// Revealing the specific software version of the user agent might allow
// the user agent to become more vulnerable to attacks against software
// that is known to contain security holes.  Implementers SHOULD make
// the User-Agent header field a configurable option.

// Example:

// 	User-Agent: Softphone Beta1.5
//
// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// User-Agent  =  "User-Agent" HCOLON server-val *(LWS server-val)
// LWS  =  [*WSP CRLF] 1*WSP ; linear whitespace

type UserAgent struct {
	field  string   // "User-Agent"
	server []string // server-val
	source string   // source string
}

func (ua *UserAgent) SetField(field string) {
	if regexp.MustCompile(`^(?i)(user-agent)$`).MatchString(field) {
		ua.field = field
	} else {
		ua.field = "User-Agent"
	}
}
func (ua *UserAgent) GetField() string {
	return ua.field
}
func (ua *UserAgent) SetServer(server ...string) {
	ua.server = server
}
func (ua *UserAgent) GetServer() []string {
	return ua.server
}
func (ua *UserAgent) GetSource() string {
	return ua.source
}
func NewUserAgent(server ...string) *UserAgent {
	return &UserAgent{
		field:  "User-Agent",
		server: server,
	}
}
func (ua *UserAgent) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(ua.field)) == 0 {
		ua.field = "User-Agent"
	}
	result.WriteString(fmt.Sprintf("%s:", ua.field))
	if ua.server != nil {
		for _, sv := range ua.server {
			result.WriteString(fmt.Sprintf(" %s", sv))
		}
	}
	result.WriteString("\r\n")
	return
}
func (ua *UserAgent) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(user-agent)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	field := regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	ua.field = field
	ua.source = raw
	ua.server = make([]string, 0)
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		rawSlice := strings.Fields(raw)
		for _, raws := range rawSlice {
			if len(strings.TrimSpace(raws)) == 0 {
				continue
			}
			raws = strings.TrimSpace(raws)
			ua.server = append(ua.server, raws)
		}
	}
}
