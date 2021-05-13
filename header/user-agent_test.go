package header

import (
	"fmt"
	"log"
	"testing"
)

func TestNewUserAgent(t *testing.T) {

}

func TestUserAgent_Field(t *testing.T) {

}

func TestUserAgent_Parser(t *testing.T) {
	raw := "User-Agent: UAS v1.0\r\n"
	ua := new(UserAgent)
	if err := ua.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(ua.Raw())
	fmt.Println(ua.String())
}

func TestUserAgent_Raw(t *testing.T) {
	ua := NewUserAgent("UAS v1.0")
	fmt.Print(ua.Raw())
}

func TestUserAgent_Server(t *testing.T) {

}

func TestUserAgent_SetField(t *testing.T) {

}

func TestUserAgent_SetServer(t *testing.T) {

}

func TestUserAgent_String(t *testing.T) {
	ua := NewUserAgent("UAS v1.0")
	fmt.Println(ua.String())
}

func TestUserAgent_Validator(t *testing.T) {
	ua := NewUserAgent("uas")
	fmt.Println(ua.Validator())
}
