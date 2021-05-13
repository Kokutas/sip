package sip

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// SIP-URI          =  "sip:" [ userinfo ] hostport
//                     uri-parameters [ headers ]
// SIPS-URI         =  "sips:" [ userinfo ] hostport
//                     uri-parameters [ headers ]
// userinfo         =  ( user / telephone-subscriber ) [ ":" password ] "@"
// user             =  1*( unreserved / escaped / user-unreserved )
// user-unreserved  =  "&" / "=" / "+" / "$" / "," / ";" / "?" / "/"
// password         =  *( unreserved / escaped /
//                     "&" / "=" / "+" / "$" / "," )
// hostport         =  host [ ":" port ]
// host             =  hostname / IPv4address / IPv6reference
// hostname         =  *( domainlabel "." ) toplabel [ "." ]
// domainlabel      =  alphanum
//                     / alphanum *( alphanum / "-" ) alphanum
// toplabel         =  ALPHA / ALPHA *( alphanum / "-" ) alphanum
// IPv4address    =  1*3DIGIT "." 1*3DIGIT "." 1*3DIGIT "." 1*3DIGIT
// IPv6reference  =  "[" IPv6address "]"
// IPv6address    =  hexpart [ ":" IPv4address ]
// hexpart        =  hexseq / hexseq "::" [ hexseq ] / "::" [ hexseq ]
// hexseq         =  hex4 *( ":" hex4)
// hex4           =  1*4HEXDIG
// port           =  1*DIGIT

//    The BNF for telephone-subscriber can be found in RFC 2806 [9].  Note,
//    however, that any characters allowed there that are not allowed in
//    the user part of the SIP URI MUST be escaped.

// uri-parameters    =  *( ";" uri-parameter)
// uri-parameter     =  transport-param / user-param / method-param
//                      / ttl-param / maddr-param / lr-param / other-param
// transport-param   =  "transport="
//                      ( "udp" / "tcp" / "sctp" / "tls"
//                      / other-transport)
// other-transport   =  token
// user-param        =  "user=" ( "phone" / "ip" / other-user)
// other-user        =  token
// method-param      =  "method=" Method
// ttl-param         =  "ttl=" ttl
// maddr-param       =  "maddr=" host
// lr-param          =  "lr"
// other-param       =  pname [ "=" pvalue ]
// pname             =  1*paramchar
// pvalue            =  1*paramchar
// paramchar         =  param-unreserved / unreserved / escaped
// param-unreserved  =  "[" / "]" / "/" / ":" / "&" / "+" / "$"

// headers         =  "?" header *( "&" header )
// header          =  hname "=" hvalue
// hname           =  1*( hnv-unreserved / unreserved / escaped )
// hvalue          =  *( hnv-unreserved / unreserved / escaped )
// hnv-unreserved  =  "[" / "]" / "/" / "?" / ":" / "+" / "$"

type SipUri struct {
	schema string
	*UserInfo
	*HostPort
	*Parameters
	headers map[string]interface{}
}

func (su *SipUri) Schema() string {
	return su.schema
}
func (su *SipUri) SetSchema(schema string) {
	su.schema = schema
}

func (su *SipUri) Headers() map[string]interface{} {
	return su.headers
}

func (su *SipUri) SetHeaders(headers map[string]interface{}) {
	su.headers = headers
}
func NewSipUri(schema string, userInfo *UserInfo, hostPort *HostPort, parameters *Parameters, headers map[string]interface{}) *SipUri {
	return &SipUri{schema: schema, UserInfo: userInfo, HostPort: hostPort, Parameters: parameters, headers: headers}
}

