package sip

import (
	"fmt"
	"regexp"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-8.1.1.3
//
// 8.1.1.3 From
// The From header field indicates the logical identity of the initiator
// of the request, possibly the user's address-of-record.  Like the To
// header field, it contains a URI and optionally a display name.  It is
// used by SIP elements to determine which processing rules to apply to
// a request (for example, automatic call rejection).  As such, it is
// very important that the From URI not contain IP addresses or the FQDN
// of the host on which the UA is running, since these are not logical
// names.

// The From header field allows for a display name.  A UAC SHOULD use
// the display name "Anonymous", along with a syntactically correct, but
// otherwise meaningless URI (like sip:thisis@anonymous.invalid), if the
// identity of the client is to remain hidden.

// Usually, the value that populates the From header field in requests
// generated by a particular UA is pre-provisioned by the user or by the
// administrators of the user's local domain.  If a particular UA is
// used by multiple users, it might have switchable profiles that
// include a URI corresponding to the identity of the profiled user.
// Recipients of requests can authenticate the originator of a request
// in order to ascertain that they are who their From header field
// claims they are (see Section 22 for more on authentication).

// The From field MUST contain a new "tag" parameter, chosen by the UAC.
// See Section 19.3 for details on choosing a tag.

// For further information on the From header field, see Section 20.20.
// Examples:

//    From: "Bob" <sips:bob@biloxi.com> ;tag=a48s
//    From: sip:+12125551212@phone2net.com;tag=887s
//    From: Anonymous <sip:c8oqz84zk7z@privacy.org>;tag=hyh8

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.20
//
// 20.20 From
// The From header field indicates the initiator of the request.  This
// may be different from the initiator of the dialog.  Requests sent by
// the callee to the caller use the callee's address in the From header
// field.

// The optional "display-name" is meant to be rendered by a human user
// interface.  A system SHOULD use the display name "Anonymous" if the
// identity of the client is to remain hidden.  Even if the "display-
// name" is empty, the "name-addr" form MUST be used if the "addr-spec"
// contains a comma, question mark, or semicolon.  Syntax issues are
// discussed in Section 7.3.1.

// Two From header fields are equivalent if their URIs match, and their
// parameters match. Extension parameters in one header field, not
// present in the other are ignored for the purposes of comparison. This
// means that the display name and presence or absence of angle brackets
// do not affect matching.

// See Section 20.10 for the rules for parsing a display name, URI and
// URI parameters, and header field parameters.

// The compact form of the From header field is f.

// Examples:

//    From: "A. G. Bell" <sip:agb@bell-telephone.com> ;tag=a48s
//    From: sip:+12125551212@server.phone2net.com;tag=887s
//    f: Anonymous <sip:c8oqz84zk7z@privacy.org>;tag=hyh8

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// From        =  ( "From" / "f" ) HCOLON from-spec
// from-spec   =  ( name-addr / addr-spec )
// 				  *( SEMI from-param )
// from-param  =  tag-param / generic-param
// tag-param   =  "tag" EQUAL token

type From struct {
	field   string //"From" / "f"
	name    string // display-name
	spec    string // named spec of URI,recommend set be uri spec <uri>,example: <sip:xxx>/"sip:xxx"/sip:xxx
	schema  string // sip,sips,tel etc.
	user    string // user part
	host    string // host part
	port    string // port part
	tag     string // tag
	generic string // generic-param
	source  string // from header line source string
}

func (f *From) SetField(field string) {
	if regexp.MustCompile(`^(?i)(from|f)$`).MatchString(field) {
		f.field = strings.Title(field)
	}
}
func (f *From) GetField() string {
	return f.field
}
func (f *From) SetName(name string) {
	f.name = name
}
func (f *From) GetName() string {
	return f.name
}
func (f *From) SetSpec(spec string) {
	f.spec = spec
}
func (f *From) GetSpec() string {
	return f.spec
}
func (f *From) SetSchema(schema string) {
	f.schema = schema
}
func (f *From) GetSchema() string {
	return f.schema
}
func (f *From) SetUser(user string) {
	f.user = user
}
func (f *From) GetUser() string {
	return f.user
}
func (f *From) SetHost(host string) {
	f.host = host
}
func (f *From) GetHost() string {
	return f.host
}
func (f *From) SetPort(port string) {
	f.port = port
}
func (f *From) GetPort() string {
	return f.port
}
func (f *From) SetTag(tag string) {
	f.tag = tag
}
func (f *From) GetTag() string {
	return f.tag
}
func (f *From) SetGeneric(generic string) {
	f.generic = generic
}
func (f *From) GetGeneric() string {
	return f.generic
}
func (f *From) SetSource(source string) {
	f.source = source
}
func (f *From) GetSource() string {
	return f.source
}

func NewFrom(name, spec, schema, user, host, port, tag, generic string) *From {
	return &From{
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

func (f *From) Raw() string {
	result := ""
	if len(strings.TrimSpace(f.field)) > 0 {
		result += fmt.Sprintf("%s:", f.field)
	} else {
		result += fmt.Sprintf("%s:", strings.Title("From"))
	}
	if len(strings.TrimSpace(f.name)) > 0 {
		if strings.Contains(f.name, "\"") {
			result += fmt.Sprintf(" %s", f.name)
		} else {
			result += fmt.Sprintf(" \"%s\"", f.name)
		}
	}
	uri := ""
	if len(strings.TrimSpace(f.schema)) > 0 {
		uri += fmt.Sprintf("%s:", strings.ToLower(f.schema))
	}
	if len(strings.TrimSpace(f.user)) > 0 {
		uri += f.user
	}
	if len(strings.TrimSpace(f.host)) > 0 {
		uri += fmt.Sprintf("@%s", f.host)
	}
	if len(strings.TrimSpace(f.port)) > 0 {
		uri += fmt.Sprintf(":%s", f.port)
	}
	if len(uri) > 0 {
		switch strings.TrimSpace(f.spec) {
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
	if len(strings.TrimSpace(f.tag)) > 0 {
		result += fmt.Sprintf(";tag=%s", f.tag)
	}
	if len(strings.TrimSpace(f.generic)) > 0 {
		result += fmt.Sprintf(";%s", f.generic)
	}
	result += "\r\n"
	return result
}

func (f *From) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// from field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(from|f)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	f.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	f.source = raw
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
			f.name = name
			raw = regexp.MustCompile(`.*`+name).ReplaceAllString(raw, "")
			raw = stringTrimPrefixAndTrimSuffix(raw, " ")
		}
	}
	//uri spec  regexp: named spec of URI
	switch {
	case regexp.MustCompile(`'.*?` + schemasRegexpStr).MatchString(raw):
		f.spec = "'"
		raw = regexp.MustCompile(`.*'`).ReplaceAllString(raw, "")
	case regexp.MustCompile(`".*?` + schemasRegexpStr).MatchString(raw):
		f.spec = "\""
		raw = regexp.MustCompile(`.*"`).ReplaceAllString(raw, "")
	case regexp.MustCompile(`<.*?` + schemasRegexpStr).MatchString(raw):
		f.spec = "<"
		raw = regexp.MustCompile(`.*<`).ReplaceAllString(raw, "")
	default:
		f.spec = ""
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// schema regexp
	schemaRegexp := regexp.MustCompile(schemasRegexpStr)
	if schemaRegexp.MatchString(raw) {
		schema := schemaRegexp.FindString(raw)
		schema = regexp.MustCompile(`:`).ReplaceAllString(schema, "")
		schema = stringTrimPrefixAndTrimSuffix(schema, " ")
		f.schema = schema
	}

	// user regexp
	userRegexp := regexp.MustCompile(schemasRegexpStr + `.*@`)
	if userRegexp.MatchString(raw) {
		user := userRegexp.FindString(raw)
		user = regexp.MustCompile(schemasRegexpStr).ReplaceAllString(user, "")
		user = regexp.MustCompile(`@`).ReplaceAllString(user, "")
		user = stringTrimPrefixAndTrimSuffix(user, " ")
		if len(user) > 0 {
			f.user = user
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
			f.host = host
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
			f.port = port
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
			f.tag = tag
			raw = regexp.MustCompile(`.*`+tag).ReplaceAllString(raw, "")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// generic regexp
	if len(raw) > 0 {
		f.generic = raw
	}
}