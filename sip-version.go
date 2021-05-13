package sip

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SIP-Version    =  "SIP" "/" 1*DIGIT "." 1*DIGIT

type SipVersion struct {
	schema  string
	version float32
}

func (sv *SipVersion) Schema() string {
	return sv.schema
}

func (sv *SipVersion) SetSchema(schema string) {
	sv.schema = schema
}

func (sv *SipVersion) Version() float32 {
	return sv.version
}

func (sv *SipVersion) SetVersion(version float32) {
	sv.version = version
}
func NewSipVersion(schema string, version float32) *SipVersion {
	return &SipVersion{schema: schema, version: version}
}

func (sv *SipVersion) Raw() (string, error) {
	result := ""
	if err := sv.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(sv.schema)) == 0 {
		sv.schema = "SIP"
	}
	if sv.version > 0 {
		result += fmt.Sprintf("%v/%1.1f", strings.ToUpper(sv.schema), sv.version)
	}
	return result, nil
}
func (sv *SipVersion) String() string {
	result:=""
	if len(strings.TrimSpace(sv.schema))>0{
		result+=fmt.Sprintf("shcema: %s,",sv.schema)
	}
	if sv.version>0{
		result+=fmt.Sprintf("version: %1.1f,",sv.version)
	}
	result = strings.TrimSuffix(result,",")
	return result
}
func (sv *SipVersion) Parser(raw string) error {
	if sv == nil {
		return errors.New("sipVersion caller is not allowed to be nil")
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
	sv.schema = regexp.MustCompile(`/`).ReplaceAllString(regexp.MustCompile(`.*/`).FindString(raw), "")
	versions := regexp.MustCompile(`\d+.*`).FindString(raw)
	version, err := strconv.ParseFloat(versions, 32)
	if err != nil {
		return err
	}
	sv.version = float32(version)
	return nil
}
func (sv *SipVersion) Validator() error {
	if sv == nil {
		return errors.New("sipVersion caller is not allowed to be nil")
	}
	// sip-schema regexp
	sipSchemaRegexpStr := `(?i)(`
	for _, v := range Schemas {
		sipSchemaRegexpStr += v + "|"
	}
	sipSchemaRegexpStr = strings.TrimSuffix(sipSchemaRegexpStr, "|")
	sipSchemaRegexpStr += ")"
	sipSchemaRegexp:=regexp.MustCompile(sipSchemaRegexpStr)
	if len(strings.TrimSpace(sv.schema))==0{
		return  errors.New("schema is not allowed to be empty")
	}
	if !sipSchemaRegexp.MatchString(sv.schema){
		return errors.New("invalid schema")
	}
	return nil
}
