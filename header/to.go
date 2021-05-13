package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"github.com/kokutas/sip"
	"github.com/kokutas/sip/util"
	"strings"
)

//The To header field specifies the logical recipient of the request.
//
//The optional "display-name" is meant to be rendered by a human-user
//interface.  The "tag" parameter serves as a general mechanism for
//dialog identification.
//
//See Section 19.3 for details of the "tag" parameter.
//Comparison of To header fields for equality is identical to
//comparison of To header fields.  See Section 20.10 for the rules
//for parsing a display name, URI and URI parameters, and header field
//parameters.
//
//The compact form of the To header field is t.
//
//The following are examples of valid To header fields:
//
//To: The Operator <sip:operator@cs.columbia.edu>;tag=287447
//t: sip:+12125551212@server.phone2net.com
//To        =  ( "To" / "t" ) HCOLON ( name-addr
//			/ addr-spec ) *( SEMI to-param )
//to-param  =  tag-param / generic-param
type To struct {
	field       string
	displayName string
	addr        *sip.SipUri
	tag         string
}

func (to *To) Field() string {
	return to.field
}

func (to *To) SetField(field string) {
	to.field = field
}

func (to *To) DisplayName() string {
	return to.displayName
}

func (to *To) SetDisplayName(displayName string) {
	to.displayName = displayName
}

func (to *To) Addr() *sip.SipUri {
	return to.addr
}

func (to *To) SetAddr(addr *sip.SipUri) {
	to.addr = addr
}

func (to *To) Tag() string {
	return to.tag
}

func (to *To) SetTag(tag string) {
	to.tag = tag
}

func NewTo(displayName string, addr *sip.SipUri, tag string) *To {
	return &To{
		field:       "To",
		displayName: displayName,
		addr:        addr,
		tag:         tag,
	}
}

func (to *To) Raw() (string, error) {
	result := ""
	if err := to.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(to.field)) == 0 {
		to.field = "To"
	}
	//The optional "display-name" is meant to be rendered by a human-user
	//interface.  The "tag" parameter serves as a general mechanism for
	//dialog identification.
	result += fmt.Sprintf("%s:", strings.Title(to.field))
	if len(strings.TrimSpace(to.displayName)) > 0 {
		if strings.Contains(to.displayName, "\"") {
			result += fmt.Sprintf(" %s", to.displayName)
		} else {
			result += fmt.Sprintf(" \"%v\"", to.displayName)
		}
		//	if !reflect.DeepEqual(nil, to.Addr) {
		//		res,err:=to.addr.Raw()
		//		if err!=nil{
		//			return "",err
		//		}
		//		result += fmt.Sprintf(" %v", res)
		//	}
		//} else {
		//	if !reflect.DeepEqual(nil, to.addr) {
		//		res,err:=to.addr.Raw()
		//		if err!=nil{
		//			return "",err
		//		}
		//		result += fmt.Sprintf(" <%v>", res)
		//	}
	}
	// If the name-addr / addr-spec need to be commented as follows, release the comment in the display name
	if !reflect.DeepEqual(nil, to.addr) {
		res, err := to.addr.Raw()
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf(" <%v>", res)
	}
	if len(strings.TrimSpace(to.tag)) > 0 {
		result += fmt.Sprintf(";tag=%v", to.tag)
	}
	result += "\r\n"
	return result, nil
}
func (to *To) String() string {
	result := ""
	if len(strings.TrimSpace(to.field)) > 0 {
		result += fmt.Sprintf("field: %s,", to.field)
	}
	if len(strings.TrimSpace(to.displayName)) > 0 {
		result += fmt.Sprintf("display-name: %s,", to.displayName)
	}
	if to.addr != nil {
		result += fmt.Sprintf("%s,", to.addr.String())
	}
	if len(strings.TrimSpace(to.tag)) > 0 {
		result += fmt.Sprintf("tag: %s,", to.tag)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (to *To) Parser(raw string) error {
	if to == nil {
		return errors.New("to caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(to).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		to.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
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
		to.displayName = util.TrimPrefixAndSuffix(displayNames, " ")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// tag
	tagRegexp := regexp.MustCompile(`;(?i)tag=.*`)
	if tagRegexp.MatchString(raw) {
		tag := tagRegexp.FindString(raw)
		raw = tagRegexp.ReplaceAllString(raw, "")
		to.tag = regexp.MustCompile(`;(?i)tag=`).ReplaceAllString(tag, "")
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
		to.addr = new(sip.SipUri)
		if err := to.addr.Parser(addr); err != nil {
			return err
		}
	}

	return nil
}
func (to *To) Validator() error {
	if to == nil {
		return errors.New("to caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(to.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(to)`).Match([]byte(to.field)) {
		return errors.New("field is not match")
	}
	return to.addr.Validator()
}
