package gb28181

import (
	"crypto/md5"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kokutas/sip"
)

// GB/T 28181-2016 IPC
type IPC struct {
	id         string
	ip         net.IP
	port       uint16
	sid        string // server id
	sip        net.IP // server ip
	sport      uint16 // server port
	transport  string
	schema     string
	version    float64
	expires    uint32 // 注册过期时间
	registerSN uint32 // 注册cseq number
	// 发送间隔
	branch    string // branch值
	localId   string // call-id 中的local-id值
	fromTag   string // from tag 值
	toTag     string // to tag 值
	userAgent []string
	realm     string // Digest realm
	nonce     string // Digest nonce
}

func (ipc *IPC) SetExpires(expires uint32) {
	ipc.expires = expires
}
func (ipc *IPC) SetRegisterSN(sn uint32) {
	ipc.registerSN = sn
}
func (ipc *IPC) SetBranch(branch string) {
	ipc.branch = branch
}
func (ipc *IPC) SetLocalId(localId string) {
	ipc.localId = localId
}
func (ipc *IPC) SetFromTag(fromTag string) {
	ipc.fromTag = fromTag
}
func (ipc *IPC) SetToTag(toTag string) {
	ipc.toTag = toTag
}
func (ipc *IPC) SetUserAgent(userAgent ...string) {
	ipc.userAgent = userAgent
}
func (ipc *IPC) SetRealm(realm string) {
	ipc.realm = realm
}
func (ipc *IPC) SetNonce(nonce string) {
	ipc.nonce = nonce
}

func NewIPC(id string, ip net.IP, port uint16, sid string, sip net.IP, sport uint16, transport string, expires uint32) *IPC {
	return &IPC{
		id:         id,
		ip:         ip,
		port:       port,
		sid:        sid,
		sip:        sip,
		sport:      sport,
		transport:  transport,
		schema:     "sip",
		version:    2.0,
		expires:    expires,
		registerSN: 1,
		userAgent:  []string{"SIP", "UAC-IPC", "com.kokutas", "V1.0.0"},
	}
}
func (ipc *IPC) Request(method string, sm *sip.SipMsg) (result strings.Builder) {
	// sm.SetStatusLine(nil) TODO : 许多参数需要从这里拿到
	reqUri := sip.NewRequestUri(
		sip.NewSipUri(
			sip.NewUserInfo(ipc.id, "", ""),
			sip.NewHostPort("", ipc.sip, nil, ipc.sport),
			nil,
			sync.Map{}))
	reqLine := sip.NewRequestLine(method, reqUri, ipc.schema, ipc.version)
	// from tag
	if len(strings.TrimSpace(ipc.fromTag)) == 0 {
		ipc.fromTag = fmt.Sprintf("%v", time.Now().UnixNano())
	}
	from := sip.NewFrom("", "<", ipc.schema, ipc.id, ipc.ip.String(), ipc.port, ipc.fromTag, sync.Map{})
	to := sip.NewTo("", "<", ipc.schema, ipc.sid, ipc.sip.String(), ipc.sport, ipc.toTag, sync.Map{})
	if regexp.MustCompile(`(?i)(register)`).MatchString(method) {
		to = sip.NewTo("", "<", ipc.schema, ipc.id, ipc.ip.String(), ipc.port, "", sync.Map{})
	}
	contact := sip.NewContact("", "<", ipc.schema, ipc.id, ipc.ip.String(), ipc.port, "", -1, sync.Map{})
	// localId
	if len(strings.TrimSpace(ipc.localId)) == 0 {
		ipc.localId = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", time.Now().UnixNano()))))
	}
	callId := sip.NewCallID(ipc.localId, ipc.ip.String())
	// branch
	if len(strings.TrimSpace(ipc.branch)) == 0 {
		fromRaw := from.Raw()
		toRaw := to.Raw()
		callIdRaw := callId.Raw()
		fromVal := regexp.MustCompile(`(?i)(from|f)( ):`).ReplaceAllString(fromRaw.String(), "")
		toVal := regexp.MustCompile(`(?i)(to|t)( ):`).ReplaceAllString(toRaw.String(), "")
		callIdVal := regexp.MustCompile(`(?i)(call-id)( ):`).ReplaceAllString(callIdRaw.String(), "")
		reqUriVal := reqUri.Raw()
		ipc.branch = sip.GenBranch(fromVal, toVal, callIdVal, reqUriVal.String())
	}
	// via
	via := sip.NewVia(ipc.schema, ipc.version, ipc.transport, ipc.sip.String(), ipc.sport, 0, "", "", ipc.branch, 1, "", sync.Map{})
	expires := sip.NewExpires(ipc.expires)
	cSeq := sip.NewCSeq(ipc.registerSN, method)
	maxForwards := sip.NewMaxForwards(70)
	contentLength := sip.NewContentLength(0)
	userAgent := sip.NewUserAgent(ipc.userAgent...)
	sm.SetRequestLine(reqLine)
	sm.SetFrom(from)
	sm.SetTo(to)
	sm.SetContact(contact)
	sm.SetCallID(callId)
	sm.SetVia(via)
	sm.SetExpires(expires)
	sm.SetCSeq(cSeq)
	sm.SetUserAgent(userAgent)
	sm.SetMaxForwards(maxForwards)
	sm.SetContentLength(contentLength)
	if len(strings.TrimSpace(ipc.nonce)) > 0 {
		realm := ipc.realm
		if len(strings.TrimSpace(realm)) == 0 {
			realm = ipc.sid[:10]
		}
		reqUriRaw := reqUri.Raw()
		dp := &sip.DigestParams{
			Algorithm: "MD5",
			Method:    method,
			URI:       reqUriRaw.String(),
			Nonce:     ipc.nonce,
		}
		response := sip.GenDigestResponse(dp)
		authorization := sip.NewAuthorization(ipc.id, realm, ipc.nonce, reqUri, response, "MD5", "", "", "", "", sync.Map{})
		sm.SetAuthorization(authorization)
	}
	res := sm.Raw()
	result.WriteString(res.String())
	return
}

func (ipc *IPC) Response(code uint, reason string, sm *sip.SipMsg) (result strings.Builder) {
	sm.SetStatusLine(sip.NewStatusLine(sm.GetRequestLine().GetSchema(), sm.GetRequestLine().GetVersion(), code, reason))
	res := sm.Raw()
	result.WriteString(res.String())
	sm.SetRequestLine(nil) // 最后做，发送前/返回信息前
	return
}

func (ipc *IPC) Start() {}
