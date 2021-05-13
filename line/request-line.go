package line

import (
	"errors"
	"fmt"
	"regexp"
	"github.com/kokutas/sip"
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
	method string
	*sip.RequestUri
	*sip.SipVersion
}

func (rl *RequestLine) Method() string {
	return rl.method
}
func (rl *RequestLine) SetMethod(method string) {
	rl.method = method
}
func NewRequestLine(method string, requestUri *sip.RequestUri, sipVersion *sip.SipVersion) *RequestLine {
	return &RequestLine{method: method, RequestUri: requestUri, SipVersion: sipVersion}
}
func (rl *RequestLine) Raw() (string, error) {
	result := ""
	if err := rl.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(rl.method)) > 0 {
		result += fmt.Sprintf("%v", strings.ToUpper(rl.method))
	}
	if rl.RequestUri != nil {
		res, err := rl.RequestUri.Raw()
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf(" %s", res)
	}
	if rl.SipVersion != nil {
		res, err := rl.SipVersion.Raw()
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf(" %s", res)
	}
	result += "\r\n"
	return result, nil
}
func (rl *RequestLine) String() string {
	result:=""
	if len(strings.TrimSpace(rl.method))>0{
		result+=fmt.Sprintf("method: %s,", rl.method)
	}
	if rl.RequestUri!=nil{
		result+=fmt.Sprintf("%s,",rl.RequestUri.String())
	}
	if rl.SipVersion!=nil{
		result+=fmt.Sprintf("%s,",rl.SipVersion.String())
	}
	result =strings.TrimSuffix(result,",")
	return result
}
func (rl *RequestLine) Parser(raw string) error {
	if rl == nil {
		return errors.New("requestLine caller is not allowed to be nil")
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
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range sip.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")"
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if methodsRegexp.MatchString(raw) {
		rl.method = methodsRegexp.FindString(raw)
		raw = methodsRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimLeft(raw, " ")
		raw = strings.TrimRight(raw," ")
		raw = strings.TrimPrefix(raw," ")
		raw = strings.TrimSuffix(raw," ")
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
		rl.SipVersion = new(sip.SipVersion)
		if err := rl.SipVersion.Parser(version); err != nil {
			return err
		}
	}
	rl.RequestUri = new(sip.RequestUri)
	return rl.RequestUri.Parser(raw)
}
func (rl *RequestLine) Validator() error {
	if rl == nil {
		return errors.New("requestLine caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(rl.method)) == 0 {
		return errors.New("method is not allowed to be empty")
	}
	if _, ok := sip.Methods[strings.ToUpper(rl.method)]; !ok {
		return errors.New("method is not support")
	}
	if err := rl.SipUri.Validator(); err != nil {
		return err
	}
	if err := rl.SipUri.Validator(); err != nil {
		return err
	}
	return rl.SipVersion.Validator()
}
