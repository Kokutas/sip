package sip

import (
	"fmt"
	"regexp"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.34
//
// 20.34 Route

// The Route header field is used to force routing for a request through
// the listed set of proxies.  Examples of the use of the Route header
// field are in Section 16.12.1.

// Example:

// 	Route: <sip:bigbox3.site3.atlanta.com;lr>,
// 			<sip:server10.biloxi.com;lr>
//
// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.15
//
// Route        =  "Route" HCOLON route-param *(COMMA route-param)
// route-param  =  name-addr *( SEMI rr-param )
// rr-param      =  generic-param
// generic-param  =  token [ EQUAL gen-value ]
// gen-value      =  token / host / quoted-string
// SEMI    =  SWS ";" SWS ; semicolon
// HCOLON  =  *( SP / HTAB ) ":" SWS

type Route struct {
	field     string // "Route"
	nameAddrs []*NameAddr
	source    string // source string
}

func (r *Route) SetField(field string) {
	if regexp.MustCompile(`^(?i)(route)$`).MatchString(field) {
		r.field = strings.Title(field)
	} else {
		r.field = "Route"
	}
}
func (r *Route) GetField() string {
	return r.field
}
func (r *Route) SetNameAddrs(nameAddrs []*NameAddr) {
	r.nameAddrs = nameAddrs
}
func (r *Route) GetNameAddrs() []*NameAddr {
	return r.nameAddrs
}
func (r *Route) GetSource() string {
	return r.source
}
func NewRoute(nameAddrs ...*NameAddr) *Route {
	return &Route{
		field:     "route",
		nameAddrs: nameAddrs,
	}
}
func (r *Route) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(r.field)) == 0 {
		r.field = "route"
	}
	result.WriteString(fmt.Sprintf("%s:", strings.Title(r.field)))
	if r.nameAddrs != nil {
		for _, nameAddr := range r.nameAddrs {
			if nameAddr != nil {
				nameAddrBuilder := nameAddr.Raw()
				result.WriteString(fmt.Sprintf(" <%s>,", nameAddrBuilder.String()))
			}
		}
	}
	temp := result.String()
	result.Reset()
	temp = strings.TrimSuffix(temp, ",")
	result.WriteString(temp)
	result.WriteString("\r\n")
	return
}
func (r *Route) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(route)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	r.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	r.source = raw
	r.nameAddrs = make([]*NameAddr, 0)
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, ",")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// name-addr regexp
	nameAddrRegexp := regexp.MustCompile(`>( )*,`)
	if nameAddrRegexp.MatchString(raw) {
		rawSlice := strings.Split(raw, ",")
		if len(rawSlice) == 1 {
			nameAddrs := regexp.MustCompile(`>`).ReplaceAllString(rawSlice[0], "")
			nameAddrs = regexp.MustCompile(`<`).ReplaceAllString(nameAddrs, "")
			nameAddrs = stringTrimPrefixAndTrimSuffix(nameAddrs, " ")
			nameAddr := new(NameAddr)
			nameAddr.Parse(nameAddrs)
			if len(nameAddr.GetSource()) > 0 {
				r.nameAddrs = append(r.nameAddrs, nameAddr)
			}
		} else {
			for _, raws := range rawSlice {
				nameAddrs := regexp.MustCompile(`>`).ReplaceAllString(raws, "")
				nameAddrs = regexp.MustCompile(`<`).ReplaceAllString(nameAddrs, "")
				nameAddrs = stringTrimPrefixAndTrimSuffix(nameAddrs, " ")
				nameAddr := new(NameAddr)
				nameAddr.Parse(nameAddrs)
				if len(nameAddr.GetSource()) > 0 {
					r.nameAddrs = append(r.nameAddrs, nameAddr)
				}
			}
		}
	} else {
		nameAddrs := regexp.MustCompile(`>`).ReplaceAllString(raw, "")
		nameAddrs = regexp.MustCompile(`<`).ReplaceAllString(nameAddrs, "")
		nameAddrs = stringTrimPrefixAndTrimSuffix(nameAddrs, " ")
		nameAddr := new(NameAddr)
		nameAddr.Parse(nameAddrs)
		if len(nameAddr.GetSource()) > 0 {
			r.nameAddrs = append(r.nameAddrs, nameAddr)
		}
	}
}
