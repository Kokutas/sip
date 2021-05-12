package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sip"
	"sip/util"
	"strconv"
	"strings"
)

// A Contact header field value provides a URI whose meaning depends on
//  the type of request or response it is in.
//  A Contact header field value can contain a display name, a URI with
//  URI parameters, and header parameters.
//  This document defines the Contact parameters "q" and "expires".
//  These parameters are only used when the Contact is present in a
//  REGISTER request or response, or in a 3xx response. Additional
//  parameters may be defined in other specifications.
//  When the header field value contains a display name, the URI
//  including all URI parameters is enclosed in "<" and ">". If no "<"
//  and ">" are present, all parameters after the URI are header
//  parameters, not URI parameters. The display name can be tokens, or a
//  quoted string, if a larger character set is desired.
//  Even if the "display-name" is empty, the "name-addr" form MUST be
//  used if the "addr-spec" contains a comma, semicolon, or question
//  mark. There may or may not be LWS between the display-name and the
//  "<".
//  These rules for parsing a display name, URI and URI parameters, and
//  header parameters also apply for the header fields To and From.
//  The Contact header field has a role similar to the Location header
//  field in HTTP. However, the HTTP header field only allows one
//  address, unquoted. Since URIs can contain commas and semicolons
//  as reserved characters, they can be mistaken for header or
//  parameter delimiters, respectively.
//  The compact form of the Contact header field is m (for "moved").
//  Examples:
//  Contact: "Mr. Watson" <sip:watson@worcester.bell-telephone.com>
//  	;q=0.7; expires=3600,
//  	"Mr. Watson" <mailto:watson@bell-telephone.com> ;q=0.1
//  m: <sips:bob@192.0.2.4>;expires=60
//
// Contact = ("Contact" / "m" ) HCOLON
//  		  ( STAR / (contact-param *(COMMA contact-param)))
// contact-param = (name-addr / addr-spec) *(SEMI contact-params)
// name-addr = [ display-name ] LAQUOT addr-spec RAQUOT
// addr-spec = SIP-URI / SIPS-URI / absoluteURI
// display-name = *(token LWS)/ quoted-string
// contact-params = c-p-q / c-p-expires
//  				/ contact-extension
// c-p-q = "q" EQUAL qvalue
// c-p-expires = "expires" EQUAL delta-seconds
// contact-extension = generic-param
// delta-seconds = 1*DIGIT
type Contact struct {
	Field       string                 `json:"field"`
	DisplayName string                 `json:"display-name"`
	Addr        *sip.SipUri            `json:"name-addr/addr-spec"`
	CPQ         float32                `json:"q"`
	CPExpires   int                    `json:"expires"`
	Extension   map[string]interface{} `json:"contact-extension"`
}

func CreateContact() sip.Sip {
	return &Contact{}
}
func NewContact(displayName string, addr *sip.SipUri, cpq float32, cpExpires int, extension map[string]interface{}) sip.Sip {
	return &Contact{
		Field:       "Contact",
		DisplayName: displayName,
		Addr:        addr,
		CPQ:         cpq,
		CPExpires:   cpExpires,
		Extension:   extension,
	}
}

