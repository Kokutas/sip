package header

import (
	"fmt"
	"log"
	"testing"
)

func TestMaxForwards_Raw(t *testing.T) {
	mf := NewMaxForwards(70)
	fmt.Print(mf.Raw())
}

func TestMaxForwards_JsonString(t *testing.T) {
	mf := NewMaxForwards(0)
	fmt.Println(mf.JsonString())
}

func TestMaxForwards_Parser(t *testing.T) {
	raw := "Max-Forwards: 70\r\n"
	mf := CreateMaxForwards()
	if err := mf.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(mf.Raw())
	fmt.Println(mf.JsonString())
}

func TestMaxForwards_Validator(t *testing.T) {
	mf := NewMaxForwards(70)
	fmt.Println(mf.Validator())
}
