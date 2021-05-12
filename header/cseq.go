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

// The CSeq header field serves as a way to identify and order
// transactions.  It consists of a sequence number and a method.  The
// method MUST match that of the request.  For non-REGISTER requests
// outside of a dialog, the sequence number value is arbitrary.  The
// sequence number value MUST be expressible as a 32-bit unsigned
// integer and MUST be less than 2**31.  As long as it follows the above
// guidelines, a client may use any mechanism it would like to select
// CSeq header field values.

// Section 12.2.1.1 discusses construction of the CSeq for requests
// within a dialog.

// Example:

// 	CSeq: 4711 INVITE
type CSeq struct {
	Field    string `json:"field"`
	Sequence uint64 `json:"sequence"`
	Method   string `json:"method"`
}

func CreateCSeq() sip.Sip {
	return &CSeq{}
}

func NewCSeq(sequence uint64, method string) sip.Sip {
	return &CSeq{
		Field:    "CSeq",
		Sequence: sequence,
		Method:   method,
	}
}
func (cseq *CSeq) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, cseq) {
		return result
	}
	result += fmt.Sprintf("%v:", cseq.Field)
	result += fmt.Sprintf(" %v", cseq.Sequence)
	if len(strings.TrimSpace(cseq.Method)) > 0 {
		result += fmt.Sprintf(" %v", strings.ToUpper(cseq.Method))
	}
	result += "\r\n"
	return result
}
func (cseq *CSeq) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, cseq) {
		return result
	}
	data, err := json.Marshal(cseq)
	if err != nil {
		return result
	}
	result += fmt.Sprintf("%s", data)
	return result
}
func (cseq *CSeq) Parser(raw string) error {
	if cseq == nil {
		return errors.New("cseq caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")

	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(cseq).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		cseq.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// seq regexp
	sequenceRegexp := regexp.MustCompile(`\d+`)
	if sequenceRegexp.MatchString(raw) {
		cq, err := strconv.Atoi(sequenceRegexp.FindString(raw))
		if err != nil {
			return err
		}
		cseq.Sequence = uint64(cq)
	}
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range sip.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")"
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if methodsRegexp.MatchString(raw) {
		cseq.Method = methodsRegexp.FindString(raw)
	}

	return nil
}
func (cseq *CSeq) Validator() error {
	if cseq == nil {
		return errors.New("cseq caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(cseq.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(cseq)`).Match([]byte(cseq.Field)) {
		return errors.New("field is not match")
	}
	if len(strings.TrimSpace(cseq.Method)) == 0 {
		return errors.New("method is not allowed to be empty")
	}
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range sip.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")"
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if !methodsRegexp.MatchString(cseq.Method) {
		return errors.New("method is not match")
	}
	return nil
}
