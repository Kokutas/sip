package sip

// 注释掉的都要优化，参照Authorization
type SipMsg struct {
	// *RequestLine
	// *StatusLine
	*Authorization
	// *Via
	// *From
	// *To
	// *CSeq
	// *CallID
	// *Expires
	// *Date
	// *Warning
	// *WWWAuthenticate
	// *Route
	// *Contact
	// *ContentLength
	// *ContentType
	// *MaxForwards
}
