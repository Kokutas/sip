package sip

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// userinfo         =  ( user / telephone-subscriber ) [ ":" password ] "@"
// user             =  1*( unreserved / escaped / user-unreserved )
// user-unreserved  =  "&" / "=" / "+" / "$" / "," / ";" / "?" / "/"
// password         =  *( unreserved / escaped /
//                     "&" / "=" / "+" / "$" / "," )

type UserInfo struct {
	user                string
	telephoneSubscriber string
	password            string
}

func (ui *UserInfo) User() string {
	return ui.user
}

func (ui *UserInfo) SetUser(user string) {
	ui.user = user
}

func (ui *UserInfo) TelephoneSubscriber() string {
	return ui.telephoneSubscriber
}

func (ui *UserInfo) SetTelephoneSubscriber(telephoneSubscriber string) {
	ui.telephoneSubscriber = telephoneSubscriber
}

func (ui *UserInfo) Password() string {
	return ui.password
}

func (ui *UserInfo) SetPassword(password string) {
	ui.password = password
}

func NewUserInfo(user string, telephoneSubscriber string, password string) *UserInfo {
	return &UserInfo{user: user, telephoneSubscriber: telephoneSubscriber, password: password}
}

func (ui *UserInfo) Raw() (string, error) {
	result := ""
	if err := ui.Validator(); err != nil {
		return result, err
	}
	switch {
	case len(strings.TrimSpace(ui.user)) > 0:
		result += ui.user
	case len(strings.TrimSpace(ui.telephoneSubscriber)) > 0:
		result += ui.telephoneSubscriber
	}
	if len(strings.TrimSpace(ui.password)) > 0 {
		result += fmt.Sprintf(":%s", ui.password)
	}
	return result, nil
}
func (ui *UserInfo) String() string {
	result:=""
	if len(strings.TrimSpace(ui.user))>0{
		result +=fmt.Sprintf("user: %s,",ui.user)
	}
	if len(strings.TrimSpace(ui.telephoneSubscriber))>0{
		result +=fmt.Sprintf("telephone-subscriber: %s,",ui.telephoneSubscriber)
	}
	if len(strings.TrimSpace(ui.password))>0{
		result +=fmt.Sprintf("password: %s,",ui.password)
	}
	result = strings.TrimSuffix(result,",")
	return result
}
func (ui *UserInfo) Parser(raw string) error {
	if ui == nil {
		return errors.New("userInfo caller is not allowed to be nil")
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
	if regexp.MustCompile(`:\w+$`).MatchString(raw) {
		passwordRaw := regexp.MustCompile(`:\w+$`).FindString(raw)
		password := strings.TrimPrefix(passwordRaw, ":")
		ui.password = password
		raw = strings.Replace(raw, passwordRaw, "", 1)
	}
	if regexp.MustCompile(`\d+-\d+.*`).MatchString(raw) {
		ui.telephoneSubscriber = raw
	} else {
		ui.user = raw
	}
	return nil
}
func (ui *UserInfo) Validator() error {
	if ui == nil {
		return errors.New("userInfo caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(ui.user)) == 0 && len(strings.TrimSpace(ui.telephoneSubscriber)) == 0 {
		return errors.New("user or telephone-subscriber must has one")
	}
	return nil
}
