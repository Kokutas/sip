package sip

// https://www.rfc-editor.org/rfc/rfc3261.html#section-20.43
//
// 20.43 Warning
// The Warning header field is used to carry additional information
// about the status of a response.  Warning header field values are sent
// with responses and contain a three-digit warning code, host name, and
// warning text.

// The "warn-text" should be in a natural language that is most likely
// to be intelligible to the human user receiving the response.  This
// decision can be based on any available knowledge, such as the
// location of the user, the Accept-Language field in a request, or the
// Content-Language field in a response.  The default language is i-
// default [21].

// The currently-defined "warn-code"s are listed below, with a
// recommended warn-text in English and a description of their meaning.
// These warnings describe failures induced by the session description.
// The first digit of warning codes beginning with "3" indicates
// warnings specific to SIP.  Warnings 300 through 329 are reserved for
// indicating problems with keywords in the session description, 330
// through 339 are warnings related to basic network services requested
// in the session description, 370 through 379 are warnings related to
// quantitative QoS parameters requested in the session description, and
// 390 through 399 are miscellaneous warnings that do not fall into one
// of the above categories.
// 300 Incompatible network protocol: One or more network protocols
// contained in the session description are not available.

// 301 Incompatible network address formats: One or more network
// address formats contained in the session description are not
// available.

// 302 Incompatible transport protocol: One or more transport
// protocols described in the session description are not
// available.

// 303 Incompatible bandwidth units: One or more bandwidth
// measurement units contained in the session description were
// not understood.

// 304 Media type not available: One or more media types contained in
// the session description are not available.

// 305 Incompatible media format: One or more media formats contained
// in the session description are not available.

// 306 Attribute not understood: One or more of the media attributes
// in the session description are not supported.

// 307 Session description parameter not understood: A parameter
// other than those listed above was not understood.

// 330 Multicast not available: The site where the user is located
// does not support multicast.

// 331 Unicast not available: The site where the user is located does
// not support unicast communication (usually due to the presence
// of a firewall).
// 370 Insufficient bandwidth: The bandwidth specified in the session
// description or defined by the media exceeds that known to be
// available.

// 399 Miscellaneous warning: The warning text can include arbitrary
// information to be presented to a human user or logged.  A
// system receiving this warning MUST NOT take any automated
// action.

// 	1xx and 2xx have been taken by HTTP/1.1.

// Additional "warn-code"s can be defined through IANA, as defined in
// Section 27.2.

// Examples:

// 	Warning: 307 isi.edu "Session parameter 'foo' not understood"
// 	Warning: 301 isi.edu "Incompatible network address type 'E.164'"

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Warning        =  "Warning" HCOLON warning-value *(COMMA warning-value)
// warning-value  =  warn-code SP warn-agent SP warn-text
// warn-code      =  3DIGIT
// warn-agent     =  hostport / pseudonym
//                   ;  the name or pseudonym of the server adding
//                   ;  the Warning header, for use in debugging
// warn-text      =  quoted-string
// pseudonym      =  token

// Warning : for use in debugging
type Warning struct {
	field     string      // "Warning"
	warnCode  uint        //  warn-code = 3DIGIT
	warnAgent string      // warn-agent =  hostport / pseudonym;the name or pseudonym of the server adding;the Warning header, for use in debugging;pseudonym = token
	warnText  string      // warn-text  =  quoted-string
	isOrder   bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order     chan string // It is convenient to record the order of the original parameter fields when parsing
	source    string      // waring source string
}

func (warning *Warning) SetField(field string) {
	warning.field = field
}
func (warning *Warning) GetField() string {
	return warning.field
}
func (warning *Warning) SetWarnCode(warnCode uint) {
	warning.warnCode = warnCode
}
func (warning *Warning) GetWarnCode() uint {
	return warning.warnCode
}
func (warning *Warning) SetWarnAgent(warnAgent string) {
	warning.warnAgent = warnAgent
}
func (warning *Warning) GetWarnAgent() string {
	return warning.warnAgent
}
func (warning *Warning) SetWarnText(warnText string) {
	warning.warnText = warnText
}
func (warning *Warning) GetWarnText() string {
	return warning.warnText
}
func (warning *Warning) GetIsOrder() bool {
	return warning.isOrder
}
func (warning *Warning) GetOrder() []string {
	result := make([]string, 0)
	if warning.order == nil {
		return result
	}
	for data := range warning.order {
		result = append(result, data)
	}
	return result
}
func (warning *Warning) SetSource(source string) {
	warning.source = source
}
func (warning *Warning) GetSource() string {
	return warning.source
}
