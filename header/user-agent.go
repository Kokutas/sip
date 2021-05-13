package header

import (
	"errors"
	"fmt"
	"regexp"
	"sip/util"
	"strings"
)

// The User-Agent header field contains information about the UAC
//  originating the request. The semantics of this header field are
//  defined in [H14.43].
//  Revealing the specific software version of the user agent might allow
//  the user agent to become more vulnerable to attacks against software
//  that is known to contain security holes. Implementers SHOULD make
//  the User-Agent header field a configurable option.
//  Example:
//  User-Agent: Softphone Beta1.5
// User-Agent = "User-Agent" HCOLON server-val *(LWS server-val)

type UserAgent struct {
	field  string
	server string
}

func (ua *UserAgent) Field() string {
	return ua.field
}

func (ua *UserAgent) SetField(field string) {
	ua.field = field
}

func (ua *UserAgent) Server() string {
	return ua.server
}

func (ua *UserAgent) SetServer(server string) {
	ua.server = server
}

func NewUserAgent(server string) *UserAgent {
	return &UserAgent{
		field:  "User-Agent",
		server: server,
	}
}

func (ua *UserAgent) Raw() (string, error) {
	result := ""
	if err := ua.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(ua.field)) == 0 {
		ua.field = "User-Agent"
	}
	result += fmt.Sprintf("%s:", strings.Title(ua.field))
	if len(strings.TrimSpace(ua.server)) > 0 {
		result += fmt.Sprintf(" %s", ua.server)
	}
	result += "\r\n"
	return result, nil
}
func (ua *UserAgent) String() string {
	result := ""
	if len(strings.TrimSpace(ua.field)) > 0 {
		result += fmt.Sprintf("field: %s,", ua.field)
	}
	if len(strings.TrimSpace(ua.server)) > 0 {
		result += fmt.Sprintf("server-val: %s,", ua.server)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (ua *UserAgent) Parser(raw string) error {
	if ua == nil {
		return errors.New("user-agent caller is not allowed via be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(user-agent).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		ua.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		ua.server = raw
	} else {
		return errors.New("server-val is not match")
	}
	return nil
}
func (ua *UserAgent) Validator() error {
	if ua == nil {
		return errors.New("user-agent caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(ua.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(user-agent)`).Match([]byte(ua.field)) {
		return errors.New("field is not match")
	}
	if len(strings.TrimSpace(ua.server)) == 0 {
		return errors.New("server-val is not allowed to be empty")
	}
	return nil
}
