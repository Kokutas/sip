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
	schema string
}

func (au *AbsoluteUri) Schema() string {
	return au.schema
}

func (au *AbsoluteUri) SetSchema(schema string) {
	au.schema = schema
}
func NewAbsoluteUri(schema string) *AbsoluteUri {
	return &AbsoluteUri{schema: schema}
}

func (au *AbsoluteUri) Raw() (string,error) {
	result := ""
	return result,nil
}
func (au *AbsoluteUri) String() string {
	result := ""
	return result
}
func (au *AbsoluteUri) Parser(raw string) error {
	return nil
}
func (au *AbsoluteUri) Validator() error {
	return nil
}
