package sip

import (
	"fmt"
	"testing"
)

func TestUserInfo_Raw(t *testing.T) {
	userInfo := NewUserInfo("34020000001320000001", "", "aap")
	result := userInfo.Raw()
	fmt.Println(result.String())
}

func TestUserInfo_Parse(t *testing.T) {
	raws := []string{
		"010-12345678",
		"010-1234567",
		"+12125551212",
		"+12125551212@phone2net.com",
		"+1-212-555-1212:1234@gateway.com;user=phone",
		"+086-13755969903",
		"17521500865:5060",
		"+17521500865:5060",
		"+086-17521500865:5060",
		"86-17521500865:5060",
		"13755969903:abcd@qq.com",
		"+13755969903:xyz",
		"+86-010-40020021",
		"86-010-40020021",
		"010-40020020",
		"+86-13523458056",
		"10-13523458056",
		"34020000001320000001",
		"34020000001320000001:i123",
		"+086-0559-6959003:kokutas@163.com",
		"sipabc:i123",
	}
	for _, raw := range raws {
		userinfo := new(UserInfo)
		userinfo.Parse(raw)
		if len(userinfo.GetSource()) > 0 {
			fmt.Println("user:", userinfo.GetUser(), ",telephone:", userinfo.GetTelephoneSubscriber(), ",password:", userinfo.GetPassword())
			result := userinfo.Raw()
			fmt.Println(result.String())
		}
	}

}
