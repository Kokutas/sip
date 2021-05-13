package header

import (
	"fmt"
	"log"
	"testing"
)

func TestExpires_Field(t *testing.T) {

}

func TestExpires_Parser(t *testing.T) {
	raw := "Expires: 3600\r\n"
	expires := new(Expires)
	if err := expires.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(expires.Raw())
	fmt.Println(expires.String())
}

func TestExpires_Raw(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Print(expires.Raw())
}

func TestExpires_Seconds(t *testing.T) {

}

func TestExpires_SetField(t *testing.T) {

}

func TestExpires_SetSeconds(t *testing.T) {

}

func TestExpires_String(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Println(expires.String())
}

func TestExpires_Validator(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Println(expires.Validator())
	expires = new(Expires)
	fmt.Println(expires.Validator())
}

func TestNewExpires(t *testing.T) {

}
