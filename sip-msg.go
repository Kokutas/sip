package sip

import "strings"

type SipMsg struct {
	*RequestLine
	*StatusLine
	*Authorization
	*CallID
	*Contact
	*ContentLength
	*ContentType
	*CSeq
	*Date
	*Expires
	*From
	*MaxForwards
	*Route
	*Subject
	*To
	*UserAgent
	*Via
	*Warning
	*WWWAuthenticate
	isOrder bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order   chan string // It is convenient to record the order of the original parameter fields when parsing
	source  string      // source string
}

func (sm *SipMsg) SetRequestLine(requestLine *RequestLine) {
	sm.RequestLine = requestLine
}
func (sm *SipMsg) GetRequestLine() *RequestLine {
	return sm.RequestLine
}
func (sm *SipMsg) SetStatusLine(statusLine *StatusLine) {
	sm.StatusLine = statusLine
}
func (sm *SipMsg) GetStatusLine() *StatusLine {
	return sm.StatusLine
}
func (sm *SipMsg) SetAuthorization(authorization *Authorization) {
	sm.Authorization = authorization
}
func (sm *SipMsg) GetAuthorization() *Authorization {
	return sm.Authorization
}
func (sm *SipMsg) SetCallID(callId *CallID) {
	sm.CallID = callId
}
func (sm *SipMsg) GetCallID() *CallID {
	return sm.CallID
}
func (sm *SipMsg) SetContact(contact *Contact) {
	sm.Contact = contact
}
func (sm *SipMsg) GetContact() *Contact {
	return sm.Contact
}
func (sm *SipMsg) SetContentLength(contentLength *ContentLength) {
	sm.ContentLength = contentLength
}
func (sm *SipMsg) GetContentLength() *ContentLength {
	return sm.ContentLength
}
func (sm *SipMsg) SetContentType(contentType *ContentType) {
	sm.ContentType = contentType
}
func (sm *SipMsg) GetContentType() *ContentType {
	return sm.ContentType
}
func (sm *SipMsg) SetCSeq(cseq *CSeq) {
	sm.CSeq = cseq
}
func (sm *SipMsg) GetCSeq() *CSeq {
	return sm.CSeq
}
func (sm *SipMsg) SetDate(date *Date) {
	sm.Date = date
}
func (sm *SipMsg) GetDate() *Date {
	return sm.Date
}
func (sm *SipMsg) SetExpires(expires *Expires) {
	sm.Expires = expires
}
func (sm *SipMsg) GetExpires() *Expires {
	return sm.Expires
}
func (sm *SipMsg) SetFrom(from *From) {
	sm.From = from
}
func (sm *SipMsg) GetFrom() *From {
	return sm.From
}
func (sm *SipMsg) SetMaxForwards(maxForwards *MaxForwards) {
	sm.MaxForwards = maxForwards
}
func (sm *SipMsg) GetMaxForwards() *MaxForwards {
	return sm.MaxForwards
}
func (sm *SipMsg) SetRoute(route *Route) {
	sm.Route = route
}
func (sm *SipMsg) GetRoute() *Route {
	return sm.Route
}
func (sm *SipMsg) SetSubject(subject *Subject) {
	sm.Subject = subject
}
func (sm *SipMsg) GetSubject() *Subject {
	return sm.Subject
}
func (sm *SipMsg) SetTo(to *To) {
	sm.To = to
}
func (sm *SipMsg) GetTo() *To {
	return sm.To
}
func (sm *SipMsg) SetUserAgent(userAgent *UserAgent) {
	sm.UserAgent = userAgent
}
func (sm *SipMsg) GetUserAgent() *UserAgent {
	return sm.UserAgent
}
func (sm *SipMsg) SetVia(via *Via) {
	sm.Via = via
}
func (sm *SipMsg) GetVia() *Via {
	return sm.Via
}
func (sm *SipMsg) SetWarning(warning *Warning) {
	sm.Warning = warning
}
func (sm *SipMsg) GetWarning() *Warning {
	return sm.Warning
}
func (sm *SipMsg) SetWWWAuthenticate(wwwAuthenticate *WWWAuthenticate) {
	sm.WWWAuthenticate = wwwAuthenticate
}
func (sm *SipMsg) GetWWWAuthenticate() *WWWAuthenticate {
	return sm.WWWAuthenticate
}
func (sm *SipMsg) GetSource() string {
	return sm.source
}
func NewSipMsg(requestLine *RequestLine, statusLine *StatusLine, authorization *Authorization, callId *CallID, contact *Contact, contentLength *ContentLength, contentType *ContentType, cseq *CSeq, date *Date, expires *Expires, from *From, maxForwards *MaxForwards, route *Route, subject *Subject, to *To, userAgent *UserAgent, via *Via, warning *Warning, wwwAuthenticate *WWWAuthenticate) *SipMsg {

	return &SipMsg{
		RequestLine:     requestLine,
		StatusLine:      statusLine,
		Authorization:   authorization,
		CallID:          callId,
		Contact:         contact,
		ContentLength:   contentLength,
		ContentType:     contentType,
		CSeq:            cseq,
		Date:            date,
		Expires:         expires,
		From:            from,
		MaxForwards:     maxForwards,
		Route:           route,
		Subject:         subject,
		To:              to,
		UserAgent:       userAgent,
		Via:             via,
		Warning:         warning,
		WWWAuthenticate: wwwAuthenticate,
		isOrder:         false,
	}
}
func CreateUacSipMsg(headerFields []string, parameters map[string]string) *SipMsg {
	sm := new(SipMsg)

	return sm
}
func CreateUasSipMsg(headerFields []string, parameters map[string]string) *SipMsg {
	sm := new(SipMsg)

	return sm
}
func (sm *SipMsg) Raw() (result strings.Builder) {
	if sm.RequestLine != nil {
		rl := sm.RequestLine.Raw()
		result.WriteString(rl.String())
	} else if sm.StatusLine != nil {
		sl := sm.StatusLine.Raw()
		result.WriteString(sl.String())
	}
	if sm.Via != nil {
		via := sm.Via.Raw()
		result.WriteString(via.String())
	}
	if sm.From != nil {
		from := sm.From.Raw()
		result.WriteString(from.String())
	}
	if sm.To != nil {
		to := sm.To.Raw()
		result.WriteString(to.String())
	}
	if sm.CallID != nil {
		callId := sm.CallID.Raw()
		result.WriteString(callId.String())
	}
	if sm.Contact != nil {
		contact := sm.Contact.Raw()
		result.WriteString(contact.String())
	}
	if sm.Route != nil {
		route := sm.Route.Raw()
		result.WriteString(route.String())
	}
	if sm.UserAgent != nil {
		userAgent := sm.UserAgent.Raw()
		result.WriteString(userAgent.String())
	}
	if sm.CSeq != nil {
		cseq := sm.CSeq.Raw()
		result.WriteString(cseq.String())
	}
	if sm.Expires != nil {
		expires := sm.Expires.Raw()
		result.WriteString(expires.String())
	}
	if sm.MaxForwards != nil {
		if sm.MaxForwards.GetForwards() == 0 {
			sm.MaxForwards.SetForwards(70)
		}
		maxForwards := sm.MaxForwards.Raw()
		result.WriteString(maxForwards.String())
	}
	contentLength := sm.ContentLength.Raw()
	result.WriteString(contentLength.String())

	if sm.WWWAuthenticate != nil {
		wwwAuthenticate := sm.WWWAuthenticate.Raw()
		result.WriteString(wwwAuthenticate.String())
	} else if sm.Authorization != nil {
		authorization := sm.Authorization.Raw()
		result.WriteString(authorization.String())
	}

	result.WriteString("\r\n")
	return
}
func (sm *SipMsg) Parse(raw string)       {}
func (sm *SipMsg) sipMsgOrder(raw string) {}
