package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"github.com/kokutas/sip"
	"github.com/kokutas/sip/util"
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
	field       string
	displayName string
	addr        *sip.SipUri
	cpq         float32
	cpExpires   int
	extension   map[string]interface{}
}

func (c *Contact) Field() string {
	return c.field
}

func (c *Contact) SetField(field string) {
	c.field = field
}

func (c *Contact) DisplayName() string {
	return c.displayName
}

func (c *Contact) SetDisplayName(displayName string) {
	c.displayName = displayName
}

func (c *Contact) Addr() *sip.SipUri {
	return c.addr
}

func (c *Contact) SetAddr(addr *sip.SipUri) {
	c.addr = addr
}

func (c *Contact) Cpq() float32 {
	return c.cpq
}

func (c *Contact) SetCpq(cpq float32) {
	c.cpq = cpq
}

func (c *Contact) CpExpires() int {
	return c.cpExpires
}

func (c *Contact) SetCpExpires(cpExpires int) {
	c.cpExpires = cpExpires
}

func (c *Contact) Extension() map[string]interface{} {
	return c.extension
}

func (c *Contact) SetExtension(extension map[string]interface{}) {
	c.extension = extension
}
func NewContact(displayName string, addr *sip.SipUri, cpq float32, cpExpires int, extension map[string]interface{}) *Contact {
	return &Contact{field: "Contact", displayName: displayName, addr: addr, cpq: cpq, cpExpires: cpExpires, extension: extension}
}

func (c *Contact) Raw() (string, error) {
	result := ""
	if err := c.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(c.field)) == 0 {
		c.field = "Contact"
	}
	result += fmt.Sprintf("%v:", strings.Title(c.field))
	if len(strings.TrimSpace(c.displayName)) > 0 {
		if strings.Contains(c.displayName, "\"") {
			result += fmt.Sprintf(" %v", c.displayName)
		} else {
			result += fmt.Sprintf(" \"%v\"", c.displayName)
		}
		//	if !reflect.DeepEqual(nil, c.addr) {
		//		res,err:=c.addr.Raw()
		//		if err!=nil{
		//			return "",err
		//		}
		//		result += fmt.Sprintf(" %v", res)
		//	}
		//} else {
		//	if !reflect.DeepEqual(nil, c.addr) {
		//		res,err:=c.addr.Raw()
		//		if err!=nil{
		//			return "",err
		//		}
		//		result += fmt.Sprintf(" <%v>", res)
		//	}
	}
	// If the name-addr / addr-spec need to be commented as follows, release the comment in the display name
	if !reflect.DeepEqual(nil, c.addr) {
		res, err := c.addr.Raw()
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf(" <%v>", res)
	}
	if c.cpq >= 0 {
		result += fmt.Sprintf(";q=%1.1f", c.cpq)
	}
	if c.cpExpires >= 0 {
		result += fmt.Sprintf(";expires=%v", c.cpExpires)
	}
	if c.extension != nil {
		extensions := ""
		for k, v := range c.extension {
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
	return result, nil
}
func (c *Contact) String() string {
	result := ""
	if len(strings.TrimSpace(c.field)) > 0 {
		result += fmt.Sprintf("field: %s,", c.field)
	}
	if len(strings.TrimSpace(c.displayName)) > 0 {
		result += fmt.Sprintf("display-name: %s,", c.displayName)
	}
	if c.addr != nil {
		result += fmt.Sprintf("%s,", c.addr.String())
	}
	if c.cpq > 0 {
		result += fmt.Sprintf("q: %1.1f,", c.cpq)
	}
	if c.cpExpires >= 0 {
		result += fmt.Sprintf("expires: %d,", c.cpExpires)
	}
	if c.extension != nil {
		result += fmt.Sprintf("c-extension: %v,", c.extension)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (c *Contact) Parser(raw string) error {
	if c == nil {
		return errors.New("c caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(c).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		c.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
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
		c.displayName = util.TrimPrefixAndSuffix(displayNames, " ")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// c-p-q regexp
	cpqRegexp := regexp.MustCompile(`;(?i)(q).*?=\d*\.\d*`)
	if cpqRegexp.MatchString(raw) {
		cpq, err := strconv.ParseFloat(regexp.MustCompile(`;(?i)(q).*?=`).ReplaceAllString(cpqRegexp.FindString(raw), ""), 10)
		if err != nil {
			return nil
		}
		c.cpq = float32(cpq)
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
		c.cpExpires = cpExpires
		raw = cpExpiresRegexp.ReplaceAllString(raw, "")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// c-extension regexp
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
			c.extension = m
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
		c.addr = new(sip.SipUri)
		if err := c.addr.Parser(addr); err != nil {
			return err
		}
	}

	return nil
}
func (c *Contact) Validator() error {

	if c == nil {
		return errors.New("contact caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(c.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(contact)`).Match([]byte(c.field)) {
		return errors.New("field is not match")
	}
	return c.addr.Validator()
}
