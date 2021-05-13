package header

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"github.com/kokutas/sip"
	"github.com/kokutas/sip/util"
)

// The Authorization header field contains authentication credentials of
//  a UA. Section 22.2 overviews the use of the Authorization header
//  field, and Section 22.4 describes the syntax and semantics when used
//  with HTTP authentication.
//  This header field, along with Proxy-Authorization, breaks the general
//  rules about multiple header field values. Although not a comma-
//  separated list, this header field name may be present multiple times,
//  and MUST NOT be combined into a single header line using the usual
//  rules described in Section 7.3.
//  In the example below, there are no quotes around the Digest
//  parameter:
//  Authorization: Digest username="Alice", realm="atlanta.com",
//  				nonce="84a4cc6f3082121f32b42a2187831a9e",
//  				response="7587245234b3434cc3412213e5f113a5432"

// Authorization = "Authorization" HCOLON credentials
// credentials = ("Digest" LWS digest-response)
//  			/ other-response
// digest-response = dig-resp *(COMMA dig-resp)
// dig-resp = username / realm / nonce / digest-uri
//  		/ dresponse / algorithm / cnonce
//  		/ opaque / message-qop
//  		/ nonce-count / auth-param
// username = "username" EQUAL username-value
// username-value = quoted-string
// digest-uri = "uri" EQUAL LDQUOT digest-uri-value RDQUOT
// digest-uri-value = rquest-uri ; Equal to request-uri as specified
//  				by HTTP/1.1
// message-qop = "qop" EQUAL qop-value
// cnonce = "cnonce" EQUAL cnonce-value
// cnonce-value = nonce-value
// nonce-count = "nc" EQUAL nc-value
// nc-value = 8LHEX
// dresponse = "response" EQUAL request-digest
// request-digest = LDQUOT 32LHEX RDQUOT
// auth-param = auth-param-name EQUAL
//  			( token / quoted-string )
// auth-param-name = token
// other-response = auth-scheme LWS auth-param
//  				*(COMMA auth-param)
// auth-scheme = token
type Authorization struct {
	field      string
	authSchema string // Basic / Digest
	userName   string
	realm      string
	nonce      string
	uri        *sip.SipUri
	dresponse  string
	algorithm  string
	cnonce     string //cnonce-value      =  nonce-value ,nonce-count       =  "nc" EQUAL nc-value ,nc-value          =  8LHEX
	opaque     string
	qop        string
	nonceCount string
	authParam  map[string]interface{}
}

func (au *Authorization) Field() string {
	return au.field
}

func (au *Authorization) SetField(field string) {
	au.field = field
}

func (au *Authorization) AuthSchema() string {
	return au.authSchema
}

func (au *Authorization) SetAuthSchema(authSchema string) {
	au.authSchema = authSchema
}

func (au *Authorization) UserName() string {
	return au.userName
}

func (au *Authorization) SetUserName(userName string) {
	au.userName = userName
}

func (au *Authorization) Realm() string {
	return au.realm
}

func (au *Authorization) SetRealm(realm string) {
	au.realm = realm
}

func (au *Authorization) Nonce() string {
	return au.nonce
}

func (au *Authorization) SetNonce(nonce string) {
	au.nonce = nonce
}

func (au *Authorization) Uri() *sip.SipUri {
	return au.uri
}

func (au *Authorization) SetUri(uri *sip.SipUri) {
	au.uri = uri
}

func (au *Authorization) Dresponse() string {
	return au.dresponse
}

func (au *Authorization) SetDresponse(dresponse string) {
	au.dresponse = dresponse
}

func (au *Authorization) Algorithm() string {
	return au.algorithm
}

func (au *Authorization) SetAlgorithm(algorithm string) {
	au.algorithm = algorithm
}

func (au *Authorization) Cnonce() string {
	return au.cnonce
}

func (au *Authorization) SetCnonce(cnonce string) {
	au.cnonce = cnonce
}

func (au *Authorization) Opaque() string {
	return au.opaque
}

func (au *Authorization) SetOpaque(opaque string) {
	au.opaque = opaque
}

func (au *Authorization) Qop() string {
	return au.qop
}

func (au *Authorization) SetQop(qop string) {
	au.qop = qop
}

func (au *Authorization) NonceCount() string {
	return au.nonceCount
}

func (au *Authorization) SetNonceCount(nonceCount string) {
	au.nonceCount = nonceCount
}

func (au *Authorization) AuthParam() map[string]interface{} {
	return au.authParam
}

