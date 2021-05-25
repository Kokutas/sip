package sip

import (
	"fmt"
	"regexp"
	"strings"
)

// userinfo         =  ( user / telephone-subscriber ) [ ":" password ] "@"
// user             =  1*( unreserved / escaped / user-unreserved )
// user-unreserved  =  "&" / "=" / "+" / "$" / "," / ";" / "?" / "/"
// password         =  *( unreserved / escaped /
// 					"&" / "=" / "+" / "$" / "," )
//
// https://www.rfc-editor.org/rfc/rfc2806
//
// telephone-subscriber  = global-phone-number / local-phone-number
// global-phone-number   = "+" base-phone-number [isdn-subaddress]
//                         [post-dial] *(area-specifier /
//                         service-provider / future-extension)
// base-phone-number     = 1*phonedigit
// local-phone-number    = 1*(phonedigit / dtmf-digit /
//                         pause-character) [isdn-subaddress]
//                         [post-dial] area-specifier
//                         *(area-specifier / service-provider /
//                         future-extension)
// isdn-subaddress       = ";isub=" 1*phonedigit
// post-dial             = ";postd=" 1*(phonedigit /
//                         dtmf-digit / pause-character)
// area-specifier        = ";" phone-context-tag "=" phone-context-ident
// phone-context-tag     = "phone-context"
// phone-context-ident   = network-prefix / private-prefix
// network-prefix        = global-network-prefix / local-network-prefix
// global-network-prefix = "+" 1*phonedigit
// local-network-prefix  = 1*(phonedigit / dtmf-digit / pause-character)
// private-prefix        = (%x21-22 / %x24-27 / %x2C / %x2F / %x3A /
//                         %x3C-40 / %x45-4F / %x51-56 / %x58-60 /
//                         %x65-6F / %x71-76 / %x78-7E)
//                         *(%x21-3A / %x3C-7E)
//                         ; Characters in URLs must follow escaping rules
//                         ; as explained in [RFC2396]
type UserInfo struct {
	user                string
	telephoneSubscriber string
	password            string
	source              string // source string
}

func (ui *UserInfo) SetUser(user string) {
	ui.user = user
}
func (ui *UserInfo) GetUser() string {
	return ui.user
}
func (ui *UserInfo) SetTelephoneSubscriber(telephoneSubscriber string) {
	ui.telephoneSubscriber = telephoneSubscriber
}
func (ui *UserInfo) GetTelephoneSubscriber() string {
	return ui.telephoneSubscriber
}
func (ui *UserInfo) SetPassword(password string) {
	ui.password = password
}
func (ui *UserInfo) GetPassword() string {
	return ui.password
}
func (ui *UserInfo) GetSource() string {
	return ui.source
}

func NewUserInfo(user string, telephoneSubscriber string, password string) *UserInfo {
	return &UserInfo{
		user:                user,
		telephoneSubscriber: telephoneSubscriber,
		password:            password,
	}
}
func (ui *UserInfo) Raw() (result strings.Builder) {
	switch {
	case len(strings.TrimSpace(ui.user)) > 0:
		result.WriteString(ui.user)
	case len(strings.TrimSpace(ui.telephoneSubscriber)) > 0:
		result.WriteString(ui.telephoneSubscriber)
	}
	if len(strings.TrimSpace(ui.password)) > 0 {
		result.WriteString(fmt.Sprintf(":%s", ui.password))
	}
	return result
}
func (ui *UserInfo) Parse(raw string) {
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return
	}
	ui.source = raw
	// password regexp
	passwordRegexp := regexp.MustCompile(`:.*(@)*?(;)*?(\?)*?`)
	if passwordRegexp.MatchString(raw) {
		password := regexp.MustCompile(`:`).ReplaceAllString(passwordRegexp.FindString(raw), "")
		password = regexp.MustCompile(`@.*`).ReplaceAllString(password, "")
		password = regexp.MustCompile(`;.*`).ReplaceAllString(password, "")
		password = regexp.MustCompile(`\?.*`).ReplaceAllString(password, "")
		ui.password = password
		raw = passwordRegexp.ReplaceAllString(raw, "")
	}
	raw = stringTrimPrefixAndTrimSuffix(raw, " ")
	// telephone-subscriber regexp
	// 1.global-phone-number
	// 2.local-phone-number
	telephoneSubscribeRegexp := regexp.MustCompile(`(^(\+)?(\d{1,3}\-)?\d+\-\d+(\-\d+)?)$|(^\+\d+)|(^\d{11}$)`)
	if len(strings.TrimSpace(raw)) > 0 {
		if telephoneSubscribeRegexp.MatchString(raw) {
			ui.telephoneSubscriber = telephoneSubscribeRegexp.FindString(raw)
		} else {
			ui.user = raw
		}
	}
}
