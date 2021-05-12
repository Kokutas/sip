package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sip"
	"sip/util"
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
	Field    string `json:"field"`
	Forwards uint8  `json:"forwards"`
}

func CreateMaxForwards() sip.Sip {
	return &MaxForwards{}
}
func NewMaxForwards(forwards uint8) sip.Sip {
	return &MaxForwards{
		Field:    "Max-Forwards",
		Forwards: forwards,
	}
}
func (mf *MaxForwards) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, mf) {
		return result
	}
	result += fmt.Sprintf("%v:", strings.Title(mf.Field))
	result += fmt.Sprintf(" %v", mf.Forwards)
	result += "\r\n"
	return result
}
func (mf *MaxForwards) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, mf) {
		return result
	}
	data, err := json.Marshal(mf)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (mf *MaxForwards) Parser(raw string) error {
	if mf == nil {
		return errors.New("max-forwards caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")

	fieldRegexp := regexp.MustCompile(`(?i)(max-forwards)`)
	if fieldRegexp.MatchString(raw) {
		mf.Field = fieldRegexp.FindString(raw)
	}
	forwardsRegexp := regexp.MustCompile(`\d+`)
	if forwardsRegexp.MatchString(raw) {
		forwards, err := strconv.Atoi(forwardsRegexp.FindString(raw))
		if err != nil {
			return err
		}
		mf.Forwards = uint8(forwards)
	} else {
		return errors.New("forwards is not match")
	}
	return nil
}
func (mf *MaxForwards) Validator() error {
	if mf == nil {
		return errors.New("max-forwards caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(mf.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(max-forwards)`).Match([]byte(mf.Field)) {
		return errors.New("field is not match")
	}

	return nil
}
