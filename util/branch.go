package util

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

// GenerateBranch branch参数的值必须用magic cookie”z9hG4bK”打头. 其它部分是对“To, From, Call-ID头域和Request-URI”按一定的算法加密后得到。 根据本标准产生的branch ID必须用”z9h64bK”开头。这7个字母是一个乱数cookie（定义成为7位的是为了保证旧版本的RFC2543实现不会产生这样的值），这样服务器收到请求之后，可以很方便的知道这个branch ID是否由本规范所产生的（就是说，全局唯一的）
func GenerateBranch(from, to, callId, reqUri string) string {
	rand.Seed(time.Now().UnixNano())
	result := fmt.Sprintf("%x",
		md5.Sum([]byte(fmt.Sprintf("%v%v%v%v%v", from, to, callId, reqUri, rand.Intn(60000)))))
	return "z9hG4bK-" + result
}
func GenerateUnixNanoBranch() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("z9hG4bK%x", md5.Sum([]byte(fmt.Sprintf("%v%v", time.Now().UnixNano(), rand.Intn(60000)))))
}
