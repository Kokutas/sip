package header

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestDate_Field(t *testing.T) {

}

func TestDate_Format(t *testing.T) {

}

func TestDate_Parser(t *testing.T) {
	raw := "Date: 2021-05-11T10:40:59.488\r\n"
	date := new(Date)
	if err := date.Parser(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(raw)
	fmt.Print(date.Raw())
}

func TestDate_Raw(t *testing.T) {
	date := NewDate(time.Now(), "2006-01-02T15:04:05.000")
	fmt.Println(time.Now()) // 2021-05-11 10:40:59.4883987 +0800 CST m=+0.004589201
	fmt.Print(date.Raw())   // Date: 2021-05-11T10:40:59.488
}

func TestDate_SetField(t *testing.T) {

}

func TestDate_SetFormat(t *testing.T) {

}

func TestDate_SetSipDate(t *testing.T) {

}

func TestDate_SipDate(t *testing.T) {

}

func TestDate_String(t *testing.T) {
	date := NewDate(time.Now(), "2006-01-02T15:04:05.000")
	fmt.Println(date.String())
}

func TestDate_Validator(t *testing.T) {
	date := NewDate(time.Now(), "2006-01-02T15:04:05.000")
	fmt.Println(date.Validator())
}

func TestNewDate(t *testing.T) {

}
