package sip

import (
	"fmt"
	"testing"
)

func TestMaxForwards_Raw(t *testing.T) {
	maxForwards := NewMaxForwards(70)
	fmt.Print(maxForwards.Raw())
}

func TestMaxForwards_Parse(t *testing.T) {
	raw := "Max-Forwards: 70\r\n"
	maxForwards := new(MaxForwards)
	maxForwards.Parse(raw)
	// maxForwards.SetForwards(5)
	fmt.Print(maxForwards.Raw())
}
