package sip

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type Parameters struct {
	transport string
	user      string
	method    string
	ttl       uint8
	maddr     string
	lr        bool
	other     map[string]interface{}
}

func (parameters *Parameters) Transport() string {
	return parameters.transport
}

func (parameters *Parameters) SetTransport(transport string) {
	parameters.transport = transport
}

func (parameters *Parameters) User() string {
	return parameters.user
}

func (parameters *Parameters) SetUser(user string) {
	parameters.user = user
}

func (parameters *Parameters) Method() string {
	return parameters.method
}

func (parameters *Parameters) SetMethod(method string) {
	parameters.method = method
}

func (parameters *Parameters) Ttl() uint8 {
	return parameters.ttl
}

func (parameters *Parameters) SetTtl(ttl uint8) {
	parameters.ttl = ttl
}

func (parameters *Parameters) Maddr() string {
	return parameters.maddr
}

func (parameters *Parameters) SetMaddr(maddr string) {
	parameters.maddr = maddr
}

func (parameters *Parameters) Lr() bool {
	return parameters.lr
}

func (parameters *Parameters) SetLr(lr bool) {
	parameters.lr = lr
}

func (parameters *Parameters) Other() map[string]interface{} {
	return parameters.other
}

func (parameters *Parameters) SetOther(other map[string]interface{}) {
	parameters.other = other
}

