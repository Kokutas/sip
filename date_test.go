package sip

import (
	"fmt"
	"testing"
	"time"
)

func TestDate_Raw(t *testing.T) {
	dates := []*Date{
		NewDate(time.Kitchen, time.Now()),
		NewDate(time.ANSIC, time.Now()),
		NewDate("2006-01-02T15:04:05.000", time.Now()),
	}
	for _, date := range dates {
		fmt.Print(date.Raw())
	}
}

func TestDate_Parse(t *testing.T) {
	raws := []string{
		"Date: 2021-05-22T20:56:57.694\r\n",
		"Time: 2021-05-22T20:56:57.694\r\n",
		"Date: 9:18PM",
		"Date: Sat May 22 21:18:38 2021",
		"Date: 2021-05-22T21:18:38.352",
	}
	for index, raw := range raws {
		date := new(Date)
		date.Parse(raw)
		if len(date.GetSource()) > 0 {
			fmt.Print(index, " ", date.Raw())
			fmt.Println(index, "field", date.GetField(), ",format:", date.GetTimeFormat(), "time:", date.GetSipDate())
		}
	}
}
