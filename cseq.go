package sip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-8.1.1.5
//
// 8.1.1.5 CSeq
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

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.16
//
// 20.16 CSeq
// A CSeq header field in a request contains a single decimal sequence
// number and the request method.  The sequence number MUST be
// expressible as a 32-bit unsigned integer.  The method part of CSeq is
// case-sensitive.  The CSeq header field serves to order transactions
// within a dialog, to provide a means to uniquely identify
// transactions, and to differentiate between new requests and request
// retransmissions.  Two CSeq header fields are considered equal if the
// sequence number and the request method are identical.
// Example:

// 	CSeq: 4711 INVITE

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// CSeq  =  "CSeq" HCOLON 1*DIGIT LWS Method

type CSeq struct {
	field   string      // "CSeq"
	number  uint32      // sequence number
	method  string      // method
	isOrder bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order   chan string // It is convenient to record the order of the original parameter fields when parsing
	source  string      // source string
}

func (cSeq *CSeq) SetField(field string) {
	if regexp.MustCompile(`^(?i)(cseq)$`).MatchString(field) {
		cSeq.field = strings.Title(field)
	} else {
		cSeq.field = "CSeq"
	}
}
func (cSeq *CSeq) GetField() string {
	return cSeq.field
}
func (cSeq *CSeq) SetNumber(number uint32) {
	cSeq.number = number
}
func (cSeq *CSeq) GetNumber() uint32 {
	return cSeq.number
}
func (cSeq *CSeq) SetMethod(method string) {
	cSeq.method = method
}
func (cSeq *CSeq) GetMethod() string {
	return cSeq.method
}
func (cSeq *CSeq) GetSource() string {
	return cSeq.source
}
func NewCSeq(number uint32, method string) *CSeq {
	return &CSeq{
		field:   "CSeq",
		number:  number,
		method:  method,
		isOrder: false,
	}
}
func (cSeq *CSeq) Raw() string {
	result := ""
	if cSeq.isOrder {
		for data := range cSeq.order {
			result += data
		}
		cSeq.isOrder = false
		result += "\r\n"
		return result
	}
	if len(strings.TrimSpace(cSeq.field)) == 0 {
		cSeq.field = "CSeq"
	}
	result += fmt.Sprintf("%s:", cSeq.field)
	result += fmt.Sprintf(" %d", cSeq.number)
	if len(strings.TrimSpace(cSeq.method)) > 0 {
		result += fmt.Sprintf(" %s", strings.ToUpper(cSeq.method))
	}
	result += "\r\n"
	return result
}
func (cSeq *CSeq) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// field regexp
	fieldRegexp := regexp.MustCompile(`^(?i)(cseq)( )*:`)
	if !fieldRegexp.MatchString(raw) {
		return
	}
	cSeq.field = regexp.MustCompile(`:`).ReplaceAllString(fieldRegexp.FindString(raw), "")
	cSeq.source = raw
	// cseq order
	cSeq.cseqOrder(raw)
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// sequence number regexp
	numberRegexp := regexp.MustCompile(`\d+`)
	if numberRegexp.MatchString(raw) {
		numbers := numberRegexp.FindString(raw)
		if len(numbers) > 0 {
			number, _ := strconv.Atoi(numbers)
			cSeq.number = uint32(number)
			raw = regexp.MustCompile(`.*`+numbers).ReplaceAllString(raw, "")
			raw = stringTrimPrefixAndTrimSuffix(raw, " ")
		}
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// method regexp string
	methodsRegexpStr := `^(?i)(`
	for _, v := range methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")( )?"
	// method regexp
	methodRegexp := regexp.MustCompile(methodsRegexpStr)
	// method regexp
	if len(raw) > 0 {
		if methodRegexp.MatchString(raw) {
			cSeq.method = methodRegexp.FindString(raw)
		}
	}
}
func (cSeq *CSeq) cseqOrder(raw string) {
	cSeq.order = make(chan string, 1024)
	cSeq.isOrder = true
	defer close(cSeq.order)
	cSeq.order <- raw
}
