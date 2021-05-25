package sip

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.17
//
// 20.17 Date
// The Date header field contains the date and time.  Unlike HTTP/1.1,
// SIP only supports the most recent RFC 1123 [20] format for dates.  As
// in [H3.3], SIP restricts the time zone in SIP-date to "GMT", while
// RFC 1123 allows any time zone.  An RFC 1123 date is case-sensitive.

// The Date header field reflects the time when the request or response
// is first sent.
// The Date header field can be used by simple end systems without a
// battery-backed clock to acquire a notion of current time.
// However, in its GMT form, it requires clients to know their offset
// from GMT.
// Example:
//    Date: Sat, 13 Nov 2010 23:29:00 GMT
//
// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Date          =  "Date" HCOLON SIP-date
// SIP-date      =  rfc1123-date
// rfc1123-date  =  wkday "," SP date1 SP time SP "GMT"
// date1         =  2DIGIT SP month SP 4DIGIT
//                  ; day month year (e.g., 02 Jun 1982)
// time          =  2DIGIT ":" 2DIGIT ":" 2DIGIT
//                  ; 00:00:00 - 23:59:59
// wkday         =  "Mon" / "Tue" / "Wed"
//                  / "Thu" / "Fri" / "Sat" / "Sun"
// month         =  "Jan" / "Feb" / "Mar" / "Apr"
//                  / "May" / "Jun" / "Jul" / "Aug"
//                  / "Sep" / "Oct" / "Nov" / "Dec"
type Date struct {
	field      string //  "Date"
	timeFormat string // default: yyyy-MM-dd'T'HH:mm:ss.SSS
	sipDate    time.Time
	source     string // source string
}

func (date *Date) SetField(field string) {
	if regexp.MustCompile(`^(?i)(date)$`).MatchString(field) {
		date.field = strings.Title(field)
	} else {
		date.field = "Date"
	}
}
func (date *Date) GetField() string {
	return date.field
}
func (date *Date) SetTimeFormat(timeFormat string) {
	date.timeFormat = timeFormat
}
func (date *Date) GetTimeFormat() string {
	return date.timeFormat
}
func (date *Date) SetSipDate(sipDate time.Time) {
	date.sipDate = sipDate
}
func (date *Date) GetSipDate() time.Time {
	return date.sipDate
}
func (date *Date) GetSource() string {
	return date.source
}

func NewDate(timeFormat string, sipDate time.Time) *Date {
	return &Date{
		field:      "Date",
		timeFormat: timeFormat,
		sipDate:    sipDate,
	}
}
func (date *Date) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(date.field)) == 0 {
		date.field = "Date"
	}
	result.WriteString(fmt.Sprintf("%s:", strings.Title(date.field)))
	if len(strings.TrimSpace(date.timeFormat)) == 0 {
		date.timeFormat = "2006-01-02T15:04:05.000"
	}
	result.WriteString(fmt.Sprintf(" %s", date.sipDate.Format(date.timeFormat)))
	result.WriteString("\r\n")
	return
}
func (date *Date) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// date field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(date)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	date.source = raw
	field := fieldRegexp.FindString(raw)
	raw = strings.TrimPrefix(raw, field)
	field = strings.ReplaceAll(field, ":", "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	date.field = field
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, ":")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		if sipDate, err := time.Parse(time.ANSIC, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.ANSIC
		} else if sipDate, err := time.Parse(time.UnixDate, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.UnixDate
		} else if sipDate, err := time.Parse(time.RubyDate, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RubyDate
		} else if sipDate, err := time.Parse(time.RFC822, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RFC822
		} else if sipDate, err := time.Parse(time.RFC822Z, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RFC822Z
		} else if sipDate, err := time.Parse(time.RFC850, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RFC850
		} else if sipDate, err := time.Parse(time.RFC1123, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RFC1123
		} else if sipDate, err := time.Parse(time.RFC1123Z, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RFC1123Z
		} else if sipDate, err := time.Parse(time.RFC3339, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RFC3339
		} else if sipDate, err := time.Parse(time.RFC3339Nano, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.RFC3339Nano
		} else if sipDate, err := time.Parse(time.Kitchen, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.Kitchen
		} else if sipDate, err := time.Parse(time.Stamp, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.Stamp
		} else if sipDate, err := time.Parse(time.StampMilli, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.StampMilli
		} else if sipDate, err := time.Parse(time.StampMicro, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.StampMicro
		} else if sipDate, err := time.Parse(time.StampNano, raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = time.StampNano
		} else if sipDate, err := time.Parse("2006-01-02T15:04:05.000", raw); err == nil {
			date.sipDate = sipDate
			date.timeFormat = "2006-01-02T15:04:05.000"
		} else {
			// date.sipDate = time.Now()
			date.timeFormat = "unknown"
		}
	}
}
