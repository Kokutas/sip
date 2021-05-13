package header

import (
	"errors"
	"fmt"
	"github.com/kokutas/sip"
	"github.com/kokutas/sip/util"
	"regexp"
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
	field      string
	authSchema string // basic / digest
	realm      string
	domain     *sip.SipUri
	nonce      string
	opaque     string
	stale      string
	algorithm  string
	qopOptions string
	authParam  map[string]interface{}
}

func (wa *WWWAuthenticate) Field() string {
	return wa.field
}

func (wa *WWWAuthenticate) SetField(field string) {
	wa.field = field
}

func (wa *WWWAuthenticate) AuthSchema() string {
	return wa.authSchema
}

func (wa *WWWAuthenticate) SetAuthSchema(authSchema string) {
	wa.authSchema = authSchema
}

func (wa *WWWAuthenticate) Realm() string {
	return wa.realm
}

func (wa *WWWAuthenticate) SetRealm(realm string) {
	wa.realm = realm
}

func (wa *WWWAuthenticate) Domain() *sip.SipUri {
	return wa.domain
}

func (wa *WWWAuthenticate) SetDomain(domain *sip.SipUri) {
	wa.domain = domain
}

func (wa *WWWAuthenticate) Nonce() string {
	return wa.nonce
}

func (wa *WWWAuthenticate) SetNonce(nonce string) {
	wa.nonce = nonce
}

func (wa *WWWAuthenticate) Opaque() string {
	return wa.opaque
}

func (wa *WWWAuthenticate) SetOpaque(opaque string) {
	wa.opaque = opaque
}

func (wa *WWWAuthenticate) Stale() string {
	return wa.stale
}

func (wa *WWWAuthenticate) SetStale(stale string) {
	wa.stale = stale
}

func (wa *WWWAuthenticate) Algorithm() string {
	return wa.algorithm
}

func (wa *WWWAuthenticate) SetAlgorithm(algorithm string) {
	wa.algorithm = algorithm
}

func (wa *WWWAuthenticate) QopOptions() string {
	return wa.qopOptions
}

func (wa *WWWAuthenticate) SetQopOptions(qopOptions string) {
	wa.qopOptions = qopOptions
}

func (wa *WWWAuthenticate) AuthParam() map[string]interface{} {
	return wa.authParam
}

func (wa *WWWAuthenticate) SetAuthParam(authParam map[string]interface{}) {
	wa.authParam = authParam
}

func NewWWWAuthenticate(authSchema string, realm string, domain *sip.SipUri, nonce string, opaque string, stale string, algorithm string, qopOptions string, authParam map[string]interface{}) *WWWAuthenticate {
	return &WWWAuthenticate{field: "WWW-Authenticate", authSchema: authSchema, realm: realm, domain: domain, nonce: nonce, opaque: opaque, stale: stale, algorithm: algorithm, qopOptions: qopOptions, authParam: authParam}
}

