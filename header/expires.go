package header

import (
	"errors"
	"fmt"
	"regexp"
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
	field    string
	seconds uint
}

func (expires *Expires) Field() string {
	return expires.field
}

func (expires *Expires) SetField(field string) {
	expires.field = field
}

func (expires *Expires) Seconds() uint {
	return expires.seconds
}

func (expires *Expires) SetSeconds(seconds uint) {
	expires.seconds = seconds
}

func NewExpires(seconds uint) *Expires {
	return &Expires{
		field:    "Expires",
		seconds: seconds,
	}
}
func (expires *Expires) Raw() (string, error) {
	result := ""
	if err := expires.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(expires.field)) == 0 {
		expires.field = "Expires"
	}
	result += fmt.Sprintf("%s:", strings.Title(expires.field))
	result += fmt.Sprintf(" %d", expires.seconds)
	result += "\r\n"
	return result, nil
}
func (expires *Expires) String() string {
	result := ""
	if len(strings.TrimSpace(expires.field)) > 0 {
		result += fmt.Sprintf("field: %s,", expires.field)
	}
	if expires.seconds >= 0 {
		result += fmt.Sprintf("delta-seconds: %d,", expires.seconds)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (expires *Expires) Parser(raw string) error {
	if expires == nil {
		return errors.New("expires caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	fieldRegexp := regexp.MustCompile(`(?i)(expires)`)
	if fieldRegexp.MatchString(raw) {
		expires.field = fieldRegexp.FindString(raw)
	}
	secondsRegexp := regexp.MustCompile(`\d+`)
	if secondsRegexp.MatchString(raw) {
		seconds, err := strconv.Atoi(secondsRegexp.FindString(raw))
		if err != nil {
			return err
		}
		expires.seconds = uint(seconds)
	} else {
		return errors.New("delta-seconds is not match")
	}
	return nil
}
func (expires *Expires) Validator() error {
	if expires == nil {
		return errors.New("expires caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(expires.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(expires)`).Match([]byte(expires.field)) {
		return errors.New("field is not match")
	}
	return nil
}
