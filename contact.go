package sip

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-8.1.1.8
//
// 8.1.1.8 Contact
// The Contact header field provides a SIP or SIPS URI that can be used
// to contact that specific instance of the UA for subsequent requests.
// The Contact header field MUST be present and contain exactly one SIP
// or SIPS URI in any request that can result in the establishment of a
// dialog.  For the methods defined in this specification, that includes
// only the INVITE request.  For these requests, the scope of the
// Contact is global.  That is, the Contact header field value contains
// the URI at which the UA would like to receive requests, and this URI
// MUST be valid even if used in subsequent requests outside of any
// dialogs.

// If the Request-URI or top Route header field value contains a SIPS
// URI, the Contact header field MUST contain a SIPS URI as well.

// For further information on the Contact header field, see Section
// 20.10.
//
// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.10
//
// 20.10 Contact
// A Contact header field value provides a URI whose meaning depends on
// the type of request or response it is in.

// A Contact header field value can contain a display name, a URI with
// URI parameters, and header parameters.

// This document defines the Contact parameters "q" and "expires".
// These parameters are only used when the Contact is present in a
// REGISTER request or response, or in a 3xx response.  Additional
// parameters may be defined in other specifications.

// When the header field value contains a display name, the URI
// including all URI parameters is enclosed in "<" and ">".  If no "<"
// and ">" are present, all parameters after the URI are header
// parameters, not URI parameters.  The display name can be tokens, or a
// quoted string, if a larger character set is desired.

// Even if the "display-name" is empty, the "name-addr" form MUST be
// used if the "addr-spec" contains a comma, semicolon, or question
// mark.  There may or may not be LWS between the display-name and the
// "<".

// These rules for parsing a display name, URI and URI parameters, and
// header parameters also apply for the header fields To and From.

// 	The Contact header field has a role similar to the Location header
// 	field in HTTP.  However, the HTTP header field only allows one
// 	address, unquoted.  Since URIs can contain commas and semicolons
// 	as reserved characters, they can be mistaken for header or
// 	parameter delimiters, respectively.

// The compact form of the Contact header field is m (for "moved").

// Examples:

// 	Contact: "Mr. Watson" <sip:watson@worcester.bell-telephone.com>
// 		;q=0.7; expires=3600,
// 		"Mr. Watson" <mailto:watson@bell-telephone.com> ;q=0.1
// 	m: <sips:bob@192.0.2.4>;expires=60

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Contact        =  ("Contact" / "m" ) HCOLON
//                   ( STAR / (contact-param *(COMMA contact-param)))
// contact-param  =  (name-addr / addr-spec) *(SEMI contact-params)
// name-addr      =  [ display-name ] LAQUOT addr-spec RAQUOT
// addr-spec      =  SIP-URI / SIPS-URI / absoluteURI
// display-name   =  *(token LWS)/ quoted-string

// contact-params     =  c-p-q / c-p-expires
//                       / contact-extension
// c-p-q              =  "q" EQUAL qvalue
// c-p-expires        =  "expires" EQUAL delta-seconds
// contact-extension  =  generic-param
// delta-seconds      =  1*DIGIT
// generic-param  =  token [ EQUAL gen-value ]
// qvalue         =  ( "0" [ "." 0*3DIGIT ] )/ ( "1" [ "." 0*3("0") ] )

type Contact struct {
	field   string      // "Contact" / "m"
	name    string      // display-name
	spec    string      // named spec of URI,recommend set be uri spec <uri>,example: <sip:xxx>/"sip:xxx"/sip:xxx
	schema  string      // sip,sips,tel etc.
	user    string      // user part
	host    string      // host part
	port    uint16      // port part
	q       string      // c-p-q  =  "q" EQUAL qvalue,qvalue = ( "0" [ "." 0*3DIGIT ] )/ ( "1" [ "." 0*3("0") ] )
	expires int         // c-p-expires =  "expires" EQUAL delta-seconds,delta-seconds = 1*DIGIT
	generic sync.Map    // generic-param,contact-extension = generic-param,generic-param =  token [ EQUAL gen-value ]
	isOrder bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order   chan string // It is convenient to record the order of the original parameter fields when parsing
	source  string      // source string
}

