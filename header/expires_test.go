package header

import (
	"fmt"
	"log"
	"testing"
)

func TestExpires_JsonString(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Println(expires.JsonString())
}

func TestExpires_Raw(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Print(expires.Raw())
}

func TestExpires_Parser(t *testing.T) {
	raw := "Expires: 3600\r\n"
	expires := CreateExpires()
	if err := expires.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(expires.Raw())
	fmt.Println(expires.JsonString())
}

func TestExpires_Validator(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Println(expires.Validator())
	// expires = CreateCSeq()
	// fmt.Println(expires.Validator())
}
