package sip

import (
	"encoding/json"
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
	Schema      string `json:"schema"`
	*UserInfo   `json:"userinfo"`
	*HostPort   `json:"hostport"`
	*Parameters `json:"uri-parameters,omitempty"`
	Headers     map[string]interface{} `json:"headers,omitempty"`
}

func CreateSipUri() Sip {
	return &SipUri{}
}
func NewSipUri(schema string, userinfo *UserInfo, hostport *HostPort, parameters *Parameters, headers map[string]interface{}) Sip {
	return &SipUri{
		schema,
		userinfo,
		 hostport,
		parameters,
		headers,
	}
}

func (su *SipUri) Raw() string {
	result := ""
	if su == nil {
		return result
	}
	if len(strings.TrimSpace(su.Schema)) > 0 {
		result += fmt.Sprintf("%s:", su.Schema)
	}
	if su.UserInfo != nil {
		result += su.UserInfo.Raw()
	}
	if su.HostPort != nil {
		result += "@" + su.HostPort.Raw()
	}
	if su.Parameters != nil {
		result += su.Parameters.Raw()
	}
	if su.Headers != nil {
		headers := ""
		for k, v := range su.Headers {
			headers += fmt.Sprintf("%v=%v&", k, v)
		}
		headers = strings.TrimSuffix(headers, "&")
		if len(headers) > 0 {
			result += "?" + headers
		}
	}
	return result
}
func (su *SipUri) JsonString() string {
	result := ""
	if su == nil {
		return result
	}
	data, err := json.Marshal(su)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (su *SipUri) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if su == nil {
		return errors.New("sipUri caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	schemaRegexp := regexp.MustCompile(`^\w+\:`)
	if schemaRegexp.MatchString(raw) {
		schema := regexp.MustCompile(`\w+`).FindString(schemaRegexp.FindString(raw))
		su.Schema = strings.ToLower(schema)
	}
	userinfoRegexp := regexp.MustCompile(`.*@`)
	if userinfoRegexp.MatchString(raw) {
		ui := userinfoRegexp.FindString(raw)
		ui = schemaRegexp.ReplaceAllString(ui, "")
		ui = regexp.MustCompile(`@`).ReplaceAllString(ui, "")
		su.UserInfo = CreateUserInfo()
		if err := su.UserInfo.Parser(ui); err != nil {
			return err
		}
	}
	hostportRegexp := regexp.MustCompile(`@.*`)
	if hostportRegexp.MatchString(raw) {
		hp := hostportRegexp.FindString(raw)
		hp = regexp.MustCompile(`@`).ReplaceAllString(hp, "")
		hp = regexp.MustCompile(`;.*`).ReplaceAllString(hp, "")
		su.HostPort = CreateHostPort().(*HostPort)
		if err := su.HostPort.Parser(hp); err != nil {
			return err
		}
	}
	parametersRegexp := regexp.MustCompile(`;.*`)
	if parametersRegexp.MatchString(raw) {
		ps := parametersRegexp.FindString(raw)
		ps = regexp.MustCompile(`\?.*`).ReplaceAllString(ps, "")
		su.Parameters = CreateParameters().(*Parameters)
		if err := su.Parameters.Parser(ps); err != nil {
			return err
		}
	}
	headersRegexp := regexp.MustCompile(`\?.*`)
	if headersRegexp.MatchString(raw) {
		su.Headers = make(map[string]interface{})
		h := headersRegexp.FindString(raw)
		h = regexp.MustCompile(`\?`).ReplaceAllString(h, "")
		hs := strings.Split(h, "&")
		for _, vh := range hs {
			vs := strings.Split(vh, "=")
			if len(vs) > 1 {
				su.Headers[vs[0]] = vs[1]
			} else {
				su.Headers[vs[0]] = ""
			}
		}
	}

	return nil
}
func (su *SipUri) Validator() error {
	if su == nil {
		return errors.New("sipUri caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(su.Schema)) == 0 {
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
