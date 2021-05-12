package sip

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewUserInfo(t *testing.T) {
	ui := NewUserInfo("34020000001320000001", "", "")
	data, err := json.Marshal(ui)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func Test_UserInfo_Raw(t *testing.T) {
	ui := NewUserInfo("34020000001320000001", "", "")
	fmt.Println(ui.Raw())
}

func Test_UserInfo_JsonString(t *testing.T) {
	ui := NewUserInfo("34020000001320000001", "", "")
	if res := ui.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func Test_UserInfo_Parser(t *testing.T) {
	raw := "+358-555-1234567"
	ui := CreateUserInfo()
	if err := ui.Parser(raw); err != nil {
		log.Fatal(err)
	}
	if res := ui.JsonString(); res != "" {
		fmt.Println(res)
	}
}

func Test_UserInfo_Validator(t *testing.T) {
	ui := NewUserInfo("", "", "")
	if err := ui.Validator(); err != nil {
		log.Fatal(err)
	}
}
