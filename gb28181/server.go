package gb28181

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kokutas/sip"
)

type Server struct {
	id        string
	realm     string
	ip        net.IP
	port      uint16
	transport string
	// conn net.Conn 改成发送和接收分离
}

func NewServer(id string, realm string, ip net.IP, port uint16, transport string) *Server {
	return &Server{
		id:        id,
		realm:     realm,
		ip:        ip,
		port:      port,
		transport: transport,
	}
}

// 暂时返回strings.Builder，后续直接发送出去
func (s *Server) Response(sm *sip.SipMsg) (result strings.Builder) {
	switch {
	case regexp.MustCompile(`(?i)(register)`).MatchString(sm.GetRequestLine().GetMethod()):

		// 判断是否能认证通过
		// 发起鉴权挑战
		{
			sm.SetStatusLine(sip.NewStatusLine("sip", 2.0, 401, sip.ClientError[401]))
			// nonce需要添加到数据库
			clientIP := net.IPv4(192, 168, 0, 108)
			nonce := sip.GenNonce(clientIP.String(), fmt.Sprintf("%v", time.Now().UnixNano()))
			sm.SetWWWAuthenticate(sip.NewWWWAuthenticate(s.realm, "", nonce, "", false, "MD5", "", sync.Map{}))
			// NOTICE : register的from 和to的uri部分不做修改，需要处理的是to tag
			sm.GetTo().SetTag(fmt.Sprintf("%v", time.Now().UnixNano()))
			// 修改User-Agent
			if sm.GetUserAgent() != nil {
				sm.GetUserAgent().SetServer("SIP", "UAS", "com.kokutas", "V1.0.0")
			}
			// 修改via
			if sm.GetVia() != nil {
				if sm.GetVia().GetRport() != 0 {
					// TODO : 从 conn 中获取
					sm.GetVia().SetRport(8899)
					sm.GetVia().SetReceived(net.IPv4(192, 168, 0, 255).String())
				}
			}
		}

		// Digest鉴权挑战n次（第一次不算）失败，判断branch是否一样（403），判断response计算不一致（403）--冻结
		// sm.SetStatusLine(sip.NewStatusLine("sip", 2.0, 403, sip.ClientError[403]))
		// Digest鉴权认证通过
		// sm.SetStatusLine(sip.NewStatusLine("sip", 2.0, 200, sip.Success[200]))
	}
	sm.SetRequestLine(nil)
	res := sm.Raw()
	result.WriteString(res.String())
	return
}

func (s *Server) Start() {
	// 所有非200类的消息都要回复告知对方已经收到，不要重发（除了catalog的xml连续结构的）
	// sm.SetStatusLine(sip.NewStatusLine("sip", 2.0, 100, sip.Informational[100]))
}
