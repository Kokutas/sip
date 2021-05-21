package sip

import (
	"fmt"
	"regexp"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-8.1.1.2
//
// 8.1.1.2 To
// The To header field first and foremost specifies the desired
// "logical" recipient of the request, or the address-of-record of the
// user or resource that is the target of this request.  This may or may
// not be the ultimate recipient of the request.  The To header field
// MAY contain a SIP or SIPS URI, but it may also make use of other URI
// schemes (the tel URL (RFC 2806 [9]), for example) when appropriate.
// All SIP implementations MUST support the SIP URI scheme.  Any
// implementation that supports TLS MUST support the SIPS URI scheme.
// The To header field allows for a display name.

// A UAC may learn how to populate the To header field for a particular
// request in a number of ways.  Usually the user will suggest the To
// header field through a human interface, perhaps inputting the URI
// manually or selecting it from some sort of address book.  Frequently,
// the user will not enter a complete URI, but rather a string of digits
// or letters (for example, "bob").  It is at the discretion of the UA
// to choose how to interpret this input.  Using the string to form the
// user part of a SIP URI implies that the UA wishes the name to be
// resolved in the domain to the right-hand side (RHS) of the at-sign in
// the SIP URI (for instance, sip:bob@example.com).  Using the string to
// form the user part of a SIPS URI implies that the UA wishes to
// communicate securely, and that the name is to be resolved in the
// domain to the RHS of the at-sign.  The RHS will frequently be the
// home domain of the requestor, which allows for the home domain to
// process the outgoing request.  This is useful for features like
// "speed dial" that require interpretation of the user part in the home
// domain.  The tel URL may be used when the UA does not wish to specify
// the domain that should interpret a telephone number that has been
// input by the user.  Rather, each domain through which the request
// passes would be given that opportunity.  As an example, a user in an
// airport might log in and send requests through an outbound proxy in
// the airport.  If they enter "411" (this is the phone number for local
// directory assistance in the United States), that needs to be
// interpreted and processed by the outbound proxy in the airport, not
// the user's home domain.  In this case, tel:411 would be the right
// choice.

// A request outside of a dialog MUST NOT contain a To tag; the tag in
// the To field of a request identifies the peer of the dialog.  Since
// no dialog is established, no tag is present.

// For further information on the To header field, see Section 20.39.
// The following is an example of a valid To header field:

//    To: Carol <sip:carol@chicago.com>

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.39
//
// 20.39 To
// The To header field specifies the logical recipient of the request.

// The optional "display-name" is meant to be rendered by a human-user
// interface.  The "tag" parameter serves as a general mechanism for
// dialog identification.

// See Section 19.3 for details of the "tag" parameter.

// Comparison of To header fields for equality is identical to
// comparison of From header fields.  See Section 20.10 for the rules
// for parsing a display name, URI and URI parameters, and header field
// parameters.

// The compact form of the To header field is t.

// The following are examples of valid To header fields:

//    To: The Operator <sip:operator@cs.columbia.edu>;tag=287447
//    t: sip:+12125551212@server.phone2net.com

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// To        =  ( "To" / "t" ) HCOLON ( name-addr
// 				/ addr-spec ) *( SEMI to-param )
// to-param  =  tag-param / generic-param

type To struct {
	field   string // "To" / "t"
	name    string // display-name
	spec    string // named spec of URI,recommend set be uri spec <uri>,example: <sip:xxx>/"sip:xxx"/sip:xxx
	schema  string // sip,sips,tel etc.
	user    string // user part
	host    string // host part
	port    string // port part
	tag     string // tag
	generic string // generic-param
	source  string // to header line source string
}

func (t *To) SetField(field string) {
	if regexp.MustCompile(`^(?i)(to|t)$`).MatchString(field) {
		t.field = strings.Title(field)
	}
}
func (t *To) GetField() string {
	return t.field
}
func (t *To) SetName(name string) {
	t.name = name
}
func (t *To) GetName() string {
	return t.name
}
func (t *To) SetSpec(spec string) {
	t.spec = spec
}
func (t *To) GetSpec() string {
	return t.spec
}
func (t *To) SetSchema(schema string) {
	t.schema = schema
}
func (t *To) GetSchema() string {
	return t.schema
}
func (t *To) SetUser(user string) {
	t.user = user
}
func (t *To) GetUser() string {
	return t.user
}
func (t *To) SetHost(host string) {
	t.host = host
}
func (t *To) GetHost() string {
	return t.host
}
func (t *To) SetPort(port string) {
	t.port = port
}
func (t *To) GetPort() string {
	return t.port
}
func (t *To) SetTag(tag string) {
	t.tag = tag
}
func (t *To) GetTag() string {
	return t.tag
}
func (t *To) SetGeneric(generic string) {
	t.generic = generic
}
func (t *To) GetGeneric() string {
	return t.generic
}
func (t *To) SetSource(source string) {
	t.source = source
}
func (t *To) GetSource() string {
	return t.source
}
func NewTo(name, spec, schema, user, host, port, tag, generic string) *To {
	return &To{
		name:    name,
		spec:    spec,
		schema:  schema,
		user:    user,
		host:    host,
		port:    port,
		tag:     tag,
		generic: generic,
	}
}

func (t *To) Raw() string {
	result := ""
	if len(strings.TrimSpace(t.field)) > 0 {
		result += fmt.Sprintf("%s:", t.field)
	} else {
		result += fmt.Sprintf("%s:", strings.Title("To"))
	}
	if len(strings.TrimSpace(t.name)) > 0 {
		if strings.Contains(t.name, "\"") {
			result += fmt.Sprintf(" %s", t.name)
		} else {
			result += fmt.Sprintf(" \"%s\"", t.name)
		}
	}
	uri := ""
	if len(strings.TrimSpace(t.schema)) > 0 {
		uri += fmt.Sprintf("%s:", strings.ToLower(t.schema))
	}
	if len(strings.TrimSpace(t.user)) > 0 {
		uri += t.user
	}
	if len(strings.TrimSpace(t.host)) > 0 {
		uri += fmt.Sprintf("@%s", t.host)
	}
	if len(strings.TrimSpace(t.port)) > 0 {
		uri += fmt.Sprintf(":%s", t.port)
	}
	if len(uri) > 0 {
		switch strings.TrimSpace(t.spec) {
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
	if len(strings.TrimSpace(t.tag)) > 0 {
		result += fmt.Sprintf(";tag=%s", t.tag)
	}
	if len(strings.TrimSpace(t.generic)) > 0 {
		result += fmt.Sprintf(";%s", t.generic)
	}
	result += "\r\n"
	return result
}

func (t *To) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// to field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(to|t)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	t.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	t.source = raw
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
			t.name = name
			raw = regexp.MustCompile(`.*`+name).ReplaceAllString(raw, "")
			raw = stringTrimPrefixAndTrimSuffix(raw, " ")
		}
	}
	//uri spec  regexp: named spec of URI
	switch {
	case regexp.MustCompile(`'.*?` + schemasRegexpStr).MatchString(raw):
		t.spec = "'"
		raw = regexp.MustCompile(`.*'`).ReplaceAllString(raw, "")
	case regexp.MustCompile(`".*?` + schemasRegexpStr).MatchString(raw):
		t.spec = "\""
		raw = regexp.MustCompile(`.*"`).ReplaceAllString(raw, "")
	case regexp.MustCompile(`<.*?` + schemasRegexpStr).MatchString(raw):
		t.spec = "<"
		raw = regexp.MustCompile(`.*<`).ReplaceAllString(raw, "")
	default:
		t.spec = ""
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// schema regexp
	schemaRegexp := regexp.MustCompile(schemasRegexpStr)
	if schemaRegexp.MatchString(raw) {
		schema := schemaRegexp.FindString(raw)
		schema = regexp.MustCompile(`:`).ReplaceAllString(schema, "")
		schema = stringTrimPrefixAndTrimSuffix(schema, " ")
		t.schema = schema
	}

	// user regexp
	userRegexp := regexp.MustCompile(schemasRegexpStr + `.*@`)
	if userRegexp.MatchString(raw) {
		user := userRegexp.FindString(raw)
		user = regexp.MustCompile(schemasRegexpStr).ReplaceAllString(user, "")
		user = regexp.MustCompile(`@`).ReplaceAllString(user, "")
		user = stringTrimPrefixAndTrimSuffix(user, " ")
		if len(user) > 0 {
			t.user = user
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
			t.host = host
			raw = regexp.MustCompile(`.*`+host).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// port regexp
	portRegexp := regexp.MustCompile(`.*?:\d+`)
	if portRegexp.MatchString(raw) {
		port := portRegexp.FindString(raw)
		port = regexp.MustCompile(`.*:`).ReplaceAllString(port, "")
		port = stringTrimPrefixAndTrimSuffix(port, " ")
		if len(port) > 0 {
			t.port = port
			raw = regexp.MustCompile(`.*`+port).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// tag regexp
	tagRegexp := regexp.MustCompile(`(?i)(tag)( )?=.*`)
	if tagRegexp.MatchString(raw) {
		tag := tagRegexp.FindString(raw)
		tag = regexp.MustCompile(`(?i)tag( )?=`).ReplaceAllString(tag, "")
		tag = regexp.MustCompile(`;.*`).ReplaceAllString(tag, "")
		tag = stringTrimPrefixAndTrimSuffix(tag, " ")
		if len(tag) > 0 {
			t.tag = tag
			raw = regexp.MustCompile(`.*`+tag).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// generic regexp
	if len(raw) > 0 {
		t.generic = raw
	}
}
