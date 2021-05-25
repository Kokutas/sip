package sip

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// uri-parameters    =  *( ";" uri-parameter)
// uri-parameter     =  transport-param / user-param / method-param
//                      / ttl-param / maddr-param / lr-param / other-param
// transport-param   =  "transport="
//                      ( "udp" / "tcp" / "sctp" / "tls"
//                      / other-transport)
// other-transport   =  token
// user-param        =  "user=" ( "phone" / "ip" / other-user)
// other-user        =  token
// method-param      =  "method=" Method
// ttl-param         =  "ttl=" ttl
// maddr-param       =  "maddr=" host
// lr-param          =  "lr"
// other-param       =  pname [ "=" pvalue ]
// pname             =  1*paramchar
// pvalue            =  1*paramchar
// paramchar         =  param-unreserved / unreserved / escaped
// param-unreserved  =  "[" / "]" / "/" / ":" / "&" / "+" / "$"
//
type Parameters struct {
	transport string      // transport-param = "transport="( "udp" / "tcp" / "sctp" / "tls"/ other-transport),other-transport = token
	user      string      // user-param =  "user=" ( "phone" / "ip" / other-user), other-user = token
	method    string      // method-param =  "method=" Method
	ttl       uint8       // ttl-param =  "ttl=" ttl
	maddr     string      // maddr-param       =  "maddr=" host
	lr        bool        // lr-param          =  "lr"
	other     sync.Map    // other-param       =  pname [ "=" pvalue ]
	isOrder   bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order     chan string // It is convenient to record the order of the original parameter fields when parsing
	source    string      // source string
}