func NewParameters(transport string, user string, method string, ttl uint8, maddr string, lr bool, other map[string]interface{}) *Parameters {
	return &Parameters{transport: transport, user: user, method: method, ttl: ttl, maddr: maddr, lr: lr, other: other}
}
func (parameters *Parameters) Raw() (string, error) {
	result := ""
	if err := parameters.Validator(); err != nil {
		return result, err
	}
	if parameters.other != nil {
		for k, v := range parameters.other {
			if regexp.MustCompile(`(?i)(rport)`).Match([]byte(k)) {
				switch v.(type) {
				case int:
					if v.(int) > 0 {
						result += fmt.Sprintf(";rport=%v", v)
					} else {
						result += ";rport"
					}
				case string:
					if len(strings.TrimSpace(v.(string))) > 0 {
						result += fmt.Sprintf(";rport=%v", v)
					} else {
						result += ";rport"
					}
				default:
					result += ";rport"
				}
				break
			}
		}
	}

	if len(strings.TrimSpace(parameters.transport)) > 0 {
		result += fmt.Sprintf(";transport=%v", strings.ToLower(parameters.transport))
	}
	if len(strings.TrimSpace(parameters.user)) > 0 {
		result += fmt.Sprintf(";user=%v", parameters.user)
	}
	if len(strings.TrimSpace(parameters.method)) > 0 {
		result += fmt.Sprintf(";method=%v", strings.ToLower(parameters.method))
	}
	if parameters.ttl > 0 {
		result += fmt.Sprintf(";ttl=%v", parameters.ttl)
	}
	if len(strings.TrimSpace(parameters.maddr)) > 0 {
		result += fmt.Sprintf(";maddr=%v", parameters.maddr)
	}
	if parameters.lr {
		result += ";lr"
	}
	received := ""
	if parameters.other != nil {
		others := ""
		for k, v := range parameters.other {
			if regexp.MustCompile(`(?i)(rport)`).Match([]byte(k)) {
				continue
			}
			if regexp.MustCompile(`(?i)(received)`).Match([]byte(k)) {
				switch v.(type) {
				case string:
					if len(strings.TrimSpace(v.(string))) > 0 {
						received = fmt.Sprintf(";received=%v", v)
					} else {
						received = ";received"
					}
				case net.IP:
					if len(strings.TrimSpace(v.(net.IP).String())) > 0 {
						received = fmt.Sprintf(";received=%s", v.(net.IP).String())
					} else {
						received = ";received"
					}
				default:
					received = ";received"
				}
				others += received
				continue
			}
			if len(strings.TrimSpace(v.(string))) == 0 {
				others += fmt.Sprintf(";%v", k)
			} else {
				others += fmt.Sprintf(";%v=%v", k, v)
			}
		}
		if len(others) > 0 {
			result += others
		}
	}

	return result, nil
}
func (parameters *Parameters) String() string {
	result := ""
	if len(strings.TrimSpace(parameters.transport)) > 0 {
		result += fmt.Sprintf("transport: %s,", parameters.transport)
	}
	if len(strings.TrimSpace(parameters.user)) > 0 {
		result += fmt.Sprintf("user: %s,", parameters.user)
	}
	if len(strings.TrimSpace(parameters.method)) > 0 {
	result+=fmt.Sprintf("method: %s,",parameters.method)
	}
	if parameters.ttl > 0 {
	result+=fmt.Sprintf("ttl: %v,",parameters.ttl)
	}
	if len(strings.TrimSpace(parameters.maddr)) > 0 {
		result+=fmt.Sprintf("maddr: %s,",parameters.maddr)
	}
	if parameters.lr {
		result+=fmt.Sprintf("lr: %v,",parameters.lr)
	}
	if parameters.other!=nil {
		result+=fmt.Sprintf("other: %v,",parameters.other)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (parameters *Parameters) Parser(raw string) error {
	if parameters == nil {
		return errors.New("parameters caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw," ")
	raw = strings.TrimPrefix(raw," ")
	raw = strings.TrimSuffix(raw," ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = strings.TrimLeft(raw, ";")
	raw = strings.TrimRight(raw,";")
	raw = strings.TrimPrefix(raw,";")
	raw = strings.TrimSuffix(raw,";")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw," ")
	raw = strings.TrimPrefix(raw," ")
	raw = strings.TrimSuffix(raw," ")
	rawSlice := strings.Split(raw, ";")

	transportRegx := regexp.MustCompile(`(?i)(transport=)\w*`)
	userRegex := regexp.MustCompile(`(?i)(user=)\w*`)
	methodRegex := regexp.MustCompile(`(?i)(method=)\w*`)
	ttlRegex := regexp.MustCompile(`(?i)(ttl=)\d*`)
	maddrRegex := regexp.MustCompile(`(?i)(maddr=).*\..*`)
	lrRegex := regexp.MustCompile(`(?i)(lr)\w*`)
	m := make(map[string]interface{})
	for _, rawv := range rawSlice {
		switch {
		case transportRegx.MatchString(rawv):
			transport := regexp.MustCompile(`(?i)(transport=)`).ReplaceAllString(transportRegx.FindString(rawv), "")
			transport = strings.TrimLeft(transport, " ")
			transport = strings.TrimRight(transport," ")
			transport = strings.TrimPrefix(transport," ")
			transport = strings.TrimSuffix(transport," ")
			parameters.transport = transport
		case userRegex.MatchString(rawv):
			user := regexp.MustCompile(`(?i)(user=)`).ReplaceAllString(userRegex.FindString(rawv), "")
			parameters.user = user
		case methodRegex.MatchString(rawv):
			method := regexp.MustCompile(`(?i)(method=)`).ReplaceAllString(methodRegex.FindString(rawv), "")
			parameters.method = method
		case ttlRegex.MatchString(rawv):
			ttlStr := regexp.MustCompile(`(?i)(ttl=)`).ReplaceAllString(ttlRegex.FindString(rawv), "")
			ttl, err := strconv.Atoi(ttlStr)
			if err != nil {
				return err
			}
			parameters.ttl = uint8(ttl)
		case maddrRegex.MatchString(rawv):
			maddr := regexp.MustCompile(`(?i)(maddr=)`).ReplaceAllString(maddrRegex.FindString(rawv), "")
			parameters.maddr = maddr
		case lrRegex.MatchString(rawv):
			raw = lrRegex.ReplaceAllString(raw, "")
			parameters.lr = true
		default:
			if strings.Contains(rawv, "=") {
				vs := strings.Split(rawv, "=")
				if len(vs) > 1 {
					m[vs[0]] = vs[1]
				} else {
					m[vs[0]] = ""
				}
			} else {
				m[rawv] = ""
			}
			if len(m) > 0 {
				parameters.other = m
			}
		}
	}
	return nil
}
func (parameters *Parameters) Validator() error {
	if parameters == nil {
		return errors.New("parameters caller is not allowed to be nil")
	}
	return nil
}
