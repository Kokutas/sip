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
	Field      string                 `json:"field"`
	AuthSchema string                 `json:"auth-scheme"` // Basic / Digest
	UserName   string                 `json:"username"`
	Realm      string                 `json:"realm"`
	Nonce      string                 `json:"nonce"`
	Uri        *sip.SipUri            `json:"digest-uri"`
	Dresponse  string                 `json:"dresponse"`
	Algorithm  string                 `json:"algorithm"`
	Cnonce     string                 `json:"cnonce"` //cnonce-value      =  nonce-value ,nonce-count       =  "nc" EQUAL nc-value ,nc-value          =  8LHEX
	Opaque     string                 `json:"opaque"`
	Qop        string                 `json:"message-qop"`
	NonceCount string                 `json:"nonce-count"`
	AuthParam  map[string]interface{} `json:"auth-param"`
}

func CreateAuthorization() sip.Sip {
	return &Authorization{}
}

func NewAuthorization(authSchema string, userName string, realm string, nonce string, uri *sip.SipUri, dresponse string, algorithm string, cnonce string, opaque string, qop string, nonceCount string, authParam map[string]interface{}) sip.Sip {
	return &Authorization{
		Field:      "Authorization",
		AuthSchema: authSchema,
		UserName:   userName,
		Realm:      realm,
		Nonce:      nonce,
		Uri:        uri,
		Dresponse:  dresponse,
		Algorithm:  algorithm,
		Cnonce:     cnonce,
		Opaque:     opaque,
		Qop:        qop,
		NonceCount: nonceCount,
		AuthParam:  authParam,
	}
}
func (au *Authorization) Raw() string {
	result := ""
	if reflect.DeepEqual(nil, au) {
		return result
	}
	result += fmt.Sprintf("%v:", strings.Title(au.Field))
	if len(strings.TrimSpace(au.AuthSchema)) == 0 {
		au.AuthSchema = "Digest"
	}
	result += fmt.Sprintf(" %v ", strings.Title(au.AuthSchema))

	if len(strings.TrimSpace(au.UserName)) > 0 {
		result += fmt.Sprintf("username=\"%v\",", au.UserName)
	}
	if len(strings.TrimSpace(au.Realm)) > 0 {
		result += fmt.Sprintf("realm=\"%v\",", au.Realm)
	}
	if len(strings.TrimSpace(au.Nonce)) > 0 {
		result += fmt.Sprintf("nonce=\"%v\",", au.Nonce)
	}
	if au.Uri != nil {
		result += fmt.Sprintf("uri=\"%v\",", au.Uri.Raw())
	}
	if len(strings.TrimSpace(au.Dresponse)) > 0 {
		result += fmt.Sprintf("response=\"%v\",", au.Dresponse)
	}
	if len(strings.TrimSpace(au.Cnonce)) > 0 {
		result += fmt.Sprintf("cnonce=\"%v\"", au.Cnonce)
	}
	if len(strings.TrimSpace(au.Opaque)) > 0 {
		result += fmt.Sprintf("opaque=\"%v\",", au.Opaque)
	}
	if len(strings.TrimSpace(au.NonceCount)) > 0 {
		//nonce-count       =  "nc" EQUAL nc-value
		//nc-value          =  8LHEX
		result += fmt.Sprintf("nc=\"%8x\",", au.NonceCount)
	}
	if len(strings.TrimSpace(au.Qop)) > 0 {
		result += fmt.Sprintf("qop=\"%v\",", au.Qop)
	}
	if au.AuthParam != nil {
		for k, v := range au.AuthParam {
			result += fmt.Sprintf("%v=\"%v\",", k, v)
		}
	}
	if len(strings.TrimSpace(au.Algorithm)) > 0 {
		result += fmt.Sprintf("algorithm=%v", strings.ToUpper(au.Algorithm))
	} else {
		result += fmt.Sprintf("algorithm=%v", strings.ToUpper("MD5"))
	}
	result += "\r\n"
	return result
}
func (au *Authorization) JsonString() string {
	result := ""
	if reflect.DeepEqual(nil, au) {
		return result
	}
	data, err := json.Marshal(au)
	if err != nil {
		return result
	}
	result = fmt.Sprintf("%s", data)
	return result
}
func (au *Authorization) Parser(raw string) error {
	if au == nil {
		return errors.New("authorization caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")

	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(authorization).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		au.Field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// auth schema regexp
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if authSchemaRegexp.MatchString(raw) {
		au.AuthSchema = authSchemaRegexp.FindString(raw)
		raw = authSchemaRegexp.ReplaceAllString(raw, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
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

	raw = util.TrimPrefixAndSuffix(raw, " ")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		switch {
		case usernameRegexp.MatchString(raws):
			usernames := usernameRegexp.FindString(raws)
			usernames = regexp.MustCompile(`(?i)(username).*?=`).ReplaceAllString(usernames, "")
			usernames = regexp.MustCompile(`"`).ReplaceAllString(usernames, "")
			au.UserName = util.TrimPrefixAndSuffix(usernames, " ")
		case realmRegexp.MatchString(raws):
			realms := realmRegexp.FindString(raws)
			realms = regexp.MustCompile(`(?i)(realm).*?=`).ReplaceAllString(realms, "")
			realms = regexp.MustCompile(`"`).ReplaceAllString(realms, "")
			au.Realm = util.TrimPrefixAndSuffix(realms, " ")
		case nonceRegexp.MatchString(raws):
			nonces := nonceRegexp.FindString(raws)
			nonces = regexp.MustCompile(`(?i)(nonce).*?=`).ReplaceAllLiteralString(nonces, "")
			nonces = regexp.MustCompile(`"`).ReplaceAllString(nonces, "")
			au.Nonce = util.TrimPrefixAndSuffix(nonces, " ")
		case uriRegexp.MatchString(raws):
			uris := uriRegexp.FindString(raws)
			uris = regexp.MustCompile(`(?i)(uri).*?=`).ReplaceAllString(uris, "")
			uris = regexp.MustCompile(`"`).ReplaceAllString(uris, "")
			uris = regexp.MustCompile(`<`).ReplaceAllString(uris, "")
			uris = regexp.MustCompile(`>`).ReplaceAllString(uris, "")
			uris = util.TrimPrefixAndSuffix(uris, " ")
			au.Uri = sip.CreateSipUri().(*sip.SipUri)
			if err := au.Uri.Parser(uris); err != nil {
				return err
			}
		case responseRegexp.MatchString(raws):
			responses := responseRegexp.FindString(raws)
			responses = regexp.MustCompile(`(?i)(response).*?=`).ReplaceAllString(responses, "")
			responses = regexp.MustCompile(`"`).ReplaceAllString(responses, "")
			au.Dresponse = util.TrimPrefixAndSuffix(responses, " ")
		case cnonceRegexp.MatchString(raws):
			cnonces := cnonceRegexp.FindString(raws)
			cnonces = regexp.MustCompile(`(?i)(cnonce).*?=`).ReplaceAllString(cnonces, "")
			cnonces = regexp.MustCompile(`"`).ReplaceAllString(cnonces, "")
			au.Cnonce = util.TrimPrefixAndSuffix(cnonces, " ")
		case opaqueRegexp.MatchString(raws):
			opaques := opaqueRegexp.FindString(raws)
			raw = regexp.MustCompile(opaques).ReplaceAllString(raw, "")
			raw = util.TrimPrefixAndSuffix(raw, " ")
			opaques = regexp.MustCompile(`(?i)(opaque).*?=`).ReplaceAllString(opaques, "")
			opaques = regexp.MustCompile(`"`).ReplaceAllString(opaques, "")
			au.Opaque = util.TrimPrefixAndSuffix(opaques, " ")
		case qopOptionsRegexp.MatchString(raws):
			qopOptions := qopOptionsRegexp.FindString(raws)
			qopOptions = regexp.MustCompile(`(?i)(qop).*?=`).ReplaceAllString(qopOptions, "")
			qopOptions = regexp.MustCompile(`"`).ReplaceAllString(qopOptions, "")
			au.Qop = util.TrimPrefixAndSuffix(qopOptions, " ")
		case nonceCountRegexp.MatchString(raws):
			nonceCounts := nonceCountRegexp.FindString(raws)
			nonceCounts = regexp.MustCompile(`(?i)(nc).*?=`).ReplaceAllString(nonceCounts, "")
			nonceCounts = regexp.MustCompile(`"`).ReplaceAllString(nonceCounts, "")
			au.NonceCount = util.TrimPrefixAndSuffix(nonceCounts, " ")
		case algorithmRegexp.MatchString(raws):
			algorithms := algorithmRegexp.FindString(raws)
			algorithms = regexp.MustCompile(`(?i)(algorithm).*?=`).ReplaceAllString(algorithms, "")
			algorithms = regexp.MustCompile(`"`).ReplaceAllString(algorithms, "")
			au.Algorithm = util.TrimPrefixAndSuffix(algorithms, " ")
		default:
			// auth-param
			if strings.Contains(raws, "=") {
				rs := strings.Split(raws, "=")
				if len(rs) > 1 {
					authParams[rs[0]] = rs[1]
				}
			}
		}
		au.AuthParam = authParams
	}
	if len(strings.TrimSpace(au.Algorithm)) == 0 {
		au.Algorithm = "MD5"
	}
	return nil
}
func (au *Authorization) Validator() error {
	if au == nil {
		return errors.New("authorization caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(au.Field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(authorization)`).Match([]byte(au.Field)) {
		return errors.New("field is not match")
	}
	if err := au.Uri.Validator(); err != nil {
		return err
	}
	return nil
}
