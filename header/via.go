package header

import (
	"errors"
	"fmt"
	"github.com/kokutas/sip"
	"github.com/kokutas/sip/util"
	"regexp"
	"strconv"
	"strings"
)

// The Via header field indicates the path taken by the request so far
//  and indicates the path that should be followed in routing responses.
//  The branch ID parameter in the Via header field values serves as a
//  transaction identifier, and is used by proxies to detect loops.
//  A Via header field value contains the transport protocol used to send
//  the message, the client’s host name or network address, and possibly
//  the port number at which it wishes to receive responses. A Via
//  header field value can also contain parameters such as "maddr",
//  "ttl", "received", and "branch", whose meaning and use are described
//  in other sections. For implementations compliant to this
//  specification, the value of the branch parameter MUST start with the
//  magic cookie "z9hG4bK", as discussed in Section 8.1.1.7.
//  Transport protocols defined here are "UDP", "TCP", "TLS", and "SCTP".
//  "TLS" means TLS over TCP. When a request is sent to a SIPS URI, the
//  protocol still indicates "SIP", and the transport protocol is TLS.
// Via: SIP/2.0/UDP erlang.bell-telephone.com:5060;branch=z9hG4bK87asdks7
// Via: SIP/2.0/UDP 192.0.2.1:5060 ;received=192.0.2.207
//  ;branch=z9hG4bK77asjd
//  The compact form of the Via header field is v.
//  In this example, the message originated from a multi-homed host with
//  two addresses, 192.0.2.1 and 192.0.2.207. The sender guessed wrong
//  as to which network interface would be used. Erlang.bell-
//  telephone.com noticed the mismatch and added a parameter to the
//  previous hop’s Via header field value, containing the address that
//  the packet actually came from.
//  The host or network address and port number are not required to
//  follow the SIP URI syntax. Specifically, LWS on either side of the
//  ":" or "/" is allowed, as shown here:
//  Via: SIP / 2.0 / UDP first.example.com: 4000;ttl=16
//  ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1
//  Even though this specification mandates that the branch parameter be
//  present in all requests, the BNF for the header field indicates that
//  it is optional. This allows interoperation with RFC 2543 elements,
//  which did not have to insert the branch parameter.
//  Two Via header fields are equal if their sent-protocol and sent-by
//  fields are equal, both have the same set of parameters, and the
//  values of all parameters are equal.

// Via = ( "Via" / "v" ) HCOLON via-parm *(COMMA via-parm)
// via-parm = sent-protocol LWS sent-by *( SEMI via-params )
// via-params = via-ttl / via-maddr
//  			/ via-received / via-branch
//  			/ via-extension
// via-ttl = "ttl" EQUAL ttl
// via-maddr = "maddr" EQUAL host
// via-received = "received" EQUAL (IPv4address / IPv6address)
// via-branch = "branch" EQUAL token
// via-extension = generic-param
// sent-protocol = protocol-name SLASH protocol-version
//  				SLASH transport
// protocol-name = "SIP" / token
// protocol-version = token
// transport = "UDP" / "TCP" / "TLS" / "SCTP"
//  			/ other-transport
// sent-by = host [ COLON port ]
// ttl = 1*3DIGIT ; 0 to 255

type Via struct {
	field           string
	*sip.SipVersion
	transport       string
	*sip.HostPort
	rport           uint16
	*sip.Parameters
	branch          string
	received        string
}

func (via *Via) Field() string {
	return via.field
}

func (via *Via) SetField(field string) {
	via.field = field
}

func (via *Via) Rport() uint16 {
	return via.rport
}

func (via *Via) SetRport(rport uint16) {
	via.rport = rport
}

func (via *Via) Branch() string {
	return via.branch
}

func (via *Via) SetBranch(branch string) {
	via.branch = branch
}

func (via *Via) Received() string {
	return via.received
}

func (via *Via) SetReceived(received string) {
	via.received = received
}


