package sip

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type NameAddr struct {
	schema    string    // sip/sips
	addr      *HostPort // host/ipv4/ipv6[port]
	parameter sync.Map
	isOrder   bool        // Determine whether the analysis is the result of the analysis and whether it is sorted during the analysis
	order     chan string // It is convenient to record the order of the original parameter fields when parsing
	source    string      // source string
}

func (na *NameAddr) SetSchema(schema string) {
	na.schema = schema
}
func (na *NameAddr) GetSchema() string {
	return na.schema
}
func (na *NameAddr) SetAddr(addr *HostPort) {
	na.addr = addr
}
func (na *NameAddr) GetAddr() *HostPort {
	return na.addr
}
func (na *NameAddr) SetParameter(parameter sync.Map) {
	na.parameter = parameter
}
func (na *NameAddr) GetParameter() sync.Map {
	return na.parameter
}
func (na *NameAddr) GetSource() string {
	return na.source
}
func NewNameAddr(schema string, addr *HostPort, parameter sync.Map) *NameAddr {
	return &NameAddr{
		schema:    schema,
		addr:      addr,
		parameter: parameter,
		isOrder:   false,
	}
}
func (na *NameAddr) Raw() (result strings.Builder) {
	if len(strings.TrimSpace(na.schema)) == 0 {
		na.schema = "sip"
	}
	result.WriteString(fmt.Sprintf("%s:", strings.ToLower(na.schema)))
	if na.addr != nil {
		addr := na.addr.Raw()
		result.WriteString(addr.String())
	}
	if na.isOrder {
		na.isOrder = false
		for orders := range na.order {
			ordersSlice := strings.Split(orders, "=")
			if len(ordersSlice) == 1 {
				if val, ok := na.parameter.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
					}

				} else {
					result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
				}
			} else {
				if val, ok := na.parameter.LoadAndDelete(ordersSlice[0]); ok {
					if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
						result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], val))
					} else {
						result.WriteString(fmt.Sprintf(";%v", ordersSlice[0]))
					}
				} else {
					result.WriteString(fmt.Sprintf(";%v=%v", ordersSlice[0], ordersSlice[1]))
				}
			}
		}
	}
	na.parameter.Range(func(key, value interface{}) bool {
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
func (na *NameAddr) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	// schema regexp
	schemasRegexpStr := `(?i)(`
	for _, v := range schemas {
		schemasRegexpStr += v + "|"
	}
	schemasRegexpStr = strings.TrimSuffix(schemasRegexpStr, "|")
	schemasRegexpStr += ")( )*:"
	schemaRegexp := regexp.MustCompile(schemasRegexpStr)
	if !schemaRegexp.MatchString(raw) {
		return
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	na.source = raw
	na.parameter = sync.Map{}
	na.addr = new(HostPort)
	schema := schemaRegexp.FindString(raw)
	raw = regexp.MustCompile(`.*`+schema).ReplaceAllString(raw, "")
	schema = stringTrimPrefixAndTrimSuffix(schema, ":")
	schema = stringTrimPrefixAndTrimSuffix(schema, " ")
	na.schema = schema
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// parameter regexp
	parameterRegexp := regexp.MustCompile(`;.*`)
	if parameterRegexp.MatchString(raw) {
		parameter := parameterRegexp.FindString(raw)
		raw = parameterRegexp.ReplaceAllString(raw, "")
		parameter = stringTrimPrefixAndTrimSuffix(parameter, ";")
		parameter = stringTrimPrefixAndTrimSuffix(parameter, " ")
		na.parameterOrder(parameter)
		rawSlice := strings.Split(parameter, ";")
		for _, raws := range rawSlice {
			if len(strings.TrimSpace(raws)) > 0 {
				if strings.Contains(raws, "=") {
					gs := strings.Split(raws, "=")
					if len(gs) > 1 {
						na.parameter.Store(gs[0], gs[1])
					} else {
						na.parameter.Store(gs[0], "")
					}
				} else {
					na.parameter.Store(raws, "")
				}
			}
		}

	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		na.addr.Parse(raw)
	}
}
func (na *NameAddr) parameterOrder(raw string) {
	na.isOrder = true
	na.order = make(chan string, 1024)
	defer close(na.order)
	raw = stringTrimPrefixAndTrimSuffix(raw, ";")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	rawSlice := strings.Split(raw, ";")
	for _, raws := range rawSlice {
		na.order <- raws
	}
}
