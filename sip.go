package sip

import (
	"strings"
)

const (
	sip  = "sip"
	sips = "sips"
	tel  = "tel"
)

// Informational  =  "100"  ;  Trying
//               /   "180"  ;  Ringing
//               /   "181"  ;  Call Is Being Forwarded
//               /   "182"  ;  Queued
//               /   "183"  ;  Session Progress
var informational = map[int]string{
	100: "Trying",
	180: "Ringing",
	181: "Call Is Being Forwarded",
	182: "Queued",
	183: "Session Progress",
}

// Success  =  "200"  ;  OK
var success = map[int]string{
	200: "OK",
}

// Redirection  =  "300"  ;  Multiple Choices
//             /   "301"  ;  Moved Permanently
//             /   "302"  ;  Moved Temporarily
//             /   "305"  ;  Use Proxy
//             /   "380"  ;  Alternative Service
var redirection = map[int]string{
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
var clientError = map[int]string{
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
var serverError = map[int]string{
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
var globalFailure = map[int]string{
	600: "Busy Everywhere",
	603: "Decline",
	604: "Does not exist anywhere",
	606: "Not Acceptable",
}

var schemas = map[string]string{
	sip:  sip,
	sips: sips,
	tel:  tel,
}
var methods = map[string]string{
	"ACK":       "ACK",
	"BYE":       "BYE",
	"CANCEL":    "CANCEL",
	"INFO":      "INFO",
	"INVITE":    "INVITE",
	"MESSAGE":   "MESSAGE",
	"NOTIFY":    "NOTIFY",
	"OPTIONS":   "OPTIONS",
	"REGISTER":  "REGISTER",
	"SUBSCRIBE": "SUBSCRIBE",
}

type SipLayer interface {
	Raw() string
	Parse()
}

func stringTrimPrefixAndTrimSuffix(source string, sub string) string {
	for strings.HasPrefix(source, sub) || strings.HasSuffix(source, sub) {
		source = strings.TrimPrefix(source, sub)
		source = strings.TrimSuffix(source, sub)
	}
	return source
}

// type generic struct {
// 	index int // order
// 	kv    map[int]map[string]interface{}
// 	gk    sync.RWMutex
// }

// func (g *generic) store(k string, v interface{}) {
// 	g.gk.Lock()
// 	defer g.gk.Unlock()
// 	g.index++
// 	g.kv[g.index] = map[string]interface{}{k: v}
// }
