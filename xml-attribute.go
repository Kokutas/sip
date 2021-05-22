package sip

type Authorization struct {
	authSchema string // auth-schema: Basic / Digest
	username   string // username
	realm      string // realm
	nonce      string // nonce
	uri        *Uri   // Uri
	response   string // response
	algorithm  string // algorithm
}