func (m *Contact) SetField(field string) {
	if regexp.MustCompile(`^(?i)(contact|m)$`).MatchString(field) {
		m.field = strings.Title(field)
	} else {
		m.field = "Contact"
	}
}
func (m *Contact) GetField() string {
	return m.field
}
func (m *Contact) SetName(name string) {
	m.name = name
}
func (m *Contact) GetName() string {
	return m.name
}
func (m *Contact) SetSpec(spec string) {
	m.spec = spec
}
func (m *Contact) GetSpec() string {
	return m.spec
}
func (m *Contact) SetSchema(schema string) {
	m.schema = schema
}
func (m *Contact) GetSchema() string {
	return m.schema
}
func (m *Contact) SetUser(user string) {
	m.user = user
}
func (m *Contact) GetUser() string {
	return m.user
}
func (m *Contact) SetHost(host string) {
	m.host = host
}
func (m *Contact) GetHost() string {
	return m.host
}
func (m *Contact) SetPort(port uint16) {
	m.port = port
}
func (m *Contact) GetPort() uint16 {
	return m.port
}
func (m *Contact) SetQ(qValue string) {
	m.q = qValue
}
func (m *Contact) GetQ() string {
	return m.q
}
func (m *Contact) SetExpires(expires int) {
	m.expires = expires
}
func (m *Contact) GetExpires() int {
	return m.expires
}
func (m *Contact) SetGeneric(generic sync.Map) {
	m.generic = generic
}
func (m *Contact) GetGeneric() sync.Map {
	return m.generic
}
func (m *Contact) GetSource() string {
	return m.source
}

func NewContact(name, spec, schema, user, host string, port uint16, q string, expires int, generic sync.Map) *Contact {
	return &Contact{
		field:   "Contact",
		name:    name,
		spec:    spec,
		schema:  schema,
		user:    user,
		host:    host,
		port:    port,
		q:       q,
		expires: expires,
		generic: generic,
		isOrder: false,
	}
}
func (m *Contact) Raw() string {
	result := ""
	if m.isOrder {
		for data := range m.order {
			result += data
		}
		m.isOrder = false
		result += "\r\n"
		return result
	}
	if len(strings.TrimSpace(m.field)) == 0 {
		m.field = "Contact"
	}
	result += fmt.Sprintf("%s:", strings.Title(m.field))
	if len(strings.TrimSpace(m.name)) > 0 {
		if strings.Contains(m.name, "\"") {
			result += fmt.Sprintf(" %s", m.name)
		} else {
			result += fmt.Sprintf(" \"%s\"", m.name)
		}
	}
	uri := ""
	if len(strings.TrimSpace(m.schema)) > 0 {
		uri += fmt.Sprintf("%s:", strings.ToLower(m.schema))
	}
	if len(strings.TrimSpace(m.user)) > 0 {
		uri += m.user
	}
	if len(strings.TrimSpace(m.host)) > 0 {
		uri += fmt.Sprintf("@%s", m.host)
	}
	if m.port > 0 {
		uri += fmt.Sprintf(":%v", m.port)
	}
	if len(uri) > 0 {
		switch strings.TrimSpace(m.spec) {
		case "\"":
			result += fmt.Sprintf(" \"%s\"", uri)
		case "'":
			result += fmt.Sprintf(" '%s'", uri)
		case "<":
			result += fmt.Sprintf(" <%s>", uri)
		default:
			result += fmt.Sprintf(" %s", uri)
		}

	}
	if len(strings.TrimSpace(m.q)) > 0 {
		result += fmt.Sprintf(";q=%s", m.q)
	}
	if m.expires >= 0 {
		result += fmt.Sprintf(";expires=%v", m.expires)
	}
	m.generic.Range(func(key, value interface{}) bool {
		if reflect.ValueOf(value).IsValid() {
			if reflect.ValueOf(value).IsZero() {
				result += fmt.Sprintf(";%v", key)
				return true
			}
			result += fmt.Sprintf(";%v=\"%v\"", key, value)
			return true
		}
		result += fmt.Sprintf(";%v", key)
		return true
	})
	result = strings.TrimSuffix(result, ";")
	result += "\r\n"
	return result
}

