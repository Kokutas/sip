package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sip"
	"sip/util"
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
	Field       string      `json:"field"`
	DisplayName string      `json:"display-name"`
	Addr        *sip.SipUri `json:"name-addr/addr-spec"`
	Tag         string      `json:"tag"`
}

func CreateTo() sip.Sip {
	return &To{}
}
func NewTo(displayName string, addr *sip.SipUri, tag string) sip.Sip {
	return &To{
		Field:       "To",
		DisplayName: displayName,
		Addr:        addr,
		Tag:         tag,
	}
}

func (to *To) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, to) {
		return result
	}
	//The optional "display-name" is meant to be rendered by a human-user
	//interface.  The "tag" parameter serves as a general mechanism for
	//dialog identification.
	result += fmt.Sprintf("%v:", strings.Title(to.Field))
	if len(strings.TrimSpace(to.DisplayName)) > 0 {
		if strings.Contains(to.DisplayName, "\"") {
			result += fmt.Sprintf(" %v", to.DisplayName)
		} else {
			result += fmt.Sprintf(" \"%v\"", to.DisplayName)
		}
		//	if !reflect.DeepEqual(nil, to.Addr) {
		//		result += fmt.Sprintf(" %v", to.Addr.Raw())
		//	}
		//} else {
		//	if !reflect.DeepEqual(nil, to.Addr) {
		//		result += fmt.Sprintf(" <%v>", to.Addr.Raw())
		//	}
	}
	// If the name-addr / addr-spec need to be commented as follows, release the comment in the display name
	if !reflect.DeepEqual(nil, to.Addr) {
		result += fmt.Sprintf(" <%v>", to.Addr.Raw())
	}
	if len(strings.TrimSpace(to.Tag)) > 0 {
		result += fmt.Sprintf(";tag=%v", to.Tag)
	}
	result += "\r\n"
	return result
}
func (to *To) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, to) {
		return result
	}
	if data, err := json.Marshal(to); err == nil {
		result = fmt.Sprintf("%s", data)
	}
	return result
}
func (to *To) Parser(raw string) error {
	if to == nil {
		return errors.New("to caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")

	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(to).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		to.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
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
		to.DisplayName = util.TrimPrefixAndSuffix(displayNames, " ")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// tag
	tagRegexp := regexp.MustCompile(`;(?i)tag=.*`)
	if tagRegexp.MatchString(raw) {
		tag := tagRegexp.FindString(raw)
		raw = tagRegexp.ReplaceAllString(raw, "")
		to.Tag = regexp.MustCompile(`;(?i)tag=`).ReplaceAllString(tag, "")
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
		to.Addr = sip.CreateSipUri().(*sip.SipUri)
		if err := to.Addr.Parser(addr); err != nil {
			return err
		}
	}

	return nil
}
func (to *To) Validator() error {
	if to == nil {
		return errors.New("to caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(to.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(to)`).Match([]byte(to.Field)) {
		return errors.New("field is not match")
	}
	return to.Addr.Validator()
}

func (to *To) String() string {
	result := ""
	if reflect.DeepEqual(nil, to) {
		return result
	}
	//The optional "display-name" is meant to be rendered by a human-user
	//interface.  The "tag" parameter serves as a general mechanism for
	//dialog identification.
	if len(strings.TrimSpace(to.DisplayName)) > 0 {
		if strings.Contains(to.DisplayName, "\"") {
			result += fmt.Sprintf("%v", to.DisplayName)
		} else {
			result += fmt.Sprintf("\"%v\"", to.DisplayName)
		}
		//	if !reflect.DeepEqual(nil, to.Addr) {
		//		result += fmt.Sprintf(" %v", to.Addr.Raw())
		//	}
		//} else {
		//	if !reflect.DeepEqual(nil, to.Addr) {
		//		result += fmt.Sprintf(" <%v>", to.Addr.Raw())
		//	}
		// If the name-addr / addr-spec need to be commented as follows, release the comment in the display name
		if !reflect.DeepEqual(nil, to.Addr) {
			result += fmt.Sprintf(" <%v>", to.Addr.Raw())
		}
	} else {
		// If the name-addr / addr-spec need to be commented as follows, release the comment in the display name
		if !reflect.DeepEqual(nil, to.Addr) {
			result += fmt.Sprintf("<%v>", to.Addr.Raw())
		}
	}

	if len(strings.TrimSpace(to.Tag)) > 0 {
		result += fmt.Sprintf(";tag=%v", to.Tag)
	}
	return result
}
