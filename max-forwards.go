package sip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-8.1.1.6
//
// 8.1.1.6 Max-Forwards
// The Max-Forwards header field serves to limit the number of hops a
// request can transit on the way to its destination.  It consists of an
// integer that is decremented by one at each hop.  If the Max-Forwards
// value reaches 0 before the request reaches its destination, it will
// be rejected with a 483(Too Many Hops) error response.

// A UAC MUST insert a Max-Forwards header field into each request it
// originates with a value that SHOULD be 70.  This number was chosen to
// be sufficiently large to guarantee that a request would not be
// dropped in any SIP network when there were no loops, but not so large
// as to consume proxy resources when a loop does occur.  Lower values
// should be used with caution and only in networks where topologies are
// known by the UA.

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.22
//
// 20.22 Max-Forwards
// The Max-Forwards header field must be used with any SIP method to
// limit the number of proxies or gateways that can forward the request
// to the next downstream server.  This can also be useful when the
// client is attempting to trace a request chain that appears to be
// failing or looping in mid-chain.

// The Max-Forwards value is an integer in the range 0-255 indicating
// the remaining number of times this request message is allowed to be
// forwarded.  This count is decremented by each server that forwards
// the request.  The recommended initial value is 70.

// This header field should be inserted by elements that can not
// otherwise guarantee loop detection.  For example, a B2BUA should
// insert a Max-Forwards header field.

// Example:

// 	Max-Forwards: 6

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
// Max-Forwards  =  "Max-Forwards" HCOLON 1*DIGIT

type MaxForwards struct {
	field    string // "Max-Forwards"
	forwards uint8  // The Max-Forwards value is an integer in the range 0-255 indicating the remaining number of times this request message is allowed to be forwarded.
	source   string // max-forwards header line source string
}

func (maxForwards *MaxForwards) SetField(field string) {
	if regexp.MustCompile(`^(?i)(max-forwards)$`).MatchString(field) {
		maxForwards.field = strings.Title(field)
	}
}
func (maxForwards *MaxForwards) GetField() string {
	return maxForwards.field
}
func (maxForwards *MaxForwards) SetForwards(forwards uint8) {
	maxForwards.forwards = forwards
}
func (maxForwards *MaxForwards) GetForwards() uint8 {
	return maxForwards.forwards
}
func (maxForwards *MaxForwards) SetSource(source string) {
	maxForwards.source = source
}
func (maxForwards *MaxForwards) GetSource() string {
	return maxForwards.source
}
func NewMaxForwards(forwards uint8) *MaxForwards {
	return &MaxForwards{
		forwards: forwards,
	}
}
func (maxForwards *MaxForwards) Raw() string {
	result := ""
	if len(strings.TrimSpace(maxForwards.field)) > 0 {
		result += fmt.Sprintf("%s:", maxForwards.field)
	} else {
		result += "Max-Forwards:"
	}
	result += fmt.Sprintf(" %d", maxForwards.forwards)
	result += "\r\n"
	return result
}

func (maxForwards *MaxForwards) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// max-forwards field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(max-forwards)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	maxForwards.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	maxForwards.source = raw
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// forwards regexp
	forwardsRegexp := regexp.MustCompile(`\d+`)
	if forwardsRegexp.MatchString(raw) {
		forwards := forwardsRegexp.FindString(raw)
		if len(forwards) > 0 {
			forward, _ := strconv.Atoi(forwards)
			maxForwards.forwards = uint8(forward)
		}
	}
}
