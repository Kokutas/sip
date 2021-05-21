package sip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.19
//
// 20.19 Expires
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
//    Expires: 5

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Expires     =  "Expires" HCOLON delta-seconds

type Expires struct {
	field  string // "Expires"
	expire uint32 // delta-seconds
	source string // to header line source string
}

func (expires *Expires) SetField(field string) {
	if regexp.MustCompile(`^(?i)(expires)$`).MatchString(field) {
		expires.field = strings.Title(field)
	}
}
func (expires *Expires) GetField() string {
	return expires.field
}
func (expires *Expires) SetExpire(expire uint32) {
	expires.expire = expire
}
func (expires *Expires) GetExpire() uint32 {
	return expires.expire
}
func (expires *Expires) SetSource(source string) {
	expires.source = source
}
func (expires *Expires) GetSource() string {
	return expires.source
}
func NewExpires(expire uint32) *Expires {
	return &Expires{
		expire: expire,
	}
}
func (expires *Expires) Raw() string {
	result := ""
	if len(strings.TrimSpace(expires.field)) > 0 {
		result += fmt.Sprintf("%s:", expires.field)
	} else {
		result += fmt.Sprintf("%s:", strings.Title("expires"))
	}
	result += fmt.Sprintf(" %d", expires.expire)
	result += "\r\n"
	return result
}
func (expires *Expires) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// expires field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(expires)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	expires.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	expires.source = raw
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// delta-seconds regexp
	deltaSecondsRegexp := regexp.MustCompile(`\d+`)
	if deltaSecondsRegexp.MatchString(raw) {
		seconds := deltaSecondsRegexp.FindString(raw)
		if len(seconds) > 0 {
			second, _ := strconv.Atoi(seconds)
			expires.expire = uint32(second)
		}
	}
}
