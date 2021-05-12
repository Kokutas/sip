package line

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sip"
	"strings"
)

// Request-Line   =  Method SP Request-URI SP SIP-Version CRLF
// Request-URI    =  SIP-URI / SIPS-URI / absoluteURI
// absoluteURI    =  scheme ":" ( hier-part / opaque-part )
// hier-part      =  ( net-path / abs-path ) [ "?" query ]
// net-path       =  "//" authority [ abs-path ]
// abs-path       =  "/" path-segments
// opaque-part    =  uric-no-slash *uric
// uric           =  reserved / unreserved / escaped
// uric-no-slash  =  unreserved / escaped / ";" / "?" / ":" / "@"
//                   / "&" / "=" / "+" / "$" / ","
// path-segments  =  segment *( "/" segment )
// segment        =  *pchar *( ";" param )
// param          =  *pchar
// pchar          =  unreserved / escaped /
//                   ":" / "@" / "&" / "=" / "+" / "$" / ","
// scheme         =  ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
// authority      =  srvr / reg-name
// srvr           =  [ [ userinfo "@" ] hostport ]
// reg-name       =  1*( unreserved / escaped / "$" / ","
//                   / ";" / ":" / "@" / "&" / "=" / "+" )
// query          =  *uric
// SIP-Version    =  "SIP" "/" 1*DIGIT "." 1*DIGIT
type RequestLine struct {
	Method          string `json:"method"`
	*sip.RequestUri `json:"request-uri"`
	*sip.SipVersion `json:"sip-version"`
}

func CreateRequestLine() sip.Sip {
	return &RequestLine{}
}
func NewRequestLine(method string, uri *sip.RequestUri, version *sip.SipVersion) sip.Sip {
	return &RequestLine{
		Method:     method,
		RequestUri: uri,
		SipVersion: version,
	}
}
func (rl *RequestLine) Raw() string {
	result := ""
	if len(strings.TrimSpace(rl.Method)) > 0 {
		result += fmt.Sprintf("%v", strings.ToUpper(rl.Method))
	}
	if rl.RequestUri != nil {
		result += fmt.Sprintf(" %v", rl.RequestUri.Raw())
	}
	if rl.SipVersion != nil {
		result += fmt.Sprintf(" %v", rl.SipVersion.Raw())
	}
	result += "\r\n"
	return result
}
func (rl *RequestLine) JsonString() string {
	result := ""
	if rl == nil {
		return result
	}
	if data, err := json.Marshal(rl); err != nil {
		return result
	} else {
		result = fmt.Sprintf("%s", data)
	}
	return result
}
func (rl *RequestLine) Parser(raw string) error {
	if rl == nil {
		return errors.New("requestLine caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimSuffix(raw, "\r")
	raw = strings.TrimSuffix(raw, "\n")
	raw = strings.TrimPrefix(raw, "\r")
	raw = strings.TrimPrefix(raw, "\n")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range sip.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")"
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if methodsRegexp.MatchString(raw) {
		rl.Method = methodsRegexp.FindString(raw)
		raw = methodsRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
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
		rl.SipVersion = sip.CreateSipVersion().(*sip.SipVersion)
		if err := rl.SipVersion.Parser(version); err != nil {
			return err
		}
	}
	rl.RequestUri = sip.CreateRequestUri().(*sip.RequestUri)
	return rl.RequestUri.Parser(raw)
}
func (rl *RequestLine) Validator() error {
	if rl == nil {
		return errors.New("RequestLine caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(rl.Method)) == 0 {
		return errors.New("method is not allowed to be empty")
	}
	if _, ok := sip.Methods[strings.ToUpper(rl.Method)]; !ok {
		return errors.New("method is not support")
	}
	if err := rl.SipUri.Validator(); err != nil {
		return err
	}

	return rl.SipVersion.Validator()
}
