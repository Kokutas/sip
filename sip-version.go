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
	Schema  string  `json:"schema"`
	Version float32 `json:"version"`
}

func CreateSipVersion() Sip {
	return &SipVersion{}
}
func NewSipVersion(schema string, version float32) Sip {
	return &SipVersion{Schema: schema, Version: version}
}
func (sv *SipVersion) Raw() string {
	result := ""
	if sv == nil {
		return result
	}
	if len(strings.TrimSpace(sv.Schema)) == 0 {
		sv.Schema = "SIP"
	}
	if sv.Version > 0 {
		result += fmt.Sprintf("%v/%1.1f", strings.ToUpper(sv.Schema), sv.Version)
	}
	return result
}
func (sv *SipVersion) JsonString() string {
	result := ""
	return result
}
func (sv *SipVersion) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimSuffix(raw, "\r")
	raw = strings.TrimSuffix(raw, "\n")
	raw = strings.TrimPrefix(raw, "\r")
	raw = strings.TrimPrefix(raw, "\n")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	if sv == nil {
		return errors.New("SipVersion caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	sv.Schema = regexp.MustCompile(`/`).ReplaceAllString(regexp.MustCompile(`.*/`).FindString(raw), "")
	versions := regexp.MustCompile(`\d+.*`).FindString(raw)
	version, err := strconv.ParseFloat(versions, 32)
	if err != nil {
		return err
	}
	sv.Version = float32(version)

	return nil
}
func (sv *SipVersion) Validator() error {
	if sv == nil {
		return errors.New("SipVersion caller is not allowed to be nil")
	}
	// sip-schema regexp
	sipSchemaRegexpStr := `(?i)(`
	for _, v := range Schemas {
		sipSchemaRegexpStr += v + "|"
	}
	sipSchemaRegexpStr = strings.TrimSuffix(sipSchemaRegexpStr, "|")
	sipSchemaRegexpStr += ")"
	sipVersionRegexp := regexp.MustCompile(sipSchemaRegexpStr + `/\d+\.\d+`)
	if !sipVersionRegexp.MatchString(sv.Raw()) {
		return errors.New("invalid sipVersion")
	}

	return nil
}