func (wa *WWWAuthenticate) Raw() (string, error) {
	result := ""
	if err := wa.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(wa.field)) == 0 {
		wa.field = "WWW-Authenticate"
	}
	result += fmt.Sprintf("%s:", wa.field)
	if len(strings.TrimSpace(wa.authSchema)) == 0 {
		wa.authSchema = "Digest"
	}
	result += fmt.Sprintf(" %s", strings.Title(wa.authSchema))

	if len(strings.TrimSpace(wa.realm)) > 0 {
		result += fmt.Sprintf(" realm=\"%v\",", wa.realm)
	}

	if wa.domain != nil {
		res, err := wa.domain.Raw()
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf("domain=\"%v\",", res)
	}
	if len(strings.TrimSpace(wa.nonce)) > 0 {
		result += fmt.Sprintf("nonce=\"%v\",", wa.nonce)
	}
	if len(strings.TrimSpace(wa.opaque)) > 0 {
		result += fmt.Sprintf("opaque=\"%v\",", wa.opaque)
	}
	if regexp.MustCompile(`(?i)(true|false)`).MatchString(wa.stale) {
		result += fmt.Sprintf("stale=%v,", strings.ToUpper(wa.stale))
	}
	if len(strings.TrimSpace(wa.qopOptions)) > 0 {
		result += fmt.Sprintf("qop=\"%v\",", wa.qopOptions)
	}
	if wa.authParam != nil {
		for k, v := range wa.authParam {
			result += fmt.Sprintf("%v=\"%v\",", k, v)
		}
	}
	if len(strings.TrimSpace(wa.algorithm)) > 0 {
		result += fmt.Sprintf("algorithm=%v", strings.ToUpper(wa.algorithm))
	}
	result += "\r\n"
	return result, nil
}
func (wa *WWWAuthenticate) String() string {
	result := ""
	if len(strings.TrimSpace(wa.field)) > 0 {
		result += fmt.Sprintf("field: %s,", wa.field)
	}
	if len(strings.TrimSpace(wa.authSchema)) > 0 {
		result += fmt.Sprintf("auth-schema: %v,", wa.authSchema)
	}
	if len(strings.TrimSpace(wa.realm)) > 0 {
		result += fmt.Sprintf("realm: %v,", wa.realm)
	}
	if wa.domain != nil {
		result += fmt.Sprintf("%s", wa.domain.String())
	}
	if len(strings.TrimSpace(wa.nonce)) > 0 {
		result += fmt.Sprintf("nonce: %s,", wa.nonce)
	}
	if len(strings.TrimSpace(wa.opaque)) > 0 {
		result += fmt.Sprintf("opaque: %v,", wa.opaque)
	}
	if regexp.MustCompile(`(?i)(true|false)`).MatchString(wa.stale) {
		result += fmt.Sprintf("stale: %v,", strings.ToUpper(wa.stale))
	}
	if len(strings.TrimSpace(wa.qopOptions)) > 0 {
		result += fmt.Sprintf("qop: %v,", wa.qopOptions)
	}
	if wa.authParam != nil {
		result += fmt.Sprintf("auth-param: %v,", wa.authParam)
	}
	if len(strings.TrimSpace(wa.algorithm)) > 0 {
		result += fmt.Sprintf("algorithm:%v,", wa.algorithm)
	}
	result = strings.TrimSuffix(result, ",")
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
		wa.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// auth schema regexp
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if authSchemaRegexp.MatchString(raw) {
		wa.authSchema = authSchemaRegexp.FindString(raw)
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
			wa.realm = util.TrimPrefixAndSuffix(realms, " ")
		case domainRegexp.MatchString(raws):
			domains := domainRegexp.FindString(raws)
			domains = regexp.MustCompile(`(?i)(domain).*?=`).ReplaceAllString(domains, "")
			domains = regexp.MustCompile(`"`).ReplaceAllString(domains, "")
			domains = util.TrimPrefixAndSuffix(domains, " ")
			wa.domain = new(sip.SipUri)
			if err := wa.domain.Parser(domains); err != nil {
				return err
			}
		case nonceRegexp.MatchString(raws):
			nonces := nonceRegexp.FindString(raws)
			nonces = regexp.MustCompile(`(?i)(nonce).*?=`).ReplaceAllString(nonces, "")
			nonces = regexp.MustCompile(`"`).ReplaceAllString(nonces, "")
			wa.nonce = util.TrimPrefixAndSuffix(nonces, " ")
		case opaqueRegexp.MatchString(raws):
			opaques := opaqueRegexp.FindString(raws)
			opaques = regexp.MustCompile(`(?i)(opaque).*?=`).ReplaceAllString(opaques, "")
			opaques = regexp.MustCompile(`"`).ReplaceAllString(opaques, "")
			wa.opaque = util.TrimPrefixAndSuffix(opaques, " ")
		case staleRegexp.MatchString(raws):
			stales := staleRegexp.FindString(raws)
			stales = regexp.MustCompile(`(?i)(stale).*?=`).ReplaceAllString(stales, "")
			stales = regexp.MustCompile(`"`).ReplaceAllString(stales, "")
			wa.stale = util.TrimPrefixAndSuffix(stales, " ")
		case algorithmRegexp.MatchString(raws):
			algorithms := algorithmRegexp.FindString(raws)
			algorithms = regexp.MustCompile(`(?i)(algorithm).*?=`).ReplaceAllString(algorithms, "")
			algorithms = regexp.MustCompile(`"`).ReplaceAllString(algorithms, "")
			wa.algorithm = util.TrimPrefixAndSuffix(algorithms, " ")
		case qopOptionsRegexp.MatchString(raws):
			qopOptions := qopOptionsRegexp.FindString(raws)
			qopOptions = regexp.MustCompile(`(?i)(qop).*?=`).ReplaceAllString(qopOptions, "")
			qopOptions = regexp.MustCompile(`"`).ReplaceAllString(qopOptions, "")
			wa.qopOptions = util.TrimPrefixAndSuffix(qopOptions, " ")
		default:
			// auth-param
			if strings.Contains(raws, "=") {
				rs := strings.Split(raws, "=")
				if len(rs) > 1 {
					authParams[rs[0]] = rs[1]
				}
			}
		}
		wa.authParam = authParams
	}
	if len(strings.TrimSpace(wa.algorithm)) == 0 {
		wa.algorithm = "MD5"
	}
	return nil
}
func (wa *WWWAuthenticate) Validator() error {
	if wa == nil {
		return errors.New("www-authenticate caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(wa.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(www-authenticate)`).Match([]byte(wa.field)) {
		return errors.New("field is not match")
	}
	if wa.domain != nil {
		if err := wa.domain.Validator(); err != nil {
			return err
		}
	}

	return nil
}