func (au *Authorization) SetAuthParam(authParam map[string]interface{}) {
	au.authParam = authParam
}

func NewAuthorization(authSchema string, userName string, realm string, nonce string, uri *sip.SipUri, dresponse string, algorithm string, cnonce string, opaque string, qop string, nonceCount string, authParam map[string]interface{}) *Authorization {
	return &Authorization{
		field:      "Authorization",
		authSchema: authSchema, userName: userName, realm: realm, nonce: nonce, uri: uri, dresponse: dresponse, algorithm: algorithm, cnonce: cnonce, opaque: opaque, qop: qop, nonceCount: nonceCount, authParam: authParam}
}

func (au *Authorization) Raw() (string, error) {
	result := ""
	if err := au.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(au.field)) == 0 {
		au.field = "Authorization"
	}
	result += fmt.Sprintf("%v:", strings.Title(au.field))

	if len(strings.TrimSpace(au.authSchema)) == 0 {
		au.authSchema = "Digest"
	}
	result += fmt.Sprintf(" %v ", strings.Title(au.authSchema))

	if len(strings.TrimSpace(au.userName)) > 0 {
		result += fmt.Sprintf("username=\"%v\",", au.userName)
	}
	if len(strings.TrimSpace(au.realm)) > 0 {
		result += fmt.Sprintf("realm=\"%v\",", au.realm)
	}
	if len(strings.TrimSpace(au.nonce)) > 0 {
		result += fmt.Sprintf("nonce=\"%v\",", au.nonce)
	}
	if au.uri != nil {
		res, err := au.uri.Raw()
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf("uri=\"%v\",", res)
	}
	if len(strings.TrimSpace(au.dresponse)) > 0 {
		result += fmt.Sprintf("response=\"%v\",", au.dresponse)
	}
	if len(strings.TrimSpace(au.cnonce)) > 0 {
		result += fmt.Sprintf("cnonce=\"%v\"", au.cnonce)
	}
	if len(strings.TrimSpace(au.opaque)) > 0 {
		result += fmt.Sprintf("opaque=\"%v\",", au.opaque)
	}
	if len(strings.TrimSpace(au.nonceCount)) > 0 {
		//nonce-count       =  "nc" EQUAL nc-value
		//nc-value          =  8LHEX
		result += fmt.Sprintf("nc=\"%8x\",", au.nonceCount)
	}
	if len(strings.TrimSpace(au.qop)) > 0 {
		result += fmt.Sprintf("qop=\"%v\",", au.qop)
	}
	if au.authParam != nil {
		for k, v := range au.authParam {
			result += fmt.Sprintf("%v=\"%v\",", k, v)
		}
	}
	if len(strings.TrimSpace(au.algorithm)) > 0 {
		result += fmt.Sprintf("algorithm=%v,", strings.ToUpper(au.algorithm))
	} else {
		result += fmt.Sprintf("algorithm=%v,", strings.ToUpper("MD5"))
	}
	result = strings.TrimSuffix(result, ",")
	result += "\r\n"
	return result, nil
}
func (au *Authorization) String() string {
	result := ""
	if len(strings.TrimSpace(au.field)) > 0 {
		result += fmt.Sprintf("field: %s,", au.field)
	}
	if len(strings.TrimSpace(au.authSchema)) == 0 {
		result += fmt.Sprintf("auth-schema: %s,", au.authSchema)
	}
	if len(strings.TrimSpace(au.userName)) > 0 {
		result += fmt.Sprintf("username: %s,", au.userName)
	}
	if len(strings.TrimSpace(au.realm)) > 0 {
		result += fmt.Sprintf("realm: %s,", au.realm)
	}
	if len(strings.TrimSpace(au.nonce)) > 0 {
		result += fmt.Sprintf("nonce: %v,", au.nonce)
	}
	if au.uri != nil {
		result += fmt.Sprintf("%s,", au.uri.String())
	}
	if len(strings.TrimSpace(au.dresponse)) > 0 {
		result += fmt.Sprintf("response: %v,", au.dresponse)
	}
	if len(strings.TrimSpace(au.cnonce)) > 0 {
		result += fmt.Sprintf("cnonce: %v,", au.cnonce)
	}
	if len(strings.TrimSpace(au.opaque)) > 0 {
		result += fmt.Sprintf("opaque: %v,", au.opaque)
	}
	if len(strings.TrimSpace(au.nonceCount)) > 0 {
		//nonce-count       =  "nc" EQUAL nc-value
		//nc-value          =  8LHEX
		result += fmt.Sprintf("nc: %8x,", au.nonceCount)
	}
	if len(strings.TrimSpace(au.qop)) > 0 {
		result += fmt.Sprintf("qop: %v,", au.qop)
	}
	if au.authParam != nil {
		result += fmt.Sprintf("auth-param: %v,", au.authParam)
	}
	if len(strings.TrimSpace(au.algorithm)) > 0 {
		result += fmt.Sprintf("algorithm: %v,", strings.ToUpper(au.algorithm))
	}
	result = strings.TrimSuffix(result, ",")
	return result
}
func (au *Authorization) Parser(raw string) error {
	if au == nil {
		return errors.New("authorization caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	//raw = strings.TrimLeft(raw, " ")
	//raw = strings.TrimRight(raw, " ")
	//raw = strings.TrimPrefix(raw, " ")
	//raw = strings.TrimSuffix(raw, " ")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}

	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(authorization).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		au.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// auth schema regexp
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if authSchemaRegexp.MatchString(raw) {
		au.authSchema = authSchemaRegexp.FindString(raw)
		raw = authSchemaRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// username regexp
	usernameRegexp := regexp.MustCompile(`(?i)(username).*?=.*`)
	// realm regexp
	realmRegexp := regexp.MustCompile(`(?i)(realm).*?=.*`)
	// nonce regexp
	nonceRegexp := regexp.MustCompile(`(?i)(nonce).*?=.*`)
	// uri regexp
	uriRegexp := regexp.MustCompile(`(?i)(uri).*?=.*`)
	// response regexp
	responseRegexp := regexp.MustCompile(`(?i)(response).*?=.*`)
	// cnonce regexp
	cnonceRegexp := regexp.MustCompile(`(?i)(cnonce).*?=.*`)
	// opaque regexp
	opaqueRegexp := regexp.MustCompile(`(?i)(opaque).*?=.*`)
	// qop-options regexp
	qopOptionsRegexp := regexp.MustCompile(`(?i)(qop).*?=.*`)
	// nonce-count regexp
	nonceCountRegexp := regexp.MustCompile(`(?i)(nc).*?=.*`)
	// algorithm regexp
	algorithmRegexp := regexp.MustCompile(`(?i)(algorithm).*?=.*`)
	// auth-param regexp
	authParams := make(map[string]interface{})

	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		switch {
		case usernameRegexp.MatchString(raws):
			usernames := usernameRegexp.FindString(raws)
			usernames = regexp.MustCompile(`(?i)(username).*?=`).ReplaceAllString(usernames, "")
			usernames = regexp.MustCompile(`"`).ReplaceAllString(usernames, "")
			usernames = strings.TrimLeft(usernames, " ")
			usernames = strings.TrimRight(usernames, " ")
			usernames = strings.TrimPrefix(usernames, " ")
			au.userName = strings.TrimSuffix(usernames, " ")
		case realmRegexp.MatchString(raws):
			realms := realmRegexp.FindString(raws)
			realms = regexp.MustCompile(`(?i)(realm).*?=`).ReplaceAllString(realms, "")
			realms = regexp.MustCompile(`"`).ReplaceAllString(realms, "")
			realms = strings.TrimLeft(realms, " ")
			realms = strings.TrimRight(realms, " ")
			realms = strings.TrimPrefix(realms, " ")
			au.realm = strings.TrimSuffix(realms, " ")
		case nonceRegexp.MatchString(raws):
			nonces := nonceRegexp.FindString(raws)
			nonces = regexp.MustCompile(`(?i)(nonce).*?=`).ReplaceAllLiteralString(nonces, "")
			nonces = regexp.MustCompile(`"`).ReplaceAllString(nonces, "")
			nonces = strings.TrimLeft(raw, " ")
			nonces = strings.TrimRight(nonces, " ")
			nonces = strings.TrimPrefix(nonces, " ")
			au.nonce = strings.TrimSuffix(nonces, " ")
		case uriRegexp.MatchString(raws):
			uris := uriRegexp.FindString(raws)
			uris = regexp.MustCompile(`(?i)(uri).*?=`).ReplaceAllString(uris, "")
			uris = regexp.MustCompile(`"`).ReplaceAllString(uris, "")
			uris = regexp.MustCompile(`<`).ReplaceAllString(uris, "")
			uris = regexp.MustCompile(`>`).ReplaceAllString(uris, "")
			uris = strings.TrimLeft(uris, " ")
			uris = strings.TrimRight(uris, " ")
			uris = strings.TrimPrefix(uris, " ")
			uris = strings.TrimSuffix(uris, " ")
			au.uri = new(sip.SipUri)
			if err := au.uri.Parser(uris); err != nil {
				return err
			}
		case responseRegexp.MatchString(raws):
			responses := responseRegexp.FindString(raws)
			responses = regexp.MustCompile(`(?i)(response).*?=`).ReplaceAllString(responses, "")
			responses = regexp.MustCompile(`"`).ReplaceAllString(responses, "")
			responses = strings.TrimLeft(responses, " ")
			responses = strings.TrimRight(responses, " ")
			responses = strings.TrimPrefix(responses, " ")
			au.dresponse = strings.TrimSuffix(responses, " ")
		case cnonceRegexp.MatchString(raws):
			cnonces := cnonceRegexp.FindString(raws)
			cnonces = regexp.MustCompile(`(?i)(cnonce).*?=`).ReplaceAllString(cnonces, "")
			cnonces = regexp.MustCompile(`"`).ReplaceAllString(cnonces, "")
			cnonces = strings.TrimLeft(cnonces, " ")
			cnonces = strings.TrimRight(cnonces, " ")
			cnonces = strings.TrimPrefix(cnonces, " ")
			au.cnonce = strings.TrimSuffix(cnonces, " ")
		case opaqueRegexp.MatchString(raws):
			opaques := opaqueRegexp.FindString(raws)
			raw = regexp.MustCompile(opaques).ReplaceAllString(raw, "")
			raw = strings.TrimLeft(raw, " ")
			raw = strings.TrimRight(raw, " ")
			raw = strings.TrimPrefix(raw, " ")
			raw = strings.TrimSuffix(raw, " ")
			opaques = regexp.MustCompile(`(?i)(opaque).*?=`).ReplaceAllString(opaques, "")
			opaques = regexp.MustCompile(`"`).ReplaceAllString(opaques, "")
			opaques = strings.TrimLeft(opaques, " ")
			raw = strings.TrimRight(raw, " ")
			opaques = strings.TrimPrefix(opaques, " ")
			au.opaque = strings.TrimSuffix(opaques, " ")
		case qopOptionsRegexp.MatchString(raws):
			qopOptions := qopOptionsRegexp.FindString(raws)
			qopOptions = regexp.MustCompile(`(?i)(qop).*?=`).ReplaceAllString(qopOptions, "")
			qopOptions = regexp.MustCompile(`"`).ReplaceAllString(qopOptions, "")
			qopOptions = strings.TrimLeft(qopOptions, " ")
			qopOptions = strings.TrimRight(qopOptions, " ")
			qopOptions = strings.TrimPrefix(qopOptions, " ")
			au.qop = strings.TrimSuffix(qopOptions, " ")
		case nonceCountRegexp.MatchString(raws):
			nonceCounts := nonceCountRegexp.FindString(raws)
			nonceCounts = regexp.MustCompile(`(?i)(nc).*?=`).ReplaceAllString(nonceCounts, "")
			nonceCounts = regexp.MustCompile(`"`).ReplaceAllString(nonceCounts, "")
			nonceCounts = strings.TrimLeft(nonceCounts, " ")
			nonceCounts = strings.TrimRight(nonceCounts, " ")
			nonceCounts = strings.TrimPrefix(nonceCounts, " ")
			au.nonceCount = strings.TrimSuffix(nonceCounts, " ")
		case algorithmRegexp.MatchString(raws):
			algorithms := algorithmRegexp.FindString(raws)
			algorithms = regexp.MustCompile(`(?i)(algorithm).*?=`).ReplaceAllString(algorithms, "")
			algorithms = regexp.MustCompile(`"`).ReplaceAllString(algorithms, "")
			au.algorithm = util.TrimPrefixAndSuffix(algorithms, " ")
		default:
			// auth-param
			if strings.Contains(raws, "=") {
				rs := strings.Split(raws, "=")
				if len(rs) > 1 {
					authParams[rs[0]] = rs[1]
				}
			}
		}
		au.authParam = authParams
	}
	if len(strings.TrimSpace(au.algorithm)) == 0 {
		au.algorithm = "MD5"
	}
	return nil
}
func (au *Authorization) Validator() error {
	if au == nil {
		return errors.New("authorization caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(au.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(authorization)`).Match([]byte(au.field)) {
		return errors.New("field is not match")
	}
	if err := au.uri.Validator(); err != nil {
		return err
	}
	return nil
}
