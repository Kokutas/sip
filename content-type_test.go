package sip

import (
	"fmt"
	"sync"
	"testing"
)

func TestContentType_Raw(t *testing.T) {
	var parameter sync.Map
	parameter.Store("protocol", "application/pkcs7-signature")
	parameter.Store("micalg", "sha1")
	parameter.Store("boundary", "boundary42")
	c := NewContentType("multipart", "signed", parameter)
	result := c.Raw()
	fmt.Println(result.String())
}

func TestContentType_Parse(t *testing.T) {
	raws := []string{
		`Content-Type: multipart/signed;protocol="application/pkcs7-signature";micalg=sha1;boundary=boundary42\r\n`,
		"Content-Type: application/sdp",
		"c: text/html; charset=ISO-8859-4",
	}
	for index, raw := range raws {
		c := new(ContentType)
		c.Parse(raw)
		if len(c.GetSource()) > 0 {
			fmt.Print("index: ", index, ",field: ", c.GetField(), ",m-type: ", c.GetMType(), ",m-subtype: ", c.GetMSubType())
			parameter := c.GetParameter()
			parameter.Range(func(key, value interface{}) bool {
				fmt.Print(" ;", key, "=", value)
				return true
			})
			fmt.Println()
			parameter.Store("protocol", "hello/world")
			result := c.Raw()
			fmt.Println(index, result.String())
		}
	}
}
