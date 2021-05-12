package sip

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"sip/util"
	"strconv"
	"strings"
)

type Parameters struct {
	Transport string                 `json:"transport,omitempty"`
	User      string                 `json:"user,omitempty"`
	Method    string                 `json:"method,omitempty"`
	Ttl       uint8                  `json:"ttl,omitempty"`
	Maddr     string                 `json:"maddr,omitempty"`
	Lr        bool                   `json:"lr,omitempty"`
	Other     map[string]interface{} `json:"other,omitempty"`
}

func CreateParameters() Sip {
	return &Parameters{}
}
func NewParameters(transport, user, method string, ttl uint8, maddr string, lr bool, other map[string]interface{}) Sip {
	return &Parameters{
		Transport: transport,
		User:      user,
		Method:    method,
		Ttl:       ttl,
		Maddr:     maddr,
		Lr:        lr,
		Other:     other,
	}
}
func (ps *Parameters) Raw() string {
	result := ""
	if ps == nil {
		return result
	}
	if ps.Other != nil {
		for k, v := range ps.Other {
			if regexp.MustCompile(`(?i)(rport)`).Match([]byte(k)) {
				switch v.(type) {
				case int:
					if v.(int) > 0 {
						result += fmt.Sprintf(";rport=%v", v)
					} else {
						result += fmt.Sprintf(";rport")
					}
				case string:
					if len(strings.TrimSpace(v.(string))) > 0 {
						result += fmt.Sprintf(";rport=%v", v)
					} else {
						result += fmt.Sprintf(";rport")
					}
				default:
					result += fmt.Sprintf(";rport")
				}
				break
			}
		}
	}

	if len(strings.TrimSpace(ps.Transport)) > 0 {
		result += fmt.Sprintf(";transport=%v", strings.ToLower(ps.Transport))
	}
	if len(strings.TrimSpace(ps.User)) > 0 {
		result += fmt.Sprintf(";user=%v", ps.User)
	}
	if len(strings.TrimSpace(ps.Method)) > 0 {
		result += fmt.Sprintf(";method=%v", strings.ToLower(ps.Method))
	}
	if ps.Ttl > 0 {
		result += fmt.Sprintf(";ttl=%v", ps.Ttl)
	}
	if len(strings.TrimSpace(ps.Maddr)) > 0 {
		result += fmt.Sprintf(";maddr=%v", ps.Maddr)
	}
	if ps.Lr {
		result += ";lr"
	}
	received := ""
	if ps.Other != nil {
		others := ""
		for k, v := range ps.Other {
			if regexp.MustCompile(`(?i)(rport)`).Match([]byte(k)) {
				continue
			}
			if regexp.MustCompile(`(?i)(received)`).Match([]byte(k)) {
				switch v.(type) {
				case string:
					if len(strings.TrimSpace(v.(string))) > 0 {
						received = fmt.Sprintf(";received=%v", v)
					} else {
						received = fmt.Sprintf(";received")
					}
				case net.IP:
					if len(strings.TrimSpace(v.(net.IP).String())) > 0 {
						received = fmt.Sprintf(";received=%s", v.(net.IP).String())
					} else {
						received = fmt.Sprintf(";received")
					}
				default:
					received = fmt.Sprintf(";received")
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

	return result
}
func (ps *Parameters) JsonString() string {
	result := ""
	if ps == nil {
		return result
	}
	if data, err := json.Marshal(ps); err != nil {
		return result
	} else {
		result = fmt.Sprintf("%s", data)
	}
	return result
}
func (ps *Parameters) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if ps == nil {
		return errors.New("Parameters caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = strings.TrimPrefix(raw, ";")
	raw = strings.TrimSuffix(raw, ";")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
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
			transport = util.TrimPrefixAndSuffix(transport, " ")
			ps.Transport = transport
		case userRegex.MatchString(rawv):
			user := regexp.MustCompile(`(?i)(user=)`).ReplaceAllString(userRegex.FindString(rawv), "")
			ps.User = user
		case methodRegex.MatchString(rawv):
			method := regexp.MustCompile(`(?i)(method=)`).ReplaceAllString(methodRegex.FindString(rawv), "")
			ps.Method = method
		case ttlRegex.MatchString(rawv):
			ttlStr := regexp.MustCompile(`(?i)(ttl=)`).ReplaceAllString(ttlRegex.FindString(rawv), "")
			ttl, err := strconv.Atoi(ttlStr)
			if err != nil {
				return err
			}
			ps.Ttl = uint8(ttl)
		case maddrRegex.MatchString(rawv):
			maddr := regexp.MustCompile(`(?i)(maddr=)`).ReplaceAllString(maddrRegex.FindString(rawv), "")
			ps.Maddr = maddr
		case lrRegex.MatchString(rawv):
			raw = lrRegex.ReplaceAllString(raw, "")
			ps.Lr = true
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
				ps.Other = m
			}
		}
	}
	return nil
}
func (ps *Parameters) Validator() error {
	if ps == nil {
		return errors.New("Parameters caller is not allowed to be nil")
	}
	return nil
}
