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

// The Expires header field gives the relative time after which the
// message (or content) expires.

// The precise meaning of this is method dependent.

// The expiration time in an INVITE does not affect the duration of the
// actual session that may result from the invitation.  Session
// description protocols may offer the ability to express time limits on
// the session duration, however.

// The value of this field is an integral number of seconds (in decimal)
// between 0 and (2**32)-1, measured from the receipt of the request.
// Example:

//       Expires: 5
//
//Expires     =  "Expires" HCOLON delta-seconds
type Expires struct {
	Field    string `json:"field"`
	Secondes uint   `json:"delta-seconds"`
}

func CreateExpires() sip.Sip {
	return &Expires{}
}
func NewExpires(secondes uint) sip.Sip {
	return &Expires{
		Field:    "Expires",
		Secondes: secondes,
	}
}
func (expires *Expires) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, expires) {
		return result
	}
	result += fmt.Sprintf("%v:", strings.Title(expires.Field))
	result += fmt.Sprintf(" %v", expires.Secondes)
	result += "\r\n"
	return result
}
func (expires *Expires) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, expires) {
		return result
	}
	data, err := json.Marshal(expires)
	if err != nil {
		return result
	}
	result += fmt.Sprintf("%s", data)
	return result
}
func (expires *Expires) Parser(raw string) error {
	if expires == nil {
		return errors.New("expires caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")

	fieldRegexp := regexp.MustCompile(`(?i)(expires)`)
	if fieldRegexp.MatchString(raw) {
		expires.Field = fieldRegexp.FindString(raw)
	}

	secondsRegexp := regexp.MustCompile(`\d+`)
	if secondsRegexp.MatchString(raw) {
		seconds, err := strconv.Atoi(secondsRegexp.FindString(raw))
		if err != nil {
			return err
		}
		expires.Secondes = uint(seconds)
	} else {
		return errors.New("delta-seconds is not match")
	}
	return nil
}
func (expires *Expires) Validator() error {
	if expires == nil {
		return errors.New("expires caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(expires.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(expires)`).Match([]byte(expires.Field)) {
		return errors.New("field is not match")
	}
	return nil
}
