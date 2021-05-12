package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sip"
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
	Field  string `json:"field"`
	Server string `json:"server-val"`
}

func CreateUserAgent() sip.Sip {
	return &UserAgent{}
}
func NewUserAgent(server string) sip.Sip {
	return &UserAgent{
		Field:  "User-Agent",
		Server: server,
	}
}

func (ua *UserAgent) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, ua) {
		return result
	}
	result += fmt.Sprintf("%v:", strings.Title(ua.Field))
	if len(strings.TrimSpace(ua.Server)) > 0 {
		result += fmt.Sprintf(" %v", ua.Server)
	}
	result += "\r\n"
	return result
}
func (ua *UserAgent) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, ua) {
		return result
	}
	data, err := json.Marshal(ua)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (ua *UserAgent) Parser(raw string) error {
	if ua == nil {
		return errors.New("user-agent caller is not allowed via be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")

	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(user-agent).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		ua.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		ua.Server = raw
	} else {
		return errors.New("server-val is not match")
	}
	return nil
}
func (ua *UserAgent) Validator() error {
	if ua == nil {
		return errors.New("user-agent caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(ua.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(user-agent)`).Match([]byte(ua.Field)) {
		return errors.New("field is not match")
	}
	if len(strings.TrimSpace(ua.Server)) == 0 {
		return errors.New("server-val is not allowed to be empty")
	}
	return nil
}