func NewVia(sipVersion *sip.SipVersion, transport string, hostPort *sip.HostPort, rport uint16, parameters *sip.Parameters, branch string, received string) *Via {
	return &Via{
		field:      "Via",
		SipVersion: sipVersion,
		transport:  transport,
		HostPort:   hostPort,
		rport:      rport,
		Parameters: parameters,
		branch:     branch,
		received:   received,
	}
}
func (via *Via) Raw() (string,error) {
	result := ""
	if err:=via.Validator();err!=nil {
		return result,err
	}
	if len(strings.TrimSpace(via.field))==0{
		via.field="Via"
	}
	result += fmt.Sprintf("%s:", strings.Title(via.field))
	if via.SipVersion != nil {
		res,err:=via.SipVersion.Raw()
		if err!=nil{
			return "",err
		}
		result += fmt.Sprintf(" %s", res)
	}
	if len(strings.TrimSpace(via.transport)) > 0 {
		result += fmt.Sprintf("/%s", strings.ToUpper(via.transport))
	}
	if via.HostPort != nil {
		res,err:=via.HostPort.Raw()
		if err!=nil{
			return "",err
		}
		result += fmt.Sprintf(" %s", res)
	}
	if via.rport == 1 {
		result += fmt.Sprintf(";%v", "rport")
	}
	if via.rport > 1 {
		result += fmt.Sprintf(";rport=%v", via.rport)
	}
	if via.Parameters != nil {
		res,err:=via.Parameters.Raw()
		if err!=nil{
			return "", err
		}
		result += fmt.Sprintf("%s", res)
	}
	if len(strings.TrimSpace(via.branch)) > 0 {
		result += fmt.Sprintf(";branch=%v", via.branch)
	}
	if len(strings.TrimSpace(via.received)) > 0 {
		result += fmt.Sprintf(";received=%v", via.received)
	}
	result += "\r\n"
	return result,nil
}
func (via *Via) String() string {
	result := ""
	if len(strings.TrimSpace(via.field))>0{
		result +=fmt.Sprintf("field: %s,",via.field)
	}
	if via.SipVersion!=nil{
		result+=fmt.Sprintf("%s,",via.SipVersion.String())
	}
	if len(strings.TrimSpace(via.transport))>0{
		result +=fmt.Sprintf("transport: %s,",via.transport)
	}
	if via.HostPort!=nil{
		result+=fmt.Sprintf("%s,",via.HostPort.String())
	}
	result +=fmt.Sprintf("rport: %d,",via.rport)
	if via.Parameters!=nil{
		result+=fmt.Sprintf("%s,",via.Parameters.String())
	}
	if len(strings.TrimSpace(via.branch))>0{
		result +=fmt.Sprintf("via-branch: %s,",via.branch)
	}
	if len(strings.TrimSpace(via.received))>0{
		result +=fmt.Sprintf("via-received: %s,",via.received)
	}
	result = strings.TrimSuffix(result,",")
	return result
}
func (via *Via) Parser(raw string) error {
	if via == nil {
		return errors.New("via caller is not allowed via be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("raw parameter is not allowed to be empty")
	}
	// filed regexp
	fieldRegexp := regexp.MustCompile(`(?i)(via).*?:`)
	if fieldRegexp.MatchString(raw) {
		field := fieldRegexp.FindString(raw)
		via.field = regexp.MustCompile(`:`).ReplaceAllString(field, "")
		raw = strings.ReplaceAll(raw, field, "")
		raw = strings.TrimSuffix(raw, " ")
		raw = strings.TrimPrefix(raw, " ")
	}
	// schemas regexp
	schemasRegexpStr := `(?i)(`
	for _, v := range sip.Schemas {
		schemasRegexpStr += v + "|"
	}
	schemasRegexpStr = strings.TrimSuffix(schemasRegexpStr, "|")
	schemasRegexpStr += ")"
	// sip-version regexp
	sipVersionRegexp := regexp.MustCompile(schemasRegexpStr + `/\d+\.\d*`)
	if sipVersionRegexp.MatchString(raw) {
		sipVersion := sipVersionRegexp.FindString(raw)
		via.SipVersion = new(sip.SipVersion)
		if err := via.SipVersion.Parser(sipVersion); err != nil {
			return err
		}
		raw = sipVersionRegexp.ReplaceAllString(raw, "")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	// transport regexp
	transportsRegexpStr := `(?i)(`
	for _, v := range sip.Transports {
		transportsRegexpStr += v + "|"
	}
	transportsRegexpStr = strings.TrimSuffix(transportsRegexpStr, "|")
	transportsRegexpStr += ")"
	transportRegexp := regexp.MustCompile(`/` + transportsRegexpStr)
	if transportRegexp.MatchString(raw) {
		via.transport = regexp.MustCompile(`/`).ReplaceAllString(transportRegexp.FindString(raw), "")
		raw = transportRegexp.ReplaceAllString(raw, "")
		raw = util.TrimPrefixAndSuffix(raw, " ")
	}
	raw = util.TrimPrefixAndSuffix(raw, " ")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// hostport regexp
	// hostportRegexp := regexp.MustCompile(`[.*\:\d+]|[.*]`)
	// rport regexp
	rportRegexp := regexp.MustCompile(`(?i)(rport).*`)
	// branch regexp
	branchRegexp := regexp.MustCompile(`(?i)(branch).*`)
	// received regexp
	receivedRegexp := regexp.MustCompile(`(?i)(received).*`)
	rawSlice := strings.Split(raw, ";")
	for k, raws := range rawSlice {
		switch {
		// case hostportRegexp.MatchString(raws):
		case k == 0:
			if regexp.MustCompile(`(?i)(received)`).MatchString(raws) {
				continue
			}
			if regexp.MustCompile(`(?i)(maddr)`).MatchString(raws) {
				continue
			}
			// raw = hostportRegexp.ReplaceAllString(raw, "")
			// raw = util.TrimPrefixAndSuffix(raw, " ")
			// hostport := hostportRegexp.FindString(raws)
			hostport := raws
			hostport = util.TrimPrefixAndSuffix(hostport, " ")
			raw = util.TrimPrefixAndSuffix(raw, " ")
			raw = strings.TrimLeft(raw, hostport)
			raw = util.TrimPrefixAndSuffix(raw, " ")
			via.HostPort = new(sip.HostPort)
			if err := via.HostPort.Parser(hostport); err != nil {
				return err
			}

		case rportRegexp.MatchString(raws):
			raw = regexp.MustCompile(rportRegexp.FindString(raws)).ReplaceAllString(raw, "")
			raw = util.TrimPrefixAndSuffix(raw, ";")
			raw = util.TrimPrefixAndSuffix(raw, " ")
			rports := regexp.MustCompile(`(?i)(rport)`).ReplaceAllString(rportRegexp.FindString(raws), "")
			rports = regexp.MustCompile(`=`).ReplaceAllString(rports, "")
			rports = util.TrimPrefixAndSuffix(rports, " ")
			if len(strings.TrimSpace(rports)) > 0 {
				rport, err := strconv.Atoi(strings.TrimSpace(rports))
				if err != nil {
					return err
				}
				via.rport = uint16(rport)
			} else {
				via.rport = 1
			}
		case branchRegexp.MatchString(raws):
			raw = regexp.MustCompile(branchRegexp.FindString(raws)).ReplaceAllString(raw, "")
			raw = util.TrimPrefixAndSuffix(raw, ";")
			raw = util.TrimPrefixAndSuffix(raw, " ")
			branchs := regexp.MustCompile(`(?i)(branch)`).ReplaceAllString(branchRegexp.FindString(raws), "")
			branchs = regexp.MustCompile(`=`).ReplaceAllString(branchs, "")
			branchs = util.TrimPrefixAndSuffix(branchs, " ")
			via.branch = branchs
		case receivedRegexp.MatchString(raws):
			raw = regexp.MustCompile(receivedRegexp.FindString(raws)).ReplaceAllString(raw, "")
			raw = util.TrimPrefixAndSuffix(raw, ";")
			raw = util.TrimPrefixAndSuffix(raw, " ")
			received := regexp.MustCompile(`(?i)(received)`).ReplaceAllString(receivedRegexp.FindString(raws), "")
			received = regexp.MustCompile(`=`).ReplaceAllString(received, "")
			received = util.TrimPrefixAndSuffix(received, " ")
			via.received = received
		}
	}
	raw = util.TrimPrefixAndSuffix(raw, ";")
	raw = util.TrimPrefixAndSuffix(raw, " ")
	// parameters regexp
	if len(strings.TrimSpace(raw)) > 0 {
		via.Parameters = new(sip.Parameters)
		if err := via.Parameters.Parser(raw); err != nil {
			return err
		}
	}

	return nil
}
func (via *Via) Validator() error {
	if via == nil {
		return errors.New("via caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(via.field)) == 0 {
		return errors.New("field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(via)`).Match([]byte(via.field)) {
		return errors.New("field is not match")
	}
	if err := via.SipVersion.Validator(); err != nil {
		return err
	}
	if len(strings.TrimSpace(via.transport)) == 0 {
		return errors.New("transport is not allowed to be empty")
	}
	// transport regexp
	transportsRegexpStr := `(?i)(`
	for _, v := range sip.Transports {
		transportsRegexpStr += v + "|"
	}
	transportsRegexpStr = strings.TrimSuffix(transportsRegexpStr, "|")
	transportsRegexpStr += ")"
	transportRegexp := regexp.MustCompile(transportsRegexpStr)
	if !transportRegexp.MatchString(via.transport) {
		return errors.New("transport is not match")
	}
	if err := via.HostPort.Validator(); err != nil {
		return err
	}
	if len(strings.TrimSpace(via.branch)) == 0 {
		return errors.New("branch is not allowed to be empty")
	}
	return nil
}
