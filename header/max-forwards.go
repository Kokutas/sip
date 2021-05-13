package header

import (
	"errors"
	"fmt"
	"regexp"
	"github.com/kokutas/sip/util"
	"strconv"
	"strings"
)

// The Max-Forwards header field must be used with any SIP method to
//  limit the number of proxies or gateways that can forward the request
//  to the next downstream server. This can also be useful when the
//  client is attempting to trace a request chain that appears to be
//  failing or looping in mid-chain.
//  The Max-Forwards value is an integer in the range 0-255 indicating
//  the remaining number of times this request message is allowed to be
//  forwarded. This count is decremented by each server that forwards
//  the request. The recommended initial value is 70.
//  This header field should be inserted by elements that can not
//  otherwise guarantee loop detection. For example, a B2BUA should
//  insert a Max-Forwards header field.
//  Example:
//  Max-Forwards: 6

//  Max-Forwards = "Max-Forwards" HCOLON 1*DIGIT

type MaxForwards struct {
	field    string
	forwards uint8
}

func (mf *MaxForwards) Field() string {
	return mf.field
}

func (mf *MaxForwards) SetField(field string) {
	mf.field = field
}

func (mf *MaxForwards) Forwards() uint8 {
	return mf.forwards
}

func (mf *MaxForwards) SetForwards(forwards uint8) {
	mf.forwards = forwards
}

func NewMaxForwards(forwards uint8) *MaxForwards {
	return &MaxForwards{
		field:    "Max-Forwards",
		forwards: forwards,
	}
}
func (mf *MaxForwards) Raw() (string, error) {
	result := ""
	if err := mf.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(mf.field)) == 0 {
		mf.field = "Max-Forwards"
	}
	result += fmt.Sprintf("%s:", strings.Title(mf.field))
	result += fmt.Sprintf(" %d", mf.forwards)
	result += "\r\n"
	return result, nil
}
func (mf *MaxForwards) String() string {
	result := ""
	if len(strings.TrimSpace(mf.field)) > 0 {
		result += fmt.Sprintf("field: %s,", mf.field)
	}
	result += fmt.Sprintf("forwards: %d,", mf.forwards)
	result = strings.TrimSuffix(result, ",")
	return result
}
func (mf *MaxForwards) Parser(raw string) error {
	if mf == nil {
		return errors.New("max-forwards caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	fieldRegexp := regexp.MustCompile(`(?i)(max-forwards)`)
	if fieldRegexp.MatchString(raw) {
		mf.field = fieldRegexp.FindString(raw)
	}
	forwardsRegexp := regexp.MustCompile(`\d+`)
	if forwardsRegexp.MatchString(raw) {
		forwards, err := strconv.Atoi(forwardsRegexp.FindString(raw))
		if err != nil {
			return err
		}
		mf.forwards = uint8(forwards)
	} else {
		return errors.New("forwards is not match")
	}
	return nil
}
func (mf *MaxForwards) Validator() error {
	if mf == nil {
		return errors.New("max-forwards caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(mf.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(max-forwards)`).Match([]byte(mf.field)) {
		return errors.New("field is not match")
	}

	return nil
}
