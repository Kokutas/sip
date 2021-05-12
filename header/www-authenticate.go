package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sip"
	"sip/util"
	"strings"
)

// A WWW-Authenticate header field value contains an authentication
//  challenge. See Section 22.2 for further details on its usage.
//  Example:
//  	WWW-Authenticate: Digest realm="atlanta.com",
//  		domain="sip:boxesbybob.com", qop="auth",
//  		nonce="f84f1cec41e6cbe5aea9c8e88d359",
//  		opaque="", stale=FALSE, algorithm=MD5

// WWW-Authenticate = "WWW-Authenticate" HCOLON challenge
// extension-header = header-name HCOLON header-value
// header-name = token
// header-value = *(TEXT-UTF8char / UTF8-CONT / LWS)
// message-body = *OCTET

// RFC 2617
//If a server receives a request for an access-protected object, and an
//   acceptable Authorization header is not sent, the server responds with
//   a "401 Unauthorized" status code, and a WWW-Authenticate header as
//   per the framework defined above, which for the digest scheme is
//   utilized as follows:
//
//      challenge        =  "Digest" digest-challenge
//
//      digest-challenge  = 1#( realm | [ domain ] | nonce |
//                          [ opaque ] |[ stale ] | [ algorithm ] |
//                          [ qop-options ] | [auth-param] )
//
//
//      domain            = "domain" "=" <"> URI ( 1*SP URI ) <">
//      URI               = absoluteURI | abs_path
//      nonce             = "nonce" "=" nonce-value
//      nonce-value       = quoted-string
//      opaque            = "opaque" "=" quoted-string
//      stale             = "stale" "=" ( "true" | "false" )
//      algorithm         = "algorithm" "=" ( "MD5" | "MD5-sess" |
//                           token )
//      qop-options       = "qop" "=" <"> 1#qop-value <">
//      qop-value         = "auth" | "auth-int" | token

// Stale
//
//一个标志，用来指示客户端先前的请求因其nonce值过期而被拒绝。如果stale是TRUE（大小写敏感），客户端可能希望用新的加密回应重新进行请求，而不用麻烦用户提供新的用户名和口令。服务器端只有在收到的请求nonce值不合法，而该nonce对应的摘要（digest）是合法的情况下（即客户端知道正确的用户名/口令），才能将stale置成TRUE值。如果stale是FALSE或其它非TRUE值，或者其stale域不存在，说明用户名、口令非法，要求输入新的值。
//
//
//
//nc值
//
//在刷新注册请求中，十六进制请求计数器(nc)必须比前一次使用的时候要大，否则攻击者可以简单的使用同样的认证信息重放老的请求
//
//Authorization: Digest username="200",realm="123.com",cnonce="6b8b4567",nc=00000001,qop=auth,            uri="sip:10.100.125.17:5060",nonce="52dc0b3a-0353-4d15-af8d-7df1b92d8422",                      response="024d327559de8cbca5406b8ffe84354f",algorithm=MD5
//
//
//
//大概的思想是  为了防止每次生成的MD5摘要值都相同，那就使用随机数nonce值，使用了随机数之后，md5的摘要值就会变化了， 但是为了性能考虑等原因， 我们希望nonce值是有有效期的，不要每次都使用不同的nonce值。在这个有效期内， MD5值又不变了，那么就又引入nc值， nc的值表示使用这个nonce的次数，要求每次加1. 所以就给重放攻击加大难度了。

type WWWAuthenticate struct {
	Field      string                 `json:"field"`
	AuthSchema string                 `json:"auth-schema"` // basic / digest
	Realm      string                 `json:"realm"`
	Domain     *sip.SipUri            `json:"domain"`
	Nonce      string                 `json:"nonce"`
	Opaque     string                 `json:"opaque"`
	Stale      string                 `json:"stale,string"`
	Algorithm  string                 `json:"algorithm"`
	QopOptions string                 `json:"qop-options"`
	AuthParam  map[string]interface{} `json:"auth-param"`
}

