package line

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sip"
	"strconv"
	"strings"
)

// Status-Line     =  SIP-Version SP Status-Code SP Reason-Phrase CRLF
// Status-Code     =  Informational
//                /   Redirection
//                /   Success
//                /   Client-Error
//                /   Server-Error
//                /   Global-Failure
//                /   extension-code
// extension-code  =  3DIGIT
// Reason-Phrase   =  *(reserved / unreserved / escaped
//                    / UTF8-NONASCII / UTF8-CONT / SP / HTAB)

// Informational  =  "100"  ;  Trying
//               /   "180"  ;  Ringing
//               /   "181"  ;  Call Is Being Forwarded
//               /   "182"  ;  Queued
//               /   "183"  ;  Session Progress
// Success  =  "200"  ;  OK

// Redirection  =  "300"  ;  Multiple Choices
//             /   "301"  ;  Moved Permanently
//             /   "302"  ;  Moved Temporarily
//             /   "305"  ;  Use Proxy
//             /   "380"  ;  Alternative Service

// Client-Error  =  "400"  ;  Bad Request
//              /   "401"  ;  Unauthorized
//              /   "402"  ;  Payment Required
//              /   "403"  ;  Forbidden
//              /   "404"  ;  Not Found
//              /   "405"  ;  Method Not Allowed
//              /   "406"  ;  Not Acceptable
//              /   "407"  ;  Proxy Authentication Required
//              /   "408"  ;  Request Timeout
//              /   "410"  ;  Gone
//              /   "413"  ;  Request Entity Too Large
//              /   "414"  ;  Request-URI Too Large
//              /   "415"  ;  Unsupported Media Type
//              /   "416"  ;  Unsupported URI Scheme
//              /   "420"  ;  Bad Extension
//              /   "421"  ;  Extension Required
//              /   "423"  ;  Interval Too Brief
//              /   "480"  ;  Temporarily not available
//              /   "481"  ;  Call Leg/Transaction Does Not Exist
//              /   "482"  ;  Loop Detected
//              /   "483"  ;  Too Many Hops
//              /   "484"  ;  Address Incomplete
//              /   "485"  ;  Ambiguous
//              /   "486"  ;  Busy Here
//              /   "487"  ;  Request Terminated
//              /   "488"  ;  Not Acceptable Here
//              /   "491"  ;  Request Pending
//              /   "493"  ;  Undecipherable

// Server-Error  =  "500"  ;  Internal Server Error
//              /   "501"  ;  Not Implemented
//              /   "502"  ;  Bad Gateway
//              /   "503"  ;  Service Unavailable
//              /   "504"  ;  Server Time-out
//              /   "505"  ;  SIP Version not supported
//              /   "513"  ;  Message Too Large
// Global-Failure  =  "600"  ;  Busy Everywhere
// 			      /   "603"  ;  Decline
// 			      /   "604"  ;  Does not exist anywhere
// 			      /   "606"  ;  Not Acceptable

type StatusLine struct {
	*sip.SipVersion `json:"SIP-Version"`
	StatusCode      int    `json:"Status-Code"`
	ReasonPhrase    string `json:"Reason-Phrase"`
}

func CreateStatusLine() sip.Sip {
	return &StatusLine{}
}
func NewStatusLine(version *sip.SipVersion, code int, reason string) sip.Sip {
	return &StatusLine{version, code, reason}
}
func (sl *StatusLine) Raw() string {
	result := ""
	if sl == nil {
		return result
	}
	if sl.SipVersion != nil {
		result += sl.SipVersion.Raw()
	}
	if sl.StatusCode > 0 && len(strings.TrimSpace(sl.ReasonPhrase)) > 0 {
		result += fmt.Sprintf(" %v %v", sl.StatusCode, sl.ReasonPhrase)
	}
	result += "\r\n"
	return result
}
func (sl *StatusLine) JsonString() string {
	result := ""
	if sl == nil {
		return result
	}
	data, err := json.Marshal(sl)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (sl *StatusLine) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimSuffix(raw, "\r")
	raw = strings.TrimSuffix(raw, "\n")
	raw = strings.TrimPrefix(raw, "\r")
	raw = strings.TrimPrefix(raw, "\n")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	if sl == nil {
		return errors.New("StatusLine caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	// sip-schema regexp
	sipSchemaRegexpStr := `(?i)(`
	for _, v := range sip.Schemas {
		sipSchemaRegexpStr += v + "|"
	}
	sipSchemaRegexpStr = strings.TrimSuffix(sipSchemaRegexpStr, "|")
	sipSchemaRegexpStr += ")"
	sipVersionRegexp := regexp.MustCompile(sipSchemaRegexpStr + `/\d+\.\d+`)
	if sipVersionRegexp.MatchString(raw) {
		version := sipVersionRegexp.FindString(raw)
		raw = sipVersionRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
		sl.SipVersion = sip.CreateSipVersion().(*sip.SipVersion)
		if err := sl.SipVersion.Parser(version); err != nil {
			return err
		}
	}
	statusCodeRegexp := regexp.MustCompile(`\d+`)
	if statusCodeRegexp.MatchString(raw) {
		code, err := strconv.Atoi(statusCodeRegexp.FindString(raw))
		if err != nil {
			return err
		}
		sl.StatusCode = code
		raw = statusCodeRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimSuffix(raw, "\r")
		raw = strings.TrimSuffix(raw, "\n")
		raw = strings.TrimPrefix(raw, "\r")
		raw = strings.TrimPrefix(raw, "\n")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	if len(strings.TrimSpace(raw)) > 0 {
		sl.ReasonPhrase = raw
	}
	return nil
}
func (sl *StatusLine) Validator() error {
	if sl == nil {
		return errors.New("StatusLine caller is not allowed to be nil")
	}
	if err := sl.SipVersion.Validator(); err != nil {
		return err
	}
	if sl.StatusCode <= 0 {
		return errors.New("invalid statusCode")
	}
	if len(strings.TrimSpace(sl.ReasonPhrase)) == 0 {
		return errors.New("reasonPhrase is not allowed to be empty")
	}
	return nil
}
