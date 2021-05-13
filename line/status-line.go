package line

import (
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
	*sip.SipVersion
	statusCode      int
	reasonPhrase    string
}
func (sl *StatusLine) StatusCode() int {
	return sl.statusCode
}

func (sl *StatusLine) SetStatusCode(statusCode int) {
	sl.statusCode = statusCode
}

func (sl *StatusLine) ReasonPhrase() string {
	return sl.reasonPhrase
}

func (sl *StatusLine) SetReasonPhrase(reasonPhrase string) {
	sl.reasonPhrase = reasonPhrase
}
func NewStatusLine(sipVersion *sip.SipVersion, statusCode int, reasonPhrase string) *StatusLine {
	return &StatusLine{SipVersion: sipVersion, statusCode: statusCode, reasonPhrase: reasonPhrase}
}

func (sl *StatusLine) Raw() (string,error) {
	result := ""
	if err:=sl.Validator();err!=nil {
		return result,err
	}
	if sl.SipVersion != nil {
		res,err:=sl.SipVersion.Raw()
		if err!=nil{
			return "", err
		}
		result +=res
	}
	if sl.statusCode > 0 && len(strings.TrimSpace(sl.reasonPhrase)) > 0 {
		result += fmt.Sprintf(" %v %v", sl.statusCode, sl.reasonPhrase)
	}
	result += "\r\n"
	return result,nil
}
func (sl *StatusLine) String() string {
	result:=""
	if sl.SipVersion!=nil{
		result+=fmt.Sprintf("%s,",sl.SipVersion.String())
	}
	if sl.statusCode>0{
		result+=fmt.Sprintf("status-code: %v,",sl.statusCode)
	}
	if len(strings.TrimSpace(sl.reasonPhrase))>0{
		result+=fmt.Sprintf("reason-phrase: %s,",sl.reasonPhrase)
	}
	result = strings.TrimSuffix(result,",")
	return result
}
func (sl *StatusLine) Parser(raw string) error {
	if sl == nil {
		return errors.New("statusLine caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw," ")
	raw = strings.TrimPrefix(raw," ")
	raw = strings.TrimSuffix(raw," ")
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
		raw = strings.TrimLeft(raw, " ")
		raw = strings.TrimRight(raw," ")
		raw = strings.TrimPrefix(raw," ")
		raw = strings.TrimSuffix(raw," ")
		sl.SipVersion = new(sip.SipVersion)
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
		sl.statusCode = code
		raw = statusCodeRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimSuffix(raw, "\r")
		raw = strings.TrimSuffix(raw, "\n")
		raw = strings.TrimPrefix(raw, "\r")
		raw = strings.TrimPrefix(raw, "\n")
		raw = strings.TrimLeft(raw, " ")
		raw = strings.TrimRight(raw," ")
		raw = strings.TrimPrefix(raw," ")
		raw = strings.TrimSuffix(raw," ")
	}
	if len(strings.TrimSpace(raw)) > 0 {
		sl.reasonPhrase = raw
	}
	return nil
}
func (sl *StatusLine) Validator() error {
	if sl == nil {
		return errors.New("statusLine caller is not allowed to be nil")
	}
	if err := sl.SipVersion.Validator(); err != nil {
		return err
	}
	if sl.statusCode <= 0 {
		return errors.New("invalid statusCode")
	}
	if len(strings.TrimSpace(sl.reasonPhrase)) == 0 {
		return errors.New("reasonPhrase is not allowed to be empty")
	}
	return nil
}
