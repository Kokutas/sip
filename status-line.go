package sip

// https://www.rfc-editor.org/rfc/rfc3261.html#section-7.2
//
// 7.2 Responses
// SIP responses are distinguished from requests by having a Status-Line
// as their start-line.  A Status-Line consists of the protocol version
// followed by a numeric Status-Code and its associated textual phrase,
// with each element separated by a single SP character.

// No CR or LF is allowed except in the final CRLF sequence.

// 	Status-Line  =  SIP-Version SP Status-Code SP Reason-Phrase CRLF

// The Status-Code is a 3-digit integer result code that indicates the
// outcome of an attempt to understand and satisfy a request.  The
// Reason-Phrase is intended to give a short textual description of the
// Status-Code.  The Status-Code is intended for use by automata,
// whereas the Reason-Phrase is intended for the human user.  A client
// is not required to examine or display the Reason-Phrase.

// While this specification suggests specific wording for the reason
// phrase, implementations MAY choose other text, for example, in the
// language indicated in the Accept-Language header field of the
// request.
// The first digit of the Status-Code defines the class of response.
// The last two digits do not have any categorization role.  For this
// reason, any response with a status code between 100 and 199 is
// referred to as a "1xx response", any response with a status code
// between 200 and 299 as a "2xx response", and so on.  SIP/2.0 allows
// six values for the first digit:

// 	1xx: Provisional -- request received, continuing to process the
// 		request;

// 	2xx: Success -- the action was successfully received, understood,
// 		and accepted;

// 	3xx: Redirection -- further action needs to be taken in order to
// 		complete the request;

// 	4xx: Client Error -- the request contains bad syntax or cannot be
// 		fulfilled at this server;

// 	5xx: Server Error -- the server failed to fulfill an apparently
// 		valid request;

// 	6xx: Global Failure -- the request cannot be fulfilled at any
// 		server.

// Section 21 defines these classes and describes the individual codes.

// https://www.rfc-editor.org/rfc/rfc3261.html#section-25.1
//
// Status-Line     =  SIP-Version SP Status-Code SP Reason-Phrase CRLF
// Status-Code     =  Informational
//                /   Redirection
//                /   Success
//                /   Client-Error
//                /   Server-Error
//                /   Global-Failure
//                /   extension-code
// extension-code  =  3DIGIT
// Reason-Phrase   =  *(reserved / unreserved / escaped
//                    / UTF8-NONASCII / UTF8-CONT / SP / HTAB)
