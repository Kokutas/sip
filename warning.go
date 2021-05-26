package sip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.43
//
// 20.43 Warning
// The Warning header field is used to carry additional information
// about the status of a response.  Warning header field values are sent
// with responses and contain a three-digit warning code, host name, and
// warning text.

// The "warn-text" should be in a natural language that is most likely
// to be intelligible to the human user receiving the response.  This
// decision can be based on any available knowledge, such as the
// location of the user, the Accept-Language field in a request, or the
// Content-Language field in a response.  The default language is i-
// default [21].

// The currently-defined "warn-code"s are listed below, with a
// recommended warn-text in English and a description of their meaning.
// These warnings describe failures induced by the session description.
// The first digit of warning codes beginning with "3" indicates
// warnings specific to SIP.  Warnings 300 through 329 are reserved for
// indicating problems with keywords in the session description, 330
// through 339 are warnings related to basic network services requested
// in the session description, 370 through 379 are warnings related to
// quantitative QoS parameters requested in the session description, and
// 390 through 399 are miscellaneous warnings that do not fall into one
// of the above categories.
// 300 Incompatible network protocol: One or more network protocols
// contained in the session description are not available.

// 301 Incompatible network address formats: One or more network
// address formats contained in the session description are not
// available.

// 302 Incompatible transport protocol: One or more transport
// protocols described in the session description are not
// available.

// 303 Incompatible bandwidth units: One or more bandwidth
// measurement units contained in the session description were
// not understood.

// 304 Media type not available: One or more media types contained in
// the session description are not available.

// 305 Incompatible media format: One or more media formats contained
// in the session description are not available.

// 306 Attribute not understood: One or more of the media attributes
// in the session description are not supported.

// 307 Session description parameter not understood: A parameter
// other than those listed above was not understood.

// 330 Multicast not available: The site where the user is located
// does not support multicast.

// 331 Unicast not available: The site where the user is located does
// not support unicast communication (usually due to the presence
// of a firewall).
// 370 Insufficient bandwidth: The bandwidth specified in the session
// description or defined by the media exceeds that known to be
// available.

// 399 Miscellaneous warning: The warning text can include arbitrary
// information to be presented to a human user or logged.  A
// system receiving this warning MUST NOT take any automated
// action.

// 	1xx and 2xx have been taken by HTTP/1.1.

// Additional "warn-code"s can be defined through IANA, as defined in
// Section 27.2.

// Examples:

// 	Warning: 307 isi.edu "Session parameter 'foo' not understood"
// 	Warning: 301 isi.edu "Incompatible network address type 'E.164'"

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Warning        =  "Warning" HCOLON warning-value *(COMMA warning-value)
// warning-value  =  warn-code SP warn-agent SP warn-text
// warn-code      =  3DIGIT
// warn-agent     =  hostport / pseudonym
//                   ;  the name or pseudonym of the server adding
//                   ;  the Warning header, for use in debugging
// warn-text      =  quoted-string
// pseudonym      =  token

// Warning : for use in debugging
type Warning struct {
	field     string // "Warning"
	warnCode  uint   //  warn-code = 3DIGIT
	warnAgent string // warn-agent =  hostport / pseudonym;the name or pseudonym of the server adding;the Warning header, for use in debugging;pseudonym = token
	warnText  string // warn-text  =  quoted-string
	source    string // source string
}

func (w *Warning) SetField(field string) {
	if regexp.MustCompile(`^(?i)(warning)$`).MatchString(field) {
		w.field = field
	} else {
		w.field = "Warning"
	}
}
func (w *Warning) GetField() string {
	return w.field
}
func (w *Warning) SetWarnCode(warnCode uint) {
	w.warnCode = warnCode
}
func (w *Warning) GetWarnCode() uint {
	return w.warnCode
}
func (w *Warning) SetWarnAgent(warnAgent string) {
	w.warnAgent = warnAgent
}
func (w *Warning) GetWarnAgent() string {
	return w.warnAgent
}
func (w *Warning) SetWarnText(warnText string) {
	w.warnText = warnText
}
func (w *Warning) GetWarnText() string {
	return w.warnText
}
func (w *Warning) GetSource() string {
	return w.source
}
func NewWarning(warnCode uint, warnAgent string, warnText string) *Warning {
	return &Warning{
		field:     "Warning",
		warnCode:  warnCode,
		warnAgent: warnAgent,
		warnText:  warnText,
	}
}
func (w *Warning) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(w.field)) == 0 {
		w.field = "Warning"
	}
	result.WriteString(fmt.Sprintf("%s:", w.field))
	if w.warnCode > 0 {
		result.WriteString(fmt.Sprintf(" %03d", w.warnCode))
	}
	if len(strings.TrimSpace(w.warnAgent)) > 0 {
		result.WriteString(fmt.Sprintf(" %s", w.warnAgent))
	}
	if len(strings.TrimSpace(w.warnText)) > 0 {
		result.WriteString(fmt.Sprintf(" \"%s\"", w.warnText))
	}
	result.WriteString("\r\n")
	return
}
func (w *Warning) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(warning)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	field := regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	w.field = field
	w.source = raw
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// warn code regexp
	warnCodeRegexp := regexp.MustCompile(`^\d+`)
	// warn agent regexp
	warnAgentRegexp := regexp.MustCompile(`^((\d+\.\d+\.\d+\.\d+)(:\d+)?)|(\w+\.\w+(\w+\.)?)`)

	if warnCodeRegexp.MatchString(raw) {
		codes := warnCodeRegexp.FindString(raw)
		raw = strings.TrimPrefix(raw, codes)
		raw = stringTrimPrefixAndTrimSuffix(raw, " ")
		code, _ := strconv.Atoi(codes)
		if code > 0 {
			w.warnCode = uint(code)
		}
	}
	if warnAgentRegexp.MatchString(raw) {
		w.warnAgent = warnAgentRegexp.FindString(raw)
		raw = strings.TrimPrefix(raw, warnAgentRegexp.FindString(raw))
		raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, "\"")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		w.warnText = raw
	}

}
