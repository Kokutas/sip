package sip

// Informational  =  "100"  ;  Trying
//               /   "180"  ;  Ringing
//               /   "181"  ;  Call Is Being Forwarded
//               /   "182"  ;  Queued
//               /   "183"  ;  Session Progress
var Informational = map[int]string{
	100: "Trying",
	180: "Ringing",
	181: "Call Is Being Forwarded",
	182: "Queued",
	183: "Session Progress",
}

// Success  =  "200"  ;  OK
var Success = map[int]string{
	200: "OK",
}

// Redirection  =  "300"  ;  Multiple Choices
//             /   "301"  ;  Moved Permanently
//             /   "302"  ;  Moved Temporarily
//             /   "305"  ;  Use Proxy
//             /   "380"  ;  Alternative Service
var Redirection = map[int]string{
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Moved Temporarily",
	305: "Use Proxy",
	380: "Alternative Service",
}

// Client-Error  =  "400"  ;  Bad Request
//              /   "401"  ;  Unauthorized
//              /   "402"  ;  Payment Required
//              /   "403"  ;  Forbidden
//              /   "404"  ;  Not Found
//              /   "405"  ;   Not Allowed
//              /   "406"  ;  Not Acceptable
//              /   "407"  ;  Proxy Authentication Required
//              /   "408"  ;  Request Timeout
//              /   "410"  ;  Gone
//              /   "413"  ;  Request Entity Too Large
//              /   "414"  ;  Request-URI Too Large
//              /   "415"  ;  Unsupported Media Type
//              /   "416"  ;  Unsupported URI Scheme
//              /   "420"  ;  Bad Extension
//              /   "421"  ;  Extension Required
//              /   "423"  ;  Interval Too Brief
//              /   "480"  ;  Temporarily not available
//              /   "481"  ;  Call Leg/Transaction Does Not Exist
//              /   "482"  ;  Loop Detected
//              /   "483"  ;  Too Many Hops
//              /   "484"  ;  Address Incomplete
//              /   "485"  ;  Ambiguous
//              /   "486"  ;  Busy Here
//              /   "487"  ;  Request Terminated
//              /   "488"  ;  Not Acceptable Here
//              /   "491"  ;  Request Pending
//              /   "493"  ;  Undecipherable
var ClientError = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: " Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	410: "Gone",
	413: "Request Entity Too Large",
	414: "Request-URI Too Large",
	415: "Unsupported Media Type",
	416: "Unsupported URI Scheme",
	420: "Bad Extension",
	421: "Extension Required",
	423: "Interval Too Brief",
	480: "Temporarily not available",
	481: "Call Leg/Transaction Does Not Exist",
	482: "Loop Detected",
	483: "Too Many Hops",
	484: "Address Incomplete",
	485: "Ambiguous",
	486: "Busy Here",
	487: "Request Terminated",
	488: "Not Acceptable Here",
	491: "Request Pending",
	493: "Undecipherable",
}

// Server-Error  =  "500"  ;  Internal Server Error
//              /   "501"  ;  Not Implemented
//              /   "502"  ;  Bad Gateway
//              /   "503"  ;  Service Unavailable
//              /   "504"  ;  Server Time-out
//              /   "505"  ;  SIP Version not supported
//              /   "513"  ;  Message Too Large
var ServerError = map[int]string{
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Server Time-out",
	505: "SIP Version not supported",
	513: "Message Too Large",
}

// Global-Failure  =  "600"  ;  Busy Everywhere
// 			      /   "603"  ;  Decline
// 			      /   "604"  ;  Does not exist anywhere
// 			      /   "606"  ;  Not Acceptable
var GlobalFailure = map[int]string{
	600: "Busy Everywhere",
	603: "Decline",
	604: "Does not exist anywhere",
	606: "Not Acceptable",
}

// INVITEm           =  %x49.4E.56.49.54.45 ; INVITE in caps
// ACKm              =  %x41.43.4B ; ACK in caps
// OPTIONSm          =  %x4F.50.54.49.4F.4E.53 ; OPTIONS in caps
// BYEm              =  %x42.59.45 ; BYE in caps
// CANCELm           =  %x43.41.4E.43.45.4C ; CANCEL in caps
// REGISTERm         =  %x52.45.47.49.53.54.45.52 ; REGISTER in caps
//             =  INVITEm / ACKm / OPTIONSm / BYEm
//                      / CANCELm / REGISTERm
//                      / extension-
// extension-  =  token
const (
	ACK       = "ACK"
	BYE       = "BYE"
	CANCEL    = "CANCEL"
	INVITE    = "INVITE"
	MESSAGE   = "MESSAGE"
	NOTIFY    = "NOTIFY"
	OPTIONS   = "OPTIONS"
	REGISTER  = "REGISTER"
	SUBSCRIBE = "SUBSCRIBE"
)

var Methods = map[string]string{
	ACK:       ACK,
	BYE:       BYE,
	CANCEL:    CANCEL,
	INVITE:    INVITE,
	MESSAGE:   MESSAGE,
	NOTIFY:    NOTIFY,
	OPTIONS:   OPTIONS,
	REGISTER:  REGISTER,
	SUBSCRIBE: SUBSCRIBE,
}

const (
	SIP  = "sip"
	SIPS = "sips"
	TEL  = "tel"
)

var Schemas = map[string]string{
	SIP:  SIP,
	SIPS: SIPS,
	TEL:  TEL,
}

const (
	UDP  = "UDP"
	TCP  = "TCP"
	TLS  = "TLS"
	SCTP = "SCTP"
)

var Transports = map[string]string{
	UDP:  UDP,
	TCP:  TCP,
	TLS:  TLS,
	SCTP: SCTP,
}

const (
	Digest string = "Digest"
	Basic  string = "Basic"
)

type Sip interface {
	Raw() string
	JsonString() string
	Parser(raw string) error
	Validator() error
}
