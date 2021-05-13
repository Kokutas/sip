package header

import (
	"fmt"
	"log"
	"testing"
)

func TestMaxForwards_Field(t *testing.T) {

}

func TestMaxForwards_Forwards(t *testing.T) {
}

func TestMaxForwards_Parser(t *testing.T) {
	raw := "Max-Forwards: 70\r\n"
	mf := new(MaxForwards)
	if err := mf.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(mf.Raw())
	fmt.Println(mf.String())
}

func TestMaxForwards_Raw(t *testing.T) {
	mf := NewMaxForwards(70)
	fmt.Print(mf.Raw())
}

func TestMaxForwards_SetField(t *testing.T) {
}

func TestMaxForwards_SetForwards(t *testing.T) {
}

func TestMaxForwards_String(t *testing.T) {
	mf := NewMaxForwards(0)
	fmt.Println(mf.String())
}

func TestMaxForwards_Validator(t *testing.T) {
	mf := NewMaxForwards(70)
	fmt.Println(mf.Validator())
}

func TestNewMaxForwards(t *testing.T) {
}