func (su *SipUri) Raw() (string, error) {
	result := ""
	if err := su.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(su.schema)) > 0 {
		result += fmt.Sprintf("%s:", su.schema)
	}
	if su.UserInfo != nil {
		res, err := su.UserInfo.Raw()
		if err != nil {
			return "", err
		}
		result += res
	}
	if su.HostPort != nil {
		res, err := su.HostPort.Raw()
		if err != nil {
			return "", err
		}
		result += "@" + res
	}
	if su.Parameters != nil {
		res, err := su.Parameters.Raw()
		if err != nil {
			return "", err
		}
		result += res
	}
	if su.headers != nil {
		headers := ""
		for k, v := range su.headers {
			headers += fmt.Sprintf("%v=%v&", k, v)
		}
		headers = strings.TrimSuffix(headers, "&")
		if len(headers) > 0 {
			result += "?" + headers
		}
	}
	return result, nil
}
func (su *SipUri) String() string {
	result := ""
	if len(strings.TrimSpace(su.schema)) > 0 {
		result += fmt.Sprintf("schema: %s,", su.schema)
	}
	if su.UserInfo!=nil {
		result += fmt.Sprintf("%s,", su.UserInfo.String())
	}
	if su.HostPort!=nil {
		result += fmt.Sprintf("%s,", su.HostPort.String())
	}
	if su.Parameters!=nil {
		result += fmt.Sprintf("%s,", su.Parameters.String())
	}
	if su.headers!=nil {
		result += fmt.Sprintf("headers: %s,", su.headers)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (su *SipUri) Parser(raw string) error {
	if su == nil {
		return errors.New("sipUri caller is not allowed to be nil")
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
	schemaRegexp := regexp.MustCompile(`^\w+:`)
	if schemaRegexp.MatchString(raw) {
		schema := regexp.MustCompile(`\w+`).FindString(schemaRegexp.FindString(raw))
		su.schema = strings.ToLower(schema)
	}
	userinfoRegexp := regexp.MustCompile(`.*@`)
	if userinfoRegexp.MatchString(raw) {
		ui := userinfoRegexp.FindString(raw)
		ui = schemaRegexp.ReplaceAllString(ui, "")
		ui = regexp.MustCompile(`@`).ReplaceAllString(ui, "")
		su.UserInfo = new(UserInfo)
		if err := su.UserInfo.Parser(ui); err != nil {
			return err
		}
	}
	hostportRegexp := regexp.MustCompile(`@.*`)
	if hostportRegexp.MatchString(raw) {
		hp := hostportRegexp.FindString(raw)
		hp = regexp.MustCompile(`@`).ReplaceAllString(hp, "")
		hp = regexp.MustCompile(`;.*`).ReplaceAllString(hp, "")
		su.HostPort = new(HostPort)
		if err := su.HostPort.Parser(hp); err != nil {
			return err
		}
	}
	parametersRegexp := regexp.MustCompile(`;.*`)
	if parametersRegexp.MatchString(raw) {
		ps := parametersRegexp.FindString(raw)
		ps = regexp.MustCompile(`\?.*`).ReplaceAllString(ps, "")
		su.Parameters = new(Parameters)
		if err := su.Parameters.Parser(ps); err != nil {
			return err
		}
	}
	headersRegexp := regexp.MustCompile(`\?.*`)
	if headersRegexp.MatchString(raw) {
		su.headers = make(map[string]interface{})
		h := headersRegexp.FindString(raw)
		h = regexp.MustCompile(`\?`).ReplaceAllString(h, "")
		hs := strings.Split(h, "&")
		for _, vh := range hs {
			vs := strings.Split(vh, "=")
			if len(vs) > 1 {
				su.headers[vs[0]] = vs[1]
			} else {
				su.headers[vs[0]] = ""
			}
		}
	}
	return nil
}
func (su *SipUri) Validator() error {
	if su == nil {
		return errors.New("sipUri caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(su.schema)) == 0 {
		return errors.New("schema is not allowed to be empty")
	}
	if err := su.UserInfo.Validator(); err != nil {
		return err
	}
	if err := su.HostPort.Validator(); err != nil {
		return err
	}
	return nil
}
