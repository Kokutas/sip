package header

import (
	"fmt"
	"log"
	"testing"
)

func TestContentLength_Field(t *testing.T) {

}

func TestContentLength_String(t *testing.T) {
	cl := NewContentLength(0)
	fmt.Println(cl.String())
}

func TestContentLength_Length(t *testing.T) {

}

func TestContentLength_Parser(t *testing.T) {
	raw := "Content-Length: 0\r\n"
	cl := new(ContentLength)
	if err := cl.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(cl.Raw())
	fmt.Println(cl.String())
}

func TestContentLength_Raw(t *testing.T) {
	cl := NewContentLength(0)
	fmt.Print(cl.Raw())
}

func TestContentLength_SetField(t *testing.T) {

}

func TestContentLength_SetLength(t *testing.T) {

}

func TestContentLength_Validator(t *testing.T) {
	cl := NewContentLength(0)
	fmt.Println(cl.Validator())
}

func TestNewContentLength(t *testing.T) {

}
