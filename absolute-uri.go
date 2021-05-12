package sip


// absoluteURI    =  scheme ":" ( hier-part / opaque-part )
// hier-part      =  ( net-path / abs-path ) [ "?" query ]
// net-path       =  "//" authority [ abs-path ]
// abs-path       =  "/" path-segments
// opaque-part    =  uric-no-slash *uric
// uric           =  reserved / unreserved / escaped
// uric-no-slash  =  unreserved / escaped / ";" / "?" / ":" / "@"
//                   / "&" / "=" / "+" / "$" / ","
// path-segments  =  segment *( "/" segment )
// segment        =  *pchar *( ";" param )
// param          =  *pchar
// pchar          =  unreserved / escaped /
//                   ":" / "@" / "&" / "=" / "+" / "$" / ","
// scheme         =  ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
// authority      =  srvr / reg-name
// srvr           =  [ [ userinfo "@" ] hostport ]
// reg-name       =  1*( unreserved / escaped / "$" / ","
//                   / ";" / ":" / "@" / "&" / "=" / "+" )
// query          =  *uric
type AbsoluteUri struct {
	Schema string `json:"schema"`
}
func CreateAbsoluteUri()Sip{
	return &AbsoluteUri{}
}
func (au *AbsoluteUri) Raw() string {
	result := ""
	return result
}
func (au *AbsoluteUri) JsonString() string {
	result := ""
	return result
}
func (au *AbsoluteUri) Parser(raw string) error {
	return nil
}
func (au *AbsoluteUri) Validator() error {
	return nil
}