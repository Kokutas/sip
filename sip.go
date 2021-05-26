package sip

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"
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
	Raw() strings.Builder
	Parse(raw string)
}

func stringTrimPrefixAndTrimSuffix(source string, sub string) string {
	for strings.HasPrefix(source, sub) || strings.HasSuffix(source, sub) {
		source = strings.TrimPrefix(source, sub)
		source = strings.TrimSuffix(source, sub)
	}
	return source
}

// GenerateBranch branch参数的值必须用magic cookie”z9hG4bK”打头. 其它部分是对“To, From, Call-ID头域和Request-URI”按一定的算法加密后得到。 根据本标准产生的branch ID必须用”z9h64bK”开头。这7个字母是一个乱数cookie（定义成为7位的是为了保证旧版本的RFC2543实现不会产生这样的值），这样服务器收到请求之后，可以很方便的知道这个branch ID是否由本规范所产生的（就是说，全局唯一的）
func GenBranch(from, to, callId, reqUri string) string {
	rand.Seed(time.Now().UnixNano())
	result := fmt.Sprintf("%x",
		md5.Sum([]byte(fmt.Sprintf("%v%v%v%v%v", from, to, callId, reqUri, rand.Intn(60000)))))
	return "z9hG4bK-" + result
}
func GenUnixNanoBranch() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("z9hG4bK%x", md5.Sum([]byte(fmt.Sprintf("%v%v", time.Now().UnixNano(), rand.Intn(60000)))))
}