func (contact *Contact) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, contact) {
		return result
	}
	result += fmt.Sprintf("%v:", strings.Title(contact.Field))
	if len(strings.TrimSpace(contact.DisplayName)) > 0 {
		if strings.Contains(contact.DisplayName, "\"") {
			result += fmt.Sprintf(" %v", contact.DisplayName)
		} else {
			result += fmt.Sprintf(" \"%v\"", contact.DisplayName)
		}
		//	if !reflect.DeepEqual(nil, contact.Addr) {
		//		result += fmt.Sprintf(" %v", contact.Addr.Raw())
		//	}
		//} else {
		//	if !reflect.DeepEqual(nil, contact.Addr) {
		//		result += fmt.Sprintf(" <%v>", contact.Addr.Raw())
		//	}
	}
	// If the name-addr / addr-spec need to be commented as follows, release the comment in the display name
	if !reflect.DeepEqual(nil, contact.Addr) {
		result += fmt.Sprintf(" <%v>", contact.Addr.Raw())
	}
	if contact.CPQ >= 0 {
		result += fmt.Sprintf(";q=%1.1f", contact.CPQ)
	}
	if contact.CPExpires >= 0 {
		result += fmt.Sprintf(";expires=%v", contact.CPExpires)
	}
	if contact.Extension != nil {
		extensions := ""
		for k, v := range contact.Extension {
			if v == nil {
				extensions += fmt.Sprintf(";%v", k)
			} else {
				extensions += fmt.Sprintf(";%v=%v", k, v)
			}
		}
		if len(strings.TrimSpace(extensions)) > 0 {
			result += extensions
		}
	}
	result += "\r\n"
	return result
}
func (contact *Contact) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, contact) {
		return result
	}
	data, err := json.Marshal(contact)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)

	return result
}
func (contact *Contact) Parser(raw string) error {
	if contact == nil {
		return errors.New("contact caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(contact).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		contact.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}

	// schemas regexp
	schemasRegexpStr := `(?i)(`
	for _, v := range sip.Schemas {
		schemasRegexpStr += v + "|"
	}
	schemasRegexpStr = strings.TrimSuffix(schemasRegexpStr, "|")
	schemasRegexpStr += ")"
	// display name regexp
	displayNameRegexp := regexp.MustCompile(`.*?` + schemasRegexpStr)
	if displayNameRegexp.MatchString(raw) {
		displayNames := displayNameRegexp.FindString(raw)
		displayNames = regexp.MustCompile(schemasRegexpStr).ReplaceAllString(displayNames, "")
		displayNames = regexp.MustCompile(`<`).ReplaceAllString(displayNames, "")
		displayNames = regexp.MustCompile(`:`).ReplaceAllString(displayNames, "")
		raw = regexp.MustCompile(`.*`+displayNames).ReplaceAllString(raw, "")
		contact.DisplayName = util.TrimPrefixAndSuffix(displayNames, " ")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// c-p-q regexp
	cpqRegexp := regexp.MustCompile(`;(?i)(q).*?=\d*\.\d*`)
	if cpqRegexp.MatchString(raw) {
		cpq, err := strconv.ParseFloat(regexp.MustCompile(`;(?i)(q).*?=`).ReplaceAllString(cpqRegexp.FindString(raw), ""), 10)
		if err != nil {
			return nil
		}
		contact.CPQ = float32(cpq)
		raw = cpqRegexp.ReplaceAllString(raw, "")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}

	// c-p-expires regexp
	cpExpiresRegexp := regexp.MustCompile(`;(?i)(expires).*?=\d*`)
	if cpExpiresRegexp.MatchString(raw) {
		cpExpires, err := strconv.Atoi(regexp.MustCompile(`;(?i)(expires).*?=`).ReplaceAllString(cpExpiresRegexp.FindString(raw), ""))
		if err != nil {
			return err
		}
		contact.CPExpires = cpExpires
		raw = cpExpiresRegexp.ReplaceAllString(raw, "")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// contact-extension regexp
	extensionRegexp := regexp.MustCompile(`;.*`)
	if extensionRegexp.MatchString(raw) {
		raw = extensionRegexp.ReplaceAllString(raw, "")
		raw = util.TrimPrefixAndSuffix(raw, " ")
		m := make(map[string]interface{})
		extension := extensionRegexp.FindString(raw)
		extensions := strings.Split(extension, ";")
		for _, v := range extensions {
			if strings.Contains(v, "=") {
				vs := strings.Split(v, "=")
				if len(vs) > 1 {
					m[vs[0]] = vs[1]
				} else {
					m[vs[0]] = ""
				}
			} else {
				m[v] = ""
			}
		}
		if len(m) > 0 {
			contact.Extension = m
		}
	}

	// addr regexp
	addrRegexp := regexp.MustCompile(schemasRegexpStr + `.*`)
	if addrRegexp.MatchString(raw) {
		addr := addrRegexp.FindString(raw)
		addr = util.TrimPrefixAndSuffix(addr, ";")
		addr = regexp.MustCompile(`<`).ReplaceAllString(addr, "")
		addr = regexp.MustCompile(`>.*`).ReplaceAllString(addr, "")
		addr = util.TrimPrefixAndSuffix(addr, ";")
		addr = util.TrimPrefixAndSuffix(addr, " ")
		contact.Addr = sip.CreateSipUri().(*sip.SipUri)
		if err := contact.Addr.Parser(addr); err != nil {
			return err
		}
	}

	return nil
}
func (contact *Contact) Validator() error {

	if contact == nil {
		return errors.New("contact caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(contact.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(contact)`).Match([]byte(contact.Field)) {
		return errors.New("field is not match")
	}
	return contact.Addr.Validator()
}
