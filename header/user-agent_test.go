package header

import (
	"fmt"
	"log"
	"testing"
)

func TestUserAgent_Raw(t *testing.T) {
	ua := NewUserAgent("UAS v1.0")
	fmt.Print(ua.Raw())
}

func TestUserAgent_JsonString(t *testing.T) {
	ua := NewUserAgent("UAS v1.0")
	fmt.Println(ua.JsonString())
}

func TestUserAgent_Parser(t *testing.T) {
	raw := "User-Agent: UAS v1.0\r\n"
	ua := CreateUserAgent()
	if err := ua.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(ua.Raw())
	fmt.Println(ua.JsonString())
}

func TestUserAgent_Validator(t *testing.T) {
	ua := NewUserAgent("uas")
	fmt.Println(ua.Validator())
}
