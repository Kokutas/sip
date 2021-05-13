package header

import (
	"errors"
	"fmt"
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
	field    string
	sequence uint64
	method   string
}

func (cseq *CSeq) Field() string {
	return cseq.field
}

func (cseq *CSeq) SetField(field string) {
	cseq.field = field
}

func (cseq *CSeq) Sequence() uint64 {
	return cseq.sequence
}

func (cseq *CSeq) SetSequence(sequence uint64) {
	cseq.sequence = sequence
}

func (cseq *CSeq) Method() string {
	return cseq.method
}

func (cseq *CSeq) SetMethod(method string) {
	cseq.method = method
}

func NewCSeq(sequence uint64, method string) *CSeq {
	return &CSeq{
		field:    "CSeq",
		sequence: sequence,
		method:   method,
	}
}
func (cseq *CSeq) Raw() (string, error) {
	result := ""
	if err := cseq.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(cseq.field)) == 0 {
		cseq.field = "CSeq"
	}
	result += fmt.Sprintf("%s:", cseq.field)
	result += fmt.Sprintf(" %d", cseq.sequence)
	if len(strings.TrimSpace(cseq.method)) > 0 {
		result += fmt.Sprintf(" %s", strings.ToUpper(cseq.method))
	}
	result += "\r\n"
	return result, nil
}
func (cseq *CSeq) String() string {
	result := ""
	if len(strings.TrimSpace(cseq.field)) > 0 {
		result += fmt.Sprintf("field: %s,", cseq.field)
	}
	result += fmt.Sprintf("sequence: %d,", cseq.sequence)
	if len(strings.TrimSpace(cseq.method)) > 0 {
		result += fmt.Sprintf("method: %s,", cseq.method)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (cseq *CSeq) Parser(raw string) error {
	if cseq == nil {
		return errors.New("cseq caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}

	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(cseq).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		cseq.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
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
		cseq.sequence = uint64(cq)
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
		cseq.method = methodsRegexp.FindString(raw)
	}

	return nil
}
func (cseq *CSeq) Validator() error {
	if cseq == nil {
		return errors.New("cseq caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(cseq.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(cseq)`).Match([]byte(cseq.field)) {
		return errors.New("field is not match")
	}
	if len(strings.TrimSpace(cseq.method)) == 0 {
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
	if !methodsRegexp.MatchString(cseq.method) {
		return errors.New("method is not match")
	}
	return nil
}
