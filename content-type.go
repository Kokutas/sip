package sip

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

//https://www.rfc-editor.org/rfc/rfc3261.html#section-20.15
//
// 20.15 Content-Type
//
// The Content-Type header field indicates the media type of the
// message-body sent to the recipient.  The "media-type" element is
// defined in [H3.7].  The Content-Type header field MUST be present if
// the body is not empty.  If the body is empty, and a Content-Type
// header field is present, it indicates that the body of the specific
// type has zero length (for example, an empty audio file).

// The compact form of the header field is c.

// Examples:

// 	Content-Type: application/sdp
//  c: text/html; charset=ISO-8859-4
//  Content-Type: application/pkcs7-mime; smime-type=enveloped-data;
//				 name=smime.p7m
//  Content-Type: multipart/signed;
//         		 protocol="application/pkcs7-signature";
//         		 micalg=sha1; boundary=boundary42
//
//https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Content-Type     =  ( "Content-Type" / "c" ) HCOLON media-type
// media-type       =  m-type SLASH m-subtype *(SEMI m-parameter)
// m-type           =  discrete-type / composite-type
// discrete-type    =  "text" / "image" / "audio" / "video"
//                     / "application" / extension-token
// composite-type   =  "message" / "multipart" / extension-token
// extension-token  =  ietf-token / x-token
// ietf-token       =  token
// x-token          =  "x-" token
// m-subtype        =  extension-token / iana-token
// iana-token       =  token
// m-parameter      =  m-attribute EQUAL m-value
// m-attribute      =  token
// m-value          =  token / quoted-string
// SLASH   =  SWS "/" SWS ; slash
//
type ContentType struct {
	field     string      //"Content-Type" / "c"
	mType     string      // media-type =  m-type SLASH m-subtype *(SEMI m-parameter),m-type = discrete-type / composite-type,discrete-type =  "text" / "image" / "audio" / "video"/ "application" / extension-token,composite-type =  "message" / "multipart" / extension-token
	mSubType  string      // m-subtype =  extension-token / iana-token
	parameter sync.Map    //m-parameter =  m-attribute EQUAL m-value, m-attribute =  token,m-value =  token / quoted-string
	isOrder   bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order     chan string // It is convenient to record the order of the original parameter fields when parsing
	source    string      // source string
}

func (c *ContentType) SetField(field string) {
	if regexp.MustCompile(`(?i)(?:^content-type|c)$`).MatchString(field) {
		c.field = field
	} else {
		c.field = "Content-Type"
	}
}
func (c *ContentType) GetField() string {
	return c.field
}
func (c *ContentType) SetMType(mType string) {
	c.mType = mType
}
func (c *ContentType) GetMType() string {
	return c.mType
}
func (c *ContentType) SetMSubType(mSubType string) {
	c.mSubType = mSubType
}
func (c *ContentType) GetMSubType() string {
	return c.mSubType
}
func (c *ContentType) SetParameter(parameter sync.Map) {
	c.parameter = parameter
}
func (c *ContentType) GetParameter() sync.Map {
	return c.parameter
}
func (c *ContentType) GetSource() string {
	return c.source
}
func NewContentType(mType string, mSubType string, parameter sync.Map) *ContentType {
	return &ContentType{
		field:     "Content-Type",
		mType:     mType,
		mSubType:  mSubType,
		parameter: parameter,
		isOrder:   false,
	}
}
func (c *ContentType) Raw() (result strings.Builder) {

	if len(strings.TrimSpace(c.field)) == 0 {
		c.field = "Content-Type"
	}
	result.WriteString(fmt.Sprintf("%s:", c.field))

	if len(strings.TrimSpace(c.mType)) > 0 {
		result.WriteString(fmt.Sprintf(" %s", c.mType))
	}
	if len(strings.TrimSpace(c.mSubType)) > 0 {
		if len(result.String()) > 0 {
			result.WriteString(fmt.Sprintf("/%s", c.mSubType))
		} else {
			result.WriteString(fmt.Sprintf(" %s", c.mSubType))
		}
	}
	if c.isOrder {
		c.isOrder = false
		for orders := range c.order {
			ordersSlice := strings.Split(orders, "=")
			if len(ordersSlice) == 1 {
				if val, ok := c.parameter.LoadAndDelete(ordersSlice[0]); ok {
					if strings.Contains(fmt.Sprintf("%v", val), "/") {
						result.WriteString(fmt.Sprintf(";%v=\"%v\"", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], val))
					}
				} else {
					result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
				}
			} else {
				if val, ok := c.parameter.LoadAndDelete(ordersSlice[0]); ok {
					if strings.Contains(fmt.Sprintf("%v", val), "/") {
						result.WriteString(fmt.Sprintf(";%v=\"%v\"", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], val))
					}
				} else {
					result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], ordersSlice[1]))
				}
			}
		}
	}
	c.parameter.Range(func(key, value interface{}) bool {
		if reflect.ValueOf(value).IsValid() {
			if reflect.ValueOf(value).IsZero() {
				result.WriteString(fmt.Sprintf(";%v", key))
				return true
			}
			if strings.Contains(fmt.Sprintf("%v", value), "/") {
				result.WriteString(fmt.Sprintf(";%v=\"%v\"", key, value))
			} else {
				result.WriteString(fmt.Sprintf(";%v=%v", key, value))
			}
			return true
		}
		result.WriteString(fmt.Sprintf(";%v", key))
		return true
	})
	result.WriteString("\r\n")
	return
}
func (c *ContentType) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(content-type|c)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	c.source = raw
	c.parameter = sync.Map{}

	field := fieldRegexp.FindString(raw)
	field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
	field = stringTrimPrefixAndTrimSuffix(field, " ")
	c.field = field
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// parameter regexp
	parameterRegexp := regexp.MustCompile(`;.*`)
	if parameterRegexp.MatchString(raw) {
		// content-type order
		c.contenttypeOrder(parameterRegexp.FindString(raw))
		rawSlice := strings.Split(parameterRegexp.FindString(raw), ";")
		for _, raws := range rawSlice {
			raws = stringTrimPrefixAndTrimSuffix(raws, " ")
			if len(strings.TrimSpace(raws)) == 0 {
				continue
			}
			kvs := strings.Split(raws, "=")
			if len(kvs) == 1 {
				c.parameter.Store(kvs[0], "")
			} else {
				c.parameter.Store(kvs[0], kvs[1])
			}
		}
		raw = parameterRegexp.ReplaceAllString(raw, "")
		raw = stringTrimPrefixAndTrimSuffix(raw, " ")
		mAndSubTypes := strings.Split(raw, "/")
		if len(mAndSubTypes) == 1 {
			c.mType = raw
		} else {
			c.mType = mAndSubTypes[0]
			c.mSubType = mAndSubTypes[1]
		}

	} else {
		mAndSubTypes := strings.Split(raw, "/")
		if len(mAndSubTypes) == 1 {
			c.mType = raw
		} else {
			c.mType = mAndSubTypes[0]
			c.mSubType = mAndSubTypes[1]
		}
	}
}

func (c *ContentType) contenttypeOrder(raw string) {
	c.isOrder = true
	c.order = make(chan string, 1024)
	defer close(c.order)
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	rawSlice := strings.Split(raw, ";")
	for _, raws := range rawSlice {
		c.order <- raws
	}
}
