package sip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc3261.html#section-7.2
//
// 7.2 Responses
// SIP responses are distinguished from requests by having a Status-Line
// as their start-line.  A Status-Line consists of the protocol version
// followed by a numeric Status-Code and its associated textual phrase,
// with each element separated by a single SP character.

// No CR or LF is allowed except in the final CRLF sequence.

// 	Status-Line  =  SIP-Version SP Status-Code SP Reason-Phrase CRLF

// The Status-Code is a 3-digit integer result code that indicates the
// outcome of an attempt to understand and satisfy a request.  The
// Reason-Phrase is intended to give a short textual description of the
// Status-Code.  The Status-Code is intended for use by automata,
// whereas the Reason-Phrase is intended for the human user.  A client
// is not required to examine or display the Reason-Phrase.

// While this specification suggests specific wording for the reason
// phrase, implementations MAY choose other text, for example, in the
// language indicated in the Accept-Language header field of the
// request.
// The first digit of the Status-Code defines the class of response.
// The last two digits do not have any categorization role.  For this
// reason, any response with a status code between 100 and 199 is
// referred to as a "1xx response", any response with a status code
// between 200 and 299 as a "2xx response", and so on.  SIP/2.0 allows
// six values for the first digit:

// 	1xx: Provisional -- request received, continuing to process the
// 		request;

// 	2xx: Success -- the action was successfully received, understood,
// 		and accepted;

// 	3xx: Redirection -- further action needs to be taken in order to
// 		complete the request;

// 	4xx: Client Error -- the request contains bad syntax or cannot be
// 		fulfilled at this server;

// 	5xx: Server Error -- the server failed to fulfill an apparently
// 		valid request;

// 	6xx: Global Failure -- the request cannot be fulfilled at any
// 		server.

// Section 21 defines these classes and describes the individual codes.

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Status-Line     =  SIP-Version SP Status-Code SP Reason-Phrase CRLF
// Status-Code     =  Informational
//                /   Redirection
//                /   Success
//                /   Client-Error
//                /   Server-Error
//                /   Global-Failure
//                /   extension-code
// extension-code  =  3DIGIT
// Reason-Phrase   =  *(reserved / unreserved / escaped
//                    / UTF8-NONASCII / UTF8-CONT / SP / HTAB)

type StatusLine struct {
	schema       string
	version      float64
	statusCode   uint
	reasonPhrase string
	isOrder      bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order        chan string // It is convenient to record the order of the original parameter fields when parsing
	source       string      // status-line source string
}

func (statusLine *StatusLine) SetSchema(schema string) {
	if regexp.MustCompile(`(?i)(sip|sips)`).MatchString(schema) {
		statusLine.schema = schema
	} else {
		statusLine.schema = "sip"
	}
}
func (statusLine *StatusLine) GetSchema() string {
	return statusLine.schema
}
func (statusLine *StatusLine) SetVersion(version float64) {
	statusLine.version = version
}
func (statusLine *StatusLine) GetVersion() float64 {
	return statusLine.version
}
func (statusLine *StatusLine) SetStatusCode(statusCode uint) {
	statusLine.statusCode = statusCode
}
func (statusLine *StatusLine) GetStatusCode() uint {
	return statusLine.statusCode
}
func (statusLine *StatusLine) SetReasonPhrase(reasonPhrease string) {
	statusLine.reasonPhrase = reasonPhrease
}
func (statusLine *StatusLine) GetReasonPhrase() string {
	return statusLine.reasonPhrase
}
func (statusLine *StatusLine) GetIsOrder() bool {
	return statusLine.isOrder
}
func (statusLine *StatusLine) GetOrder() []string {
	result := make([]string, 0)
	if statusLine.order == nil {
		return result
	}
	for data := range statusLine.order {
		result = append(result, data)
	}
	return result
}
func (statusLine *StatusLine) SetSource(source string) {
	statusLine.source = source
}
func (statusLine *StatusLine) GetSource() string {
	return statusLine.source
}
func NewStatusLine(schema string, version float64, statusCode uint, reasonPhrase string) *StatusLine {
	return &StatusLine{
		schema:       schema,
		version:      version,
		statusCode:   statusCode,
		reasonPhrase: reasonPhrase,
		isOrder:      false,
		order:        make(chan string, 1024),
	}
}
func (statusLine *StatusLine) Raw() string {
	result := ""
	if statusLine.isOrder {
		for data := range statusLine.order {
			result += data
		}
		statusLine.isOrder = false
		result += "\r\n"
		return result
	}
	if len(strings.TrimSpace(statusLine.schema)) == 0 {
		statusLine.schema = "sip"
	}
	// schema: sip,sips,tel etc.
	if len(strings.TrimSpace(statusLine.schema)) > 0 {
		result += strings.ToUpper(statusLine.schema)
	}
	// version: 2.0
	result += fmt.Sprintf("/%1.1f", statusLine.version)
	result += fmt.Sprintf(" %03d", statusLine.statusCode)
	if len(strings.TrimSpace(statusLine.reasonPhrase)) > 0 {
		result += fmt.Sprintf(" %s", statusLine.reasonPhrase)
	}
	result += "\r\n"
	return result
}
func (statusLine *StatusLine) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// schema regexp string
	schemasRegexpStr := `^(?i)(`
	for _, v := range schemas {
		schemasRegexpStr += v + "|"
	}
	schemasRegexpStr = strings.TrimSuffix(schemasRegexpStr, "|")
	schemasRegexpStr += ")( )?"
	// schema and version regexp
	schemaAndVersionRegexp := regexp.MustCompile(schemasRegexpStr + `/( )?\d\.\d`)

	if !schemaAndVersionRegexp.MatchString(raw) {
		return
	}
	statusLine.source = raw
	// status-line order
	statusLine.statuslineOrder(raw)
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// schema regexp
	schemaRegexp := regexp.MustCompile(schemasRegexpStr)
	// version regexp
	versionRegexp := regexp.MustCompile(`\d\.[0-9]{1}`)
	if schemaAndVersionRegexp.MatchString(raw) {
		schemaAndVersion := schemaAndVersionRegexp.FindString(raw)
		schemaAndVersion = stringTrimPrefixAndTrimSuffix(schemaAndVersion, " ")
		if schemaRegexp.MatchString(schemaAndVersion) {
			schema := schemaRegexp.FindString(schemaAndVersion)
			schema = stringTrimPrefixAndTrimSuffix(schema, " ")
			statusLine.schema = schema
		}
		if versionRegexp.MatchString(schemaAndVersion) {
			versions := versionRegexp.Find([]byte(schemaAndVersion))
			version, _ := strconv.ParseFloat(string(versions), 64)
			statusLine.version = version
		}
		raw = strings.ReplaceAll(raw, schemaAndVersion, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// status-code regexp
	statusCodeRegexp := regexp.MustCompile(`\d{3}`)
	if statusCodeRegexp.MatchString(raw) {
		statusCodes := statusCodeRegexp.FindString(raw)
		statusCode, _ := strconv.Atoi(statusCodes)
		statusLine.statusCode = uint(statusCode)
		raw = strings.ReplaceAll(raw, statusCodes, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		statusLine.reasonPhrase = raw
	}
}
func (statusLine *StatusLine) statuslineOrder(raw string) {
	if statusLine.order == nil {
		statusLine.order = make(chan string, 1024)
	}
	statusLine.isOrder = true
	defer close(statusLine.order)
	statusLine.order <- raw
}
