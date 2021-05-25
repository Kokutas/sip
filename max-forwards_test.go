package sip

import (
	"fmt"
	"testing"
)

func TestMaxForwards_Raw(t *testing.T) {
	maxForwards := NewMaxForwards(70)
	result := maxForwards.Raw()
	fmt.Print(result.String())
}

func TestMaxForwards_Parse(t *testing.T) {
	raw := "Max-Forwards: 70\r\n"
	maxForwards := new(MaxForwards)
	maxForwards.Parse(raw)
	if len(maxForwards.GetSource()) > 0 {
		maxForwards.SetForwards(5)
		result := maxForwards.Raw()
		fmt.Print(result.String())
	}
}
