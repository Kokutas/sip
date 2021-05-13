package sip

import (
	"fmt"
	"log"
	"testing"
)

func TestNewUserInfo(t *testing.T) {
	ui := NewUserInfo("34020000001320000001", "", "")
	fmt.Printf("%s\r\n", ui)
}

func TestUserInfo_Parser(t *testing.T) {
	raw := "+358-555-1234567"
	ui := new(UserInfo)
	if err := ui.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.String())
}

func TestUserInfo_Password(t *testing.T) {

}

func TestUserInfo_Raw(t *testing.T) {
	ui := NewUserInfo("34020000001320000001", "", "")
	fmt.Println(ui.Raw())
}

func TestUserInfo_SetPassword(t *testing.T) {

}

func TestUserInfo_SetTelephoneSubscriber(t *testing.T) {

}

func TestUserInfo_SetUser(t *testing.T) {

}

func TestUserInfo_String(t *testing.T) {
	ui := NewUserInfo("34020000001320000001", "", "")
	fmt.Println(ui.String())
}

func TestUserInfo_TelephoneSubscriber(t *testing.T) {

}

func TestUserInfo_User(t *testing.T) {

}

func TestUserInfo_Validator(t *testing.T) {
	ui := NewUserInfo("123", "", "")
	if err := ui.Validator(); err != nil {
		log.Fatal(err)
	}
}
