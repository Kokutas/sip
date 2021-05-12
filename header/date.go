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
	"time"
)

// The Date header field contains the date and time. Unlike HTTP/1.1,
//  SIP only supports the most recent RFC 1123 [20] format for dates. As
//  in [H3.3], SIP restricts the time zone in SIP-date to "GMT", while
//  RFC 1123 allows any time zone. An RFC 1123 date is case-sensitive.
//  The Date header field reflects the time when the request or response
//  is first sent.
//  The Date header field can be used by simple end systems without a
//  battery-backed clock to acquire a notion of current time.
//  However, in its GMT form, it requires clients to know their offset
//  from GMT.
//  Example:
//  Date: Sat, 13 Nov 2010 23:29:00 GMT

// Date = "Date" HCOLON SIP-date
// SIP-date = rfc1123-date
// rfc1123-date = wkday "," SP date1 SP time SP "GMT"
// date1 = 2DIGIT SP month SP 4DIGIT
//  		; day month year (e.g., 02 Jun 1982)
// time = 2DIGIT ":" 2DIGIT ":" 2DIGIT
//  		; 00:00:00 - 23:59:59
// wkday = "Mon" / "Tue" / "Wed"
//  		/ "Thu" / "Fri" / "Sat" / "Sun"
// month = "Jan" / "Feb" / "Mar" / "Apr"
//  		/ "May" / "Jun" / "Jul" / "Aug"
//  		/ "Sep" / "Oct" / "Nov" / "Dec"
type Date struct {
	Field   string    `json:"field"`
	SipDate time.Time `json:"sip-date"`
	Format  string    `json:"format"`
}

func CreateDate() sip.Sip {
	return &Date{}
}
func NewDate(sipDate time.Time, format string) sip.Sip {
	return &Date{
		Field:   "Date",
		SipDate: sipDate,
		Format:  format,
	}
}
func (date *Date) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, date) {
		return result
	}
	result += fmt.Sprintf("%v:", date.Field)
	if len(strings.TrimSpace(date.SipDate.String())) > 0 {
		if len(strings.TrimSpace(date.Format)) > 0 {
			result += fmt.Sprintf(" %v", date.SipDate.Format(date.Format))
		} else {
			result += fmt.Sprintf(" %v", date.SipDate.Format("2006-01-02T15:04:05.000"))
		}
	} else {
		if len(strings.TrimSpace(date.Format)) > 0 {
			result += fmt.Sprintf(" %v", time.Now().Format(date.Format))
		} else {
			result += fmt.Sprintf(" %v", time.Now().Format("2006-01-02T15:04:05.000"))
		}
	}
	result += "\r\n"
	return result
}
func (date *Date) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, date) {
		return result
	}
	data, err := json.Marshal(date)
	if err != nil {
		return result
	}
	result += fmt.Sprintf("%s", data)
	return result
}
func (date *Date) Parser(raw string) error {
	if date == nil {
		return errors.New("date caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(date).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		date.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(date.Format)) == 0 {
		date.Format = "2006-01-02T15:04:05.000"
	}
	tp, err := time.Parse(date.Format, raw)
	if err != nil {
		return err
	}
	date.SipDate = tp
	return nil
}
func (date *Date) Validator() error {
	if date == nil {
		return errors.New("date caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(date.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(date)`).Match([]byte(date.Field)) {
		return errors.New("field is not match")
	}
	return nil
}
