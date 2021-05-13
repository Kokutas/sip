package main

import (
	"fmt"
	"sip"
	"sip/line"
)

func main() {
	v:=sip.NewSipVersion(sip.SIP, 2.0)
	sl := line.NewStatusLine(v, 604, sip.GlobalFailure[604])
	fmt.Print(sl.Raw())
}