func CreateWWWAuthenticate() sip.Sip {
	return &WWWAuthenticate{}
}
func NewWWWAuthenticate(authSchema string, realm string, domain *sip.SipUri, nonce string, opaque string, stale string, algorithm string, qopOptions string, authParam map[string]interface{}) sip.Sip {
	return &WWWAuthenticate{
		Field:      "WWW-Authenticate",
		AuthSchema: authSchema,
		Realm:      realm,
		Domain:     domain,
		Nonce:      nonce,
		Opaque:     opaque,
		Stale:      stale,
		Algorithm:  algorithm,
		QopOptions: qopOptions,
		AuthParam:  authParam,
	}
}
func (wa *WWWAuthenticate) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, wa) {
		return result
	}
	result += fmt.Sprintf("%v:", strings.Title(wa.Field))
	if len(strings.TrimSpace(wa.AuthSchema)) == 0 {
		wa.AuthSchema = "Digest"
	}
	result += fmt.Sprintf(" %v ", strings.Title(wa.AuthSchema))

	if len(strings.TrimSpace(wa.Realm)) > 0 {
		result += fmt.Sprintf("realm=\"%v\",", wa.Realm)
	}

	if wa.Domain != nil {
		result += fmt.Sprintf("domain=\"%v\",", wa.Domain.Raw())
	}
	if len(strings.TrimSpace(wa.Nonce)) > 0 {
		result += fmt.Sprintf("nonce=\"%v\",", wa.Nonce)
	}
	if len(strings.TrimSpace(wa.Opaque)) > 0 {
		result += fmt.Sprintf("opaque=\"%v\",", wa.Opaque)
	}
	if regexp.MustCompile(`(?i)(true|false)`).MatchString(wa.Stale) {
		result += fmt.Sprintf("stale=%v,", strings.ToUpper(wa.Stale))
	}
	if len(strings.TrimSpace(wa.QopOptions)) > 0 {
		result += fmt.Sprintf("qop=\"%v\",", wa.QopOptions)
	}
	if wa.AuthParam != nil {
		for k, v := range wa.AuthParam {
			result += fmt.Sprintf("%v=\"%v\",", k, v)
		}
	}
	if len(strings.TrimSpace(wa.Algorithm)) > 0 {
		result += fmt.Sprintf("algorithm=%v", strings.ToUpper(wa.Algorithm))
	}
	result += "\r\n"
	return result
}
func (wa *WWWAuthenticate) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, wa) {
		return result
	}
	data, err := json.Marshal(wa)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (wa *WWWAuthenticate) Parser(raw string) error {
	if wa == nil {
		return errors.New("www-authenticate caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")

	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(www-authenticate).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		wa.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// auth schema regexp
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if authSchemaRegexp.MatchString(raw) {
		wa.AuthSchema = authSchemaRegexp.FindString(raw)
		raw = authSchemaRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}

	// realm regexp
	realmRegexp := regexp.MustCompile(`(?i)(realm).*?=.*`)
	// domain regexp
	domainRegexp := regexp.MustCompile(`(?i)(domain).*?=.*`)
	// nonce regexp
	nonceRegexp := regexp.MustCompile(`(?i)(nonce).*?=.*`)
	// opaque regexp
	opaqueRegexp := regexp.MustCompile(`(?i)(opaque).*?=.*`)
	// stale regexp
	staleRegexp := regexp.MustCompile(`(?i)(stale).*?=.*`)
	// algorithm regexp
	algorithmRegexp := regexp.MustCompile(`(?i)(algorithm).*?=.*`)
	// qop-options regexp
	qopOptionsRegexp := regexp.MustCompile(`(?i)(qop).*?=.*`)
	// auth-param regexp
	authParams := make(map[string]interface{})

	raw = util.TrimPrefixAndSuffix(raw, " ")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		switch {
		case realmRegexp.MatchString(raws):
			realms := realmRegexp.FindString(raws)
			realms = regexp.MustCompile(`(?i)(realm).*?=`).ReplaceAllString(realms, "")
			realms = regexp.MustCompile(`"`).ReplaceAllString(realms, "")
			wa.Realm = util.TrimPrefixAndSuffix(realms, " ")
		case domainRegexp.MatchString(raws):
			domains := domainRegexp.FindString(raws)
			domains = regexp.MustCompile(`(?i)(domain).*?=`).ReplaceAllString(domains, "")
			domains = regexp.MustCompile(`"`).ReplaceAllString(domains, "")
			domains = util.TrimPrefixAndSuffix(domains, " ")
			wa.Domain = sip.CreateSipUri().(*sip.SipUri)
			if err := wa.Domain.Parser(domains); err != nil {
				return err
			}
		case nonceRegexp.MatchString(raws):
			nonces := nonceRegexp.FindString(raws)
			nonces = regexp.MustCompile(`(?i)(nonce).*?=`).ReplaceAllString(nonces, "")
			nonces = regexp.MustCompile(`"`).ReplaceAllString(nonces, "")
			wa.Nonce = util.TrimPrefixAndSuffix(nonces, " ")
		case opaqueRegexp.MatchString(raws):
			opaques := opaqueRegexp.FindString(raws)
			opaques = regexp.MustCompile(`(?i)(opaque).*?=`).ReplaceAllString(opaques, "")
			opaques = regexp.MustCompile(`"`).ReplaceAllString(opaques, "")
			wa.Opaque = util.TrimPrefixAndSuffix(opaques, " ")
		case staleRegexp.MatchString(raws):
			stales := staleRegexp.FindString(raws)
			stales = regexp.MustCompile(`(?i)(stale).*?=`).ReplaceAllString(stales, "")
			stales = regexp.MustCompile(`"`).ReplaceAllString(stales, "")
			wa.Stale = util.TrimPrefixAndSuffix(stales, " ")
		case algorithmRegexp.MatchString(raws):
			algorithms := algorithmRegexp.FindString(raws)
			algorithms = regexp.MustCompile(`(?i)(algorithm).*?=`).ReplaceAllString(algorithms, "")
			algorithms = regexp.MustCompile(`"`).ReplaceAllString(algorithms, "")
			wa.Algorithm = util.TrimPrefixAndSuffix(algorithms, " ")
		case qopOptionsRegexp.MatchString(raws):
			qopOptions := qopOptionsRegexp.FindString(raws)
			qopOptions = regexp.MustCompile(`(?i)(qop).*?=`).ReplaceAllString(qopOptions, "")
			qopOptions = regexp.MustCompile(`"`).ReplaceAllString(qopOptions, "")
			wa.QopOptions = util.TrimPrefixAndSuffix(qopOptions, " ")
		default:
			// auth-param
			if strings.Contains(raws, "=") {
				rs := strings.Split(raws, "=")
				if len(rs) > 1 {
					authParams[rs[0]] = rs[1]
				}
			}
		}
		wa.AuthParam = authParams
	}
	if len(strings.TrimSpace(wa.Algorithm)) == 0 {
		wa.Algorithm = "MD5"
	}
	return nil
}
func (wa *WWWAuthenticate) Validator() error {
	if wa == nil {
		return errors.New("www-authenticate caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(wa.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(www-authenticate)`).Match([]byte(wa.Field)) {
		return errors.New("field is not match")
	}
	if wa.Domain != nil {
		if err := wa.Domain.Validator(); err != nil {
			return err
		}
	}

	return nil
}
