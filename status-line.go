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
	source       string // source string
}

func (sl *StatusLine) SetSchema(schema string) {
	if regexp.MustCompile(`(?i)(sip|sips)`).MatchString(schema) {
		sl.schema = schema
	} else {
		sl.schema = "sip"
	}
}
func (sl *StatusLine) GetSchema() string {
	return sl.schema
}
func (sl *StatusLine) SetVersion(version float64) {
	sl.version = version
}
func (sl *StatusLine) GetVersion() float64 {
	return sl.version
}
func (sl *StatusLine) SetStatusCode(statusCode uint) {
	sl.statusCode = statusCode
}
func (sl *StatusLine) GetStatusCode() uint {
	return sl.statusCode
}
func (sl *StatusLine) SetReasonPhrase(reasonPhrease string) {
	sl.reasonPhrase = reasonPhrease
}
func (sl *StatusLine) GetReasonPhrase() string {
	return sl.reasonPhrase
}
func (sl *StatusLine) GetSource() string {
	return sl.source
}
func NewStatusLine(schema string, version float64, statusCode uint, reasonPhrase string) *StatusLine {
	return &StatusLine{
		schema:       schema,
		version:      version,
		statusCode:   statusCode,
		reasonPhrase: reasonPhrase,
	}
}
func (sl *StatusLine) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(sl.schema)) == 0 {
		sl.schema = "sip"
	}
	// schema: sip,sips,tel etc.
	if len(strings.TrimSpace(sl.schema)) > 0 {
		result.WriteString(strings.ToUpper(sl.schema))
	}
	// version: 2.0
	result.WriteString(fmt.Sprintf("/%1.1f", sl.version))
	result.WriteString(fmt.Sprintf(" %03d", sl.statusCode))
	if len(strings.TrimSpace(sl.reasonPhrase)) > 0 {
		result.WriteString(fmt.Sprintf(" %s", sl.reasonPhrase))
	}
	result.WriteString("\r\n")
	return result
}
func (sl *StatusLine) Parse(raw string) {
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
	sl.source = raw
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
			sl.schema = schema
		}
		if versionRegexp.MatchString(schemaAndVersion) {
			versions := versionRegexp.Find([]byte(schemaAndVersion))
			version, _ := strconv.ParseFloat(string(versions), 64)
			sl.version = version
		}
		raw = strings.ReplaceAll(raw, schemaAndVersion, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// status-code regexp
	statusCodeRegexp := regexp.MustCompile(`\d+`)
	if statusCodeRegexp.MatchString(raw) {
		statusCodes := statusCodeRegexp.FindString(raw)
		statusCode, _ := strconv.Atoi(statusCodes)
		sl.statusCode = uint(statusCode)
		raw = strings.ReplaceAll(raw, statusCodes, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		sl.reasonPhrase = raw
	}
}