func (m *Contact) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(contact|m)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	m.source = raw
	m.generic = sync.Map{}
	// contact order
	m.contactOrder(raw)
	m.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")

	// schema regexp
	schemasRegexpStr := `(?i)(`
	for _, v := range schemas {
		schemasRegexpStr += v + "|"
	}
	schemasRegexpStr = strings.TrimSuffix(schemasRegexpStr, "|")
	schemasRegexpStr += ")( )?:"

	// display-name regexp
	nameRegexp := regexp.MustCompile(`.*` + schemasRegexpStr)
	if nameRegexp.MatchString(raw) {
		name := nameRegexp.FindString(raw)
		name = regexp.MustCompile(schemasRegexpStr+`$`).ReplaceAllString(name, "")
		name = regexp.MustCompile(`<$`).ReplaceAllString(name, "")
		name = stringTrimPrefixAndTrimSuffix(name, " ")
		if len(name) > 0 {
			m.name = name
			raw = regexp.MustCompile(`.*`+name).ReplaceAllString(raw, "")
			raw = stringTrimPrefixAndTrimSuffix(raw, " ")
		}
	}
	//uri spec  regexp: named spec of URI
	switch {
	case regexp.MustCompile(`'.*?` + schemasRegexpStr).MatchString(raw):
		m.spec = "'"
		raw = regexp.MustCompile(`.*'`).ReplaceAllString(raw, "")
	case regexp.MustCompile(`".*?` + schemasRegexpStr).MatchString(raw):
		m.spec = "\""
		raw = regexp.MustCompile(`.*"`).ReplaceAllString(raw, "")
	case regexp.MustCompile(`<.*?` + schemasRegexpStr).MatchString(raw):
		m.spec = "<"
		raw = regexp.MustCompile(`.*<`).ReplaceAllString(raw, "")
	default:
		m.spec = ""
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// schema regexp
	schemaRegexp := regexp.MustCompile(schemasRegexpStr)
	if schemaRegexp.MatchString(raw) {
		schema := schemaRegexp.FindString(raw)
		schema = regexp.MustCompile(`:`).ReplaceAllString(schema, "")
		schema = stringTrimPrefixAndTrimSuffix(schema, " ")
		m.schema = schema
	}
	// user regexp
	userRegexp := regexp.MustCompile(schemasRegexpStr + `.*@`)
	if userRegexp.MatchString(raw) {
		user := userRegexp.FindString(raw)
		user = regexp.MustCompile(schemasRegexpStr).ReplaceAllString(user, "")
		user = regexp.MustCompile(`@`).ReplaceAllString(user, "")
		user = stringTrimPrefixAndTrimSuffix(user, " ")
		if len(user) > 0 {
			m.user = user
			raw = regexp.MustCompile(`.*`+user).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// host regexp
	hostRegexp := regexp.MustCompile(`@.*`)
	if hostRegexp.MatchString(raw) {
		host := hostRegexp.FindString(raw)
		host = regexp.MustCompile(`;.*`).ReplaceAllString(host, "")
		host = regexp.MustCompile(`:.*`).ReplaceAllString(host, "")
		host = regexp.MustCompile(`@`).ReplaceAllString(host, "")
		host = stringTrimPrefixAndTrimSuffix(host, " ")
		if len(host) > 0 {
			m.host = host
			raw = regexp.MustCompile(`.*`+host).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// port regexp
	portRegexp := regexp.MustCompile(`.*?:\d+`)
	if portRegexp.MatchString(raw) {
		ports := portRegexp.FindString(raw)
		ports = regexp.MustCompile(`.*:`).ReplaceAllString(ports, "")
		ports = stringTrimPrefixAndTrimSuffix(ports, " ")
		if len(ports) > 0 {
			port, _ := strconv.Atoi(ports)
			m.port = uint16(port)
			raw = regexp.MustCompile(`.*`+ports).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// q regexp
	qRegexp := regexp.MustCompile(`((?i)(?:^q))( )*=`)
	// expires regexp
	expiresRegexp := regexp.MustCompile(`((?i)(?:^expires))( )*=`)
	rawSlice := strings.Split(raw, ";")
	for _, raws := range rawSlice {
		raws = stringTrimPrefixAndTrimSuffix(raws, " ")
		switch {
		case qRegexp.MatchString(raws):
			q := qRegexp.ReplaceAllString(raws, "")
			q = regexp.MustCompile(`"`).ReplaceAllString(q, "")
			if len(q) > 0 {
				m.q = q
			}
		case expiresRegexp.MatchString(raws):
			expires := expiresRegexp.ReplaceAllString(raws, "")
			expires = regexp.MustCompile(`"`).ReplaceAllString(expires, "")
			if len(expires) > 0 {
				expire, _ := strconv.Atoi(expires)
				m.expires = expire
			}
		default:
			// generic regexp
			kvs := strings.Split(raws, "=")
			if len(kvs) == 1 {
				m.generic.Store(kvs[0], "")
			} else {
				m.generic.Store(kvs[0], kvs[1])
			}
		}
	}

}
func (m *Contact) contactOrder(raw string) {
	m.order = make(chan string, 1024)
	m.isOrder = true
	defer close(m.order)
	m.order <- raw
}