func (p *Parameters) SetTransport(transport string) {
	p.transport = transport
}
func (p *Parameters) GetTransport() string {
	return p.transport
}
func (p *Parameters) SetUser(user string) {
	p.user = user
}
func (p *Parameters) GetUser() string {
	return p.user
}
func (p *Parameters) SetMethod(method string) {
	p.method = method
}
func (p *Parameters) GetMethod() string {
	return p.method
}
func (p *Parameters) SetTtl(ttl uint8) {
	p.ttl = ttl
}
func (p *Parameters) GetTtl() uint8 {
	return p.ttl
}
func (p *Parameters) SetMaddr(maddr string) {
	p.maddr = maddr
}
func (p *Parameters) GetMaddr() string {
	return p.maddr
}
func (p *Parameters) SetLr(lr bool) {
	p.lr = lr
}
func (p *Parameters) GetLr() bool {
	return p.lr
}
func (p *Parameters) SetOther(other sync.Map) {
	p.other = other
}
func (p *Parameters) GetOther() sync.Map {
	return p.other
}
func (p *Parameters) GetSource() string {
	return p.source
}
func NewParameters(transport string, user string, method string, ttl uint8, maddr string, lr bool, other sync.Map) *Parameters {
	return &Parameters{
		transport: transport,
		user:      user,
		method:    method,
		ttl:       ttl,
		maddr:     maddr,
		lr:        lr,
		other:     other,
		isOrder:   false,
	}
}
func (p *Parameters) Raw() (result strings.Builder) {
	if p.isOrder {
		p.isOrder = false
		for orders := range p.order {
			if regexp.MustCompile(`((?i)(?:^transport))( )*=`).MatchString(orders) {
				// transport-param = "transport="( "udp" / "tcp" / "sctp" / "tls"/ other-transport),other-transport = token
				if len(strings.TrimSpace(p.transport)) > 0 {
					result.WriteString(fmt.Sprintf(";transport=%s", strings.ToLower(p.transport)))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(?:^user))( )*=`).MatchString(orders) {
				// user-param =  "user=" ( "phone" / "ip" / other-user), other-user = token
				if len(strings.TrimSpace(p.user)) > 0 {
					result.WriteString(fmt.Sprintf(";user=%s", p.user))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(?:^method))( )*=`).MatchString(orders) {
				// method-param =  "method=" Method
				if len(strings.TrimSpace(p.method)) > 0 {
					result.WriteString(fmt.Sprintf(";method=%s", p.method))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(?:^ttl))( )*=`).MatchString(orders) {
				// ttl-param =  "ttl=" ttl
				if p.ttl > 0 {
					result.WriteString(fmt.Sprintf(";ttl=%d", p.ttl))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(?:^maddr))( )*=`).MatchString(orders) {
				// maddr-param       =  "maddr=" host
				if len(strings.TrimSpace(p.maddr)) > 0 {
					result.WriteString(fmt.Sprintf(";maddr=%s", p.maddr))
				}
				continue
			}
			if regexp.MustCompile(`((?i)(?:^lr))( )*=`).MatchString(orders) {
				// lr-param          =  "lr"
				if p.lr {
					result.WriteString(";lr")
				}
				continue
			}
			ordersSlice := strings.Split(orders, "=")

			if len(ordersSlice) == 1 {
				if val, ok := p.other.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
					}
				} else {
					result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
				}
			} else {
				if val, ok := p.other.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
					}
				} else {
					if len(strings.TrimSpace(fmt.Sprintf("%v", ordersSlice[1]))) > 0 {
						result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], ordersSlice[1]))
					} else {
						result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
					}
				}
			}

		}

	} else {
		// transport-param = "transport="( "udp" / "tcp" / "sctp" / "tls"/ other-transport),other-transport = token
		if len(strings.TrimSpace(p.transport)) > 0 {
			result.WriteString(fmt.Sprintf(";transport=%s", strings.ToLower(p.transport)))
		}
		// user-param =  "user=" ( "phone" / "ip" / other-user), other-user = token
		if len(strings.TrimSpace(p.user)) > 0 {
			result.WriteString(fmt.Sprintf(";user=%s", p.user))
		}
		// method-param =  "method=" Method
		if len(strings.TrimSpace(p.method)) > 0 {
			result.WriteString(fmt.Sprintf(";method=%s", p.method))
		}
		// ttl-param =  "ttl=" ttl
		if p.ttl > 0 {
			result.WriteString(fmt.Sprintf(";ttl=%d", p.ttl))
		}
		// maddr-param       =  "maddr=" host
		if len(strings.TrimSpace(p.maddr)) > 0 {
			result.WriteString(fmt.Sprintf(";maddr=%s", p.maddr))
		}
		// lr-param          =  "lr"
		if p.lr {
			result.WriteString(";lr")
		}
	}

	// other     sync.Map
	p.other.Range(func(key, value interface{}) bool {
		if reflect.ValueOf(value).IsValid() {
			if reflect.ValueOf(value).IsZero() {
				result.WriteString(fmt.Sprintf(";%v", key))
				return true
			}
			result.WriteString(fmt.Sprintf(";%v=%v", key, value))
			return true
		}
		result.WriteString(fmt.Sprintf(";%v", key))
		return true
	})
	return
}
func (p *Parameters) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	p.source = raw
	// parameters other order
	p.otherOrder(raw)
	// tranport parameter regexp
	transportRegexp := regexp.MustCompile(`(?i)(transport).*`)
	// user parameter regexp
	userRegexp := regexp.MustCompile(`(?i)(user).*`)
	// method parameter regexp
	methodRegexp := regexp.MustCompile(`(?i)(method).*`)
	// ttl parameter regexp
	ttlRegexp := regexp.MustCompile(`(?i)(ttl).*`)
	// maddr parameter regexp
	maddrRegexp := regexp.MustCompile(`(?i)(maddr).*`)
	// lr parameterregexp
	lrRegexp := regexp.MustCompile(`(?i)(lr).*`)

	rawSlice := strings.Split(raw, ";")
	for _, raws := range rawSlice {
		switch {
		case transportRegexp.MatchString(raws):
			transport := regexp.MustCompile(`(?i)(transport)`).ReplaceAllString(raws, "")
			transport = regexp.MustCompile(`.*=`).ReplaceAllString(transport, "")
			transport = stringTrimPrefixAndTrimSuffix(transport, " ")
			if len(transport) > 0 {
				p.transport = transport
			}
		case userRegexp.MatchString(raws):
			user := regexp.MustCompile(`(?i)(user)`).ReplaceAllString(raws, "")
			user = regexp.MustCompile(`.*=`).ReplaceAllString(user, "")
			user = stringTrimPrefixAndTrimSuffix(user, " ")
			if len(user) > 0 {
				p.user = user
			}
		case methodRegexp.MatchString(raws):
			method := regexp.MustCompile(`(?i)(method)`).ReplaceAllString(raws, "")
			method = regexp.MustCompile(`.*=`).ReplaceAllString(method, "")
			method = stringTrimPrefixAndTrimSuffix(method, " ")
			if len(method) > 0 {
				p.method = method
			}
		case ttlRegexp.MatchString(raws):
			ttlStr := regexp.MustCompile(`(?i)(ttl)`).ReplaceAllString(raws, "")
			ttlStr = regexp.MustCompile(`.*=`).ReplaceAllString(ttlStr, "")
			ttlStr = stringTrimPrefixAndTrimSuffix(ttlStr, " ")
			if len(ttlStr) > 0 {
				ttl, _ := strconv.Atoi(ttlStr)
				p.ttl = uint8(ttl)
			}
		case maddrRegexp.MatchString(raws):
			maddr := regexp.MustCompile(`(?i)(maddr)`).ReplaceAllString(raws, "")
			maddr = regexp.MustCompile(`.*=`).ReplaceAllString(maddr, "")
			maddr = stringTrimPrefixAndTrimSuffix(maddr, " ")
			if len(maddr) > 0 {
				p.maddr = maddr
			}
		case lrRegexp.MatchString(raws):
			p.lr = true
		default:
			if len(strings.TrimSpace(raws)) > 0 {
				if strings.Contains(raws, "=") {
					gs := strings.Split(raws, "=")
					if len(gs) > 1 {
						p.other.Store(gs[0], gs[1])
					} else {
						p.other.Store(gs[0], "")
					}
				} else {
					p.other.Store(raws, "")
				}
			}
		}
	}
}
func (p *Parameters) otherOrder(parameter string) {
	p.isOrder = true
	p.order = make(chan string, 1024)
	defer close(p.order)
	parameter = stringTrimPrefixAndTrimSuffix(parameter, ";")
	parameter = stringTrimPrefixAndTrimSuffix(parameter, " ")
	parameters := strings.Split(parameter, ";")
	for _, data := range parameters {
		p.order <- data
	}
}
