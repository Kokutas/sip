package sip

import (
	"encoding/json"
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
	User                string `json:"user"`
	TelephoneSubscriber string `json:"telephone-subscriber"`
	Password            string `json:"password,omitempty"`
}

func CreateUserInfo() *UserInfo {
	return &UserInfo{}
}
func NewUserInfo(user, telephoneSubscriber, password string) Sip {
	return &UserInfo{
		User:                user,
		TelephoneSubscriber: telephoneSubscriber,
		Password:            password,
	}
}
func (ui *UserInfo) Raw() string {
	result := ""
	if ui == nil {
		return result
	}
	switch {
	case len(strings.TrimSpace(ui.User)) > 0:
		result += ui.User
	case len(strings.TrimSpace(ui.TelephoneSubscriber)) > 0:
		result += ui.TelephoneSubscriber
	}
	if len(strings.TrimSpace(ui.Password)) > 0 {
		result += fmt.Sprintf(":%v", ui.Password)
	}
	return result
}
func (ui *UserInfo) JsonString() string {
	result := ""
	if ui == nil {
		return result
	}
	data, err := json.Marshal(ui)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (ui *UserInfo) Parser(raw string) error {
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if ui == nil {
		return errors.New("UserInfo caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	if regexp.MustCompile(`\:\w+$`).MatchString(raw) {
		passwordRaw := regexp.MustCompile(`\:\w+$`).FindString(raw)
		password := strings.TrimPrefix(passwordRaw, ":")
		ui.Password = password
		raw = strings.Replace(raw, passwordRaw, "", 1)
	}
	if regexp.MustCompile(`\d+\-\d+.*`).MatchString(raw) {
		ui.TelephoneSubscriber = raw
	} else {
		ui.User = raw
	}
	return nil
}
func (ui *UserInfo) Validator() error {
	if ui == nil {
		return errors.New("UserInfo caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(ui.User)) == 0 && len(strings.TrimSpace(ui.TelephoneSubscriber)) == 0 {
		return errors.New("user or telephone-subscriber must has one")
	}
	return nil
}
