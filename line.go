package sip

// mixing request lines and status lines
type Line struct {
	// Method     []byte // Sip Method eg INVITE etc
	// UriType    string // Type of URI sip, sips, tel etc
	// StatusCode []byte // Status Code eg 100
	// StatusDesc []byte // Status Code Description eg trying
	// User       []byte // User part
	// Host       []byte // Host part
	// Port       []byte // Port number
	// UserType   []byte // User Type
	// Src        []byte // Full source if needed
}
