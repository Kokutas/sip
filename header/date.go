package header

import (
	"errors"
	"fmt"
	"regexp"
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
	field   string
	sipDate time.Time
	format  string
}

func (date *Date) Field() string {
	return date.field
}

func (date *Date) SetField(field string) {
	date.field = field
}

func (date *Date) SipDate() time.Time {
	return date.sipDate
}

func (date *Date) SetSipDate(sipDate time.Time) {
	date.sipDate = sipDate
}

func (date *Date) Format() string {
	return date.format
}

func (date *Date) SetFormat(format string) {
	date.format = format
}

func NewDate(sipDate time.Time, format string) *Date {
	return &Date{
		field:   "Date",
		sipDate: sipDate,
		format:  format,
	}
}
func (date *Date) Raw() (string, error) {
	result := ""
	if err:=date.Validator();err!=nil{
		return result,err
	}
	if len(strings.TrimSpace(date.field))==0{
		date.field = "Date"
	}
	result += fmt.Sprintf("%s:", date.field)
	if len(strings.TrimSpace(date.sipDate.String())) > 0 {
		if len(strings.TrimSpace(date.format)) > 0 {
			result += fmt.Sprintf(" %s", date.sipDate.Format(date.format))
		} else {
			result += fmt.Sprintf(" %s", date.sipDate.Format("2006-01-02T15:04:05.000"))
		}
	} else {
		if len(strings.TrimSpace(date.format)) > 0 {
			result += fmt.Sprintf(" %s", time.Now().Format(date.format))
		} else {
			result += fmt.Sprintf(" %s", time.Now().Format("2006-01-02T15:04:05.000"))
		}
	}
	result += "\r\n"
	return result,nil
}
func (date *Date) String() string {
	result := ""
	if len(strings.TrimSpace(date.field)) > 0 {
		result += fmt.Sprintf("field: %s,", date.field)
	}
	if len(strings.TrimSpace(date.sipDate.String())) > 0 {
		result += fmt.Sprintf("sip-date: %s,", date.sipDate.String())
	}
	if len(strings.TrimSpace(date.format)) > 0 {
		result += fmt.Sprintf("format: %s,", date.format)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (date *Date) Parser(raw string) error {
	if date == nil {
		return errors.New("date caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(date).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		date.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(date.format)) == 0 {
		date.format = "2006-01-02T15:04:05.000"
	}
	tp, err := time.Parse(date.format, raw)
	if err != nil {
		return err
	}
	date.sipDate = tp
	return nil
}
func (date *Date) Validator() error {
	if date == nil {
		return errors.New("date caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(date.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(date)`).Match([]byte(date.field)) {
		return errors.New("field is not match")
	}
	return nil
}
