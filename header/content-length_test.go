package header

import (
	"fmt"
	"log"
	"testing"
)

func TestContentLength_Raw(t *testing.T) {
	cl := NewContentLength(0)
	fmt.Print(cl.Raw())
}

func TestContentLength_JsonString(t *testing.T) {
	cl := NewContentLength(0)
	fmt.Println(cl.JsonString())
}

func TestContentLength_Parser(t *testing.T) {
	raw := "Content-Length: 0\r\n"
	cl := CreateContentLength()
	if err := cl.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(cl.Raw())
	fmt.Println(cl.JsonString())
}

func TestContentLength_Validator(t *testing.T) {
	cl := NewContentLength(0)
	fmt.Println(cl.Validator())
}
