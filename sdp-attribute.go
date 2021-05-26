package sip

// https://www.rfc-editor.org/rfc/rfc2327.html
//
//https://www.rfc-editor.org/rfc/rfc2327.html#section-6
//

// 6.  SDP Specification

// SDP session descriptions are entirely textual using the ISO 10646
// character set in UTF-8 encoding. SDP field names and attributes names
// use only the US-ASCII subset of UTF-8, but textual fields and
// attribute values may use the full ISO 10646 character set.  The
// textual form, as opposed to a binary encoding such as ASN/1 or XDR,
// was chosen to enhance portability, to enable a variety of transports
// to be used (e.g, session description in a MIME email message) and to
// allow flexible, text-based toolkits (e.g., Tcl/Tk ) to be used to
// generate and to process session descriptions.  However, since the
// total bandwidth allocated to all SAP announcements is strictly
// limited, the encoding is deliberately compact.  Also, since
// announcements may be transported via very unreliable means (e.g.,
// email) or damaged by an intermediate caching server, the encoding was
// designed with strict order and formatting rules so that most errors
// would result in malformed announcements which could be detected
// easily and discarded. This also allows rapid discarding of encrypted
// announcements for which a receiver does not have the correct key.

// An SDP session description consists of a number of lines of text of
// the form <type>=<value> <type> is always exactly one character and is
// case-significant.  <value> is a structured text string whose format
// depends on <type>.  It also will be case-significant unless a
// specific field defines otherwise.  Whitespace is not permitted either
// side of the `=' sign. In general <value> is either a number of fields
// delimited by a single space character or a free format string.

// A session description consists of a session-level description
// (details that apply to the whole session and all media streams) and
// optionally several media-level descriptions (details that apply onto
// to a single media stream).

// An announcement consists of a session-level section followed by zero
// or more media-level sections.  The session-level part starts with a
// `v=' line and continues to the first media-level section.  The media
// description starts with an `m=' line and continues to the next media
// description or end of the whole session description.  In general,
// session-level values are the default for all media unless overridden
// by an equivalent media-level value.

// When SDP is conveyed by SAP, only one session description is allowed
// per packet.  When SDP is conveyed by other means, many SDP session
// descriptions may be concatenated together (the `v=' line indicating
// the start of a session description terminates the previous
// description).  Some lines in each description are required and some
// are optional but all must appear in exactly the order given here (the
// fixed order greatly enhances error detection and allows for a simple
// parser). Optional items are marked with a `*'.

// Session description
// 	v=  (protocol version)
// 	o=  (owner/creator and session identifier).
// 	s=  (session name)
// 	i=* (session information)
// 	u=* (URI of description)
// 	e=* (email address)
// 	p=* (phone number)
// 	c=* (connection information - not required if included in all media)
// 	b=* (bandwidth information)
// 	One or more time descriptions (see below)
// 	z=* (time zone adjustments)
// 	k=* (encryption key)
// 	a=* (zero or more session attribute lines)
// 	Zero or more media descriptions (see below)

// Time description
// 	t=  (time the session is active)
// 	r=* (zero or more repeat times)

// Media description
// 	m=  (media name and transport address)
// 	i=* (media title)
// 	c=* (connection information - optional if included at session-level)
// 	b=* (bandwidth information)
// 	k=* (encryption key)
// 	a=* (zero or more media attribute lines)

// The set of `type' letters is deliberately small and not intended to
// be extensible -- SDP parsers must completely ignore any announcement
// that contains a `type' letter that it does not understand. The
// `attribute' mechanism ("a=" described below) is the primary means for
// extending SDP and tailoring it to particular applications or media.
// Some attributes (the ones listed in this document) have a defined
// meaning but others may be added on an application-, media- or
// session-specific basis.  A session directory must ignore any
// attribute it doesn't understand.

// The connection (`c=') and attribute (`a=') information in the
// session-level section applies to all the media of that session unless
// overridden by connection information or an attribute of the same name
// in the media description.  For instance, in the example below, each
// media behaves as if it were given a `recvonly' attribute.

// An example SDP description is:

// 	v=0
// 	o=mhandley 2890844526 2890842807 IN IP4 126.16.64.4
// 	s=SDP Seminar
// 	i=A Seminar on the session description protocol
// 	u=http://www.cs.ucl.ac.uk/staff/M.Handley/sdp.03.ps
// 	e=mjh@isi.edu (Mark Handley)
// 	c=IN IP4 224.2.17.12/127
// 	t=2873397496 2873404696
// 	a=recvonly
// 	m=audio 49170 RTP/AVP 0
// 	m=video 51372 RTP/AVP 31
// 	m=application 32416 udp wb
// 	a=orient:portrait

// Text records such as the session name and information are bytes
// strings which may contain any byte with the exceptions of 0x00 (Nul),
// 0x0a (ASCII newline) and 0x0d (ASCII carriage return).  The sequence
// CRLF (0x0d0a) is used to end a record, although parsers should be
// tolerant and also accept records terminated with a single newline
// character.  By default these byte strings contain ISO-10646
// characters in UTF-8 encoding, but this default may be changed using
// the `charset' attribute.

// Protocol Version

// v=0

// The "v=" field gives the version of the Session Description Protocol.
// There is no minor version number.

// Origin

// o=<username> <session id> <version> <network type> <address type>
// <address>

// The "o=" field gives the originator of the session (their username
// and the address of the user's host) plus a session id and session
// version number.

// <username> is the user's login on the originating host, or it is "-"
// if the originating host does not support the concept of user ids.
// <username> must not contain spaces.  <session id> is a numeric string
// such that the tuple of <username>, <session id>, <network type>,
// <address type> and <address> form a globally unique identifier for
// the session.

// The method of <session id> allocation is up to the creating tool, but
// it has been suggested that a Network Time Protocol (NTP) timestamp be
// used to ensure uniqueness [1].

// <version> is a version number for this announcement.  It is needed
// for proxy announcements to detect which of several announcements for
// the same session is the most recent.  Again its usage is up to the
// creating tool, so long as <version> is increased when a modification
// is made to the session data.  Again, it is recommended (but not
// mandatory) that an NTP timestamp is used.

// <network type> is a text string giving the type of network.
// Initially "IN" is defined to have the meaning "Internet".  <address
// type> is a text string giving the type of the address that follows.
// Initially "IP4" and "IP6" are defined.  <address> is the globally
// unique address of the machine from which the session was created.
// For an address type of IP4, this is either the fully-qualified domain
// name of the machine, or the dotted-decimal representation of the IP
// version 4 address of the machine.  For an address type of IP6, this
// is either the fully-qualified domain name of the machine, or the
// compressed textual representation of the IP version 6 address of the
// machine.  For both IP4 and IP6, the fully-qualified domain name is
// the form that SHOULD be given unless this is unavailable, in which
// case the globally unique address may be substituted.  A local IP
// address MUST NOT be used in any context where the SDP description
// might leave the scope in which the address is meaningful.

// In general, the "o=" field serves as a globally unique identifier for
// this version of this session description, and the subfields excepting
// the version taken together identify the session irrespective of any
// modifications.

// Session Name

// s=<session name>

// The "s=" field is the session name.  There must be one and only one
// "s=" field per session description, and it must contain ISO 10646
// characters (but see also the `charset' attribute below).

// Session and Media Information

// i=<session description>

// The "i=" field is information about the session.  There may be at
// most one session-level "i=" field per session description, and at
// most one "i=" field per media. Although it may be omitted, this is
// discouraged for session announcements, and user interfaces for
// composing sessions should require text to be entered.  If it is
// present it must contain ISO 10646 characters (but see also the
// `charset' attribute below).

// A single "i=" field can also be used for each media definition.  In
// media definitions, "i=" fields are primarily intended for labeling
// media streams. As such, they are most likely to be useful when a
// single session has more than one distinct media stream of the same
// media type.  An example would be two different whiteboards, one for
// slides and one for feedback and questions.

// URI

// u=<URI>

// o A URI is a Universal Resource Identifier as used by WWW clients

// o The URI should be a pointer to additional information about the
// 	conference

// o This field is optional, but if it is present it should be specified
// 	before the first media field

// o No more than one URI field is allowed per session description

// Email Address and Phone Number

// e=<email address>
// p=<phone number>

// o These specify contact information for the person responsible for
// 	the conference.  This is not necessarily the same person that
// 	created the conference announcement.

// o Either an email field or a phone field must be specified.
// 	Additional email and phone fields are allowed.

// o If these are present, they should be specified before the first
// 	media field.

// o More than one email or phone field can be given for a session
// 	description.

// o Phone numbers should be given in the conventional international

// 	format - preceded by a "+ and the international country code.
// 	There must be a space or a hyphen ("-") between the country code
// 	and the rest of the phone number.  Spaces and hyphens may be used
// 	to split up a phone field to aid readability if desired. For
// 	example:

// p=+44-171-380-7777 or p=+1 617 253 6011
// o Both email addresses and phone numbers can have an optional free
// text string associated with them, normally giving the name of the
// person who may be contacted.  This should be enclosed in
// parenthesis if it is present.  For example:

// 					e=mjh@isi.edu (Mark Handley)

// The alternative RFC822 name quoting convention is also allowed for
// both email addresses and phone numbers.  For example,

// 					e=Mark Handley <mjh@isi.edu>

// The free text string should be in the ISO-10646 character set with
// UTF-8 encoding, or alternatively in ISO-8859-1 or other encodings
// if the appropriate charset session-level attribute is set.

// Connection Data

// c=<network type> <address type> <connection address>

// The "c=" field contains connection data.

// A session announcement must contain one "c=" field in each media
// description (see below) or a "c=" field at the session-level.  It may
// contain a session-level "c=" field and one additional "c=" field per
// media description, in which case the per-media values override the
// session-level settings for the relevant media.

// The first sub-field is the network type, which is a text string
// giving the type of network.  Initially "IN" is defined to have the
// meaning "Internet".

// The second sub-field is the address type.  This allows SDP to be used
// for sessions that are not IP based.  Currently only IP4 is defined.

// The third sub-field is the connection address.  Optional extra
// subfields may be added after the connection address depending on the
// value of the <address type> field.

// For IP4 addresses, the connection address is defined as follows:

// o Typically the connection address will be a class-D IP multicast

// group address.  If the session is not multicast, then the
// connection address contains the fully-qualified domain name or the
// unicast IP address of the expected data source or data relay or
// data sink as determined by additional attribute fields. It is not
// expected that fully-qualified domain names or unicast addresses
// will be given in a session description that is communicated by a
// multicast announcement, though this is not prohibited.  If a
// unicast data stream is to pass through a network address
// translator, the use of a fully-qualified domain name rather than an
// unicast IP address is RECOMMENDED.  In other cases, the use of an
// IP address to specify a particular interface on a multi-homed host
// might be required.  Thus this specification leaves the decision as
// to which to use up to the individual application, but all
// applications MUST be able to cope with receiving both formats.

// o Conferences using an IP multicast connection address must also have
// a time to live (TTL) value present in addition to the multicast
// address.  The TTL and the address together define the scope with
// which multicast packets sent in this conference will be sent. TTL
// values must be in the range 0-255.

// The TTL for the session is appended to the address using a slash as
// a separator.  An example is:

// 						c=IN IP4 224.2.1.1/127

// Hierarchical or layered encoding schemes are data streams where the
// encoding from a single media source is split into a number of
// layers.  The receiver can choose the desired quality (and hence
// bandwidth) by only subscribing to a subset of these layers.  Such
// layered encodings are normally transmitted in multiple multicast
// groups to allow multicast pruning.  This technique keeps unwanted
// traffic from sites only requiring certain levels of the hierarchy.
// For applications requiring multiple multicast groups, we allow the
// following notation to be used for the connection address:

// 	<base multicast address>/<ttl>/<number of addresses>

// If the number of addresses is not given it is assumed to be one.
// Multicast addresses so assigned are contiguously allocated above
// the base address, so that, for example:

// 				c=IN IP4 224.2.1.1/127/3

// would state that addresses 224.2.1.1, 224.2.1.2 and 224.2.1.3 are
// to be used at a ttl of 127.  This is semantically identical to
// including multiple "c=" lines in a media description:

// 	c=IN IP4 224.2.1.1/127
// 	c=IN IP4 224.2.1.2/127
// 	c=IN IP4 224.2.1.3/127
// 	Multiple addresses or "c=" lines can only be specified on a per-
// 	media basis, and not for a session-level "c=" field.

// 	It is illegal for the slash notation described above to be used for
// 	IP unicast addresses.

// Bandwidth

// b=<modifier>:<bandwidth-value>

// o This specifies the proposed bandwidth to be used by the session or
// 	media, and is optional.

// o <bandwidth-value> is in kilobits per second

// o <modifier> is a single alphanumeric word giving the meaning of the
// 	bandwidth figure.

// o Two modifiers are initially defined:

// CT Conference Total: An implicit maximum bandwidth is associated with
// 	each TTL on the Mbone or within a particular multicast
// 	administrative scope region (the Mbone bandwidth vs. TTL limits are
// 	given in the MBone FAQ). If the bandwidth of a session or media in
// 	a session is different from the bandwidth implicit from the scope,
// 	a `b=CT:...' line should be supplied for the session giving the
// 	proposed upper limit to the bandwidth used. The primary purpose of
// 	this is to give an approximate idea as to whether two or more
// 	conferences can co-exist simultaneously.

// AS Application-Specific Maximum: The bandwidth is interpreted to be
// 	application-specific, i.e., will be the application's concept of
// 	maximum bandwidth.  Normally this will coincide with what is set on
// 	the application's "maximum bandwidth" control if applicable.

// 	Note that CT gives a total bandwidth figure for all the media at
// 	all sites.  AS gives a bandwidth figure for a single media at a
// 	single site, although there may be many sites sending
// 	simultaneously.

// o Extension Mechanism: Tool writers can define experimental bandwidth
// 	modifiers by prefixing their modifier with "X-". For example:

// 								b=X-YZ:128

// 	SDP parsers should ignore bandwidth fields with unknown modifiers.
// 	Modifiers should be alpha-numeric and, although no length limit is
// 	given, they are recommended to be short.
// 	Times, Repeat Times and Time Zones

// 	t=<start time>  <stop time>

// 	o "t=" fields specify the start and stop times for a conference
// 	session.  Multiple "t=" fields may be used if a session is active
// 	at multiple irregularly spaced times; each additional "t=" field
// 	specifies an additional period of time for which the session will
// 	be active.  If the session is active at regular times, an "r="
// 	field (see below) should be used in addition to and following a
// 	"t=" field - in which case the "t=" field specifies the start and
// 	stop times of the repeat sequence.

// 	o The first and second sub-fields give the start and stop times for
// 	the conference respectively.  These values are the decimal
// 	representation of Network Time Protocol (NTP) time values in
// 	seconds [1].  To convert these values to UNIX time, subtract
// 	decimal 2208988800.

// 	o If the stop-time is set to zero, then the session is not bounded,
// 	though it will not become active until after the start-time.  If
// 	the start-time is also zero, the session is regarded as permanent.

// 	User interfaces should strongly discourage the creation of
// 	unbounded and permanent sessions as they give no information about
// 	when the session is actually going to terminate, and so make
// 	scheduling difficult.

// 	The general assumption may be made, when displaying unbounded
// 	sessions that have not timed out to the user, that an unbounded
// 	session will only be active until half an hour from the current
// 	time or the session start time, whichever is the later.  If
// 	behaviour other than this is required, an end-time should be given
// 	and modified as appropriate when new information becomes available
// 	about when the session should really end.

// 	Permanent sessions may be shown to the user as never being active
// 	unless there are associated repeat times which state precisely when
// 	the session will be active.  In general, permanent sessions should
// 	not be created for any session expected to have a duration of less
// 	than 2 months, and should be discouraged for sessions expected to
// 	have a duration of less than 6 months.

// 	r=<repeat interval> <active duration> <list of offsets from start-
// 	time>

// 	o "r=" fields specify repeat times for a session.  For example, if
// 	a session is active at 10am on Monday and 11am on Tuesday for one
// 	hour each week for three months, then the <start time> in the
// 	corresponding "t=" field would be the NTP representation of 10am on
// 	the first Monday, the <repeat interval> would be 1 week, the
// 	<active duration> would be 1 hour, and the offsets would be zero
// 	and 25 hours. The corresponding "t=" field stop time would be the
// 	NTP representation of the end of the last session three months
// 	later. By default all fields are in seconds, so the "r=" and "t="
// 	fields might be:

// 							t=3034423619 3042462419
// 							r=604800 3600 0 90000

// 	To make announcements more compact, times may also be given in units
// 	of days, hours or minutes. The syntax for these is a number
// 	immediately followed by a single case-sensitive character.
// 	Fractional units are not allowed - a smaller unit should be used
// 	instead.  The following unit specification characters are allowed:

// 						d - days (86400 seconds)
// 						h - minutes (3600 seconds)
// 						m - minutes (60 seconds)
// 		s - seconds (allowed for completeness but not recommended)

// 	Thus, the above announcement could also have been written:

// 								r=7d 1h 0 25h

// 	Monthly and yearly repeats cannot currently be directly specified
// 	with a single SDP repeat time - instead separate "t" fields should
// 	be used to explicitly list the session times.

// 		z=<adjustment time> <offset> <adjustment time> <offset> ....

// 	o To schedule a repeated session which spans a change from daylight-
// 	saving time to standard time or vice-versa, it is necessary to
// 	specify offsets from the base repeat times. This is required
// 	because different time zones change time at different times of day,
// 	different countries change to or from daylight time on different
// 	dates, and some countries do not have daylight saving time at all.

// 	Thus in order to schedule a session that is at the same time winter
// 	and summer, it must be possible to specify unambiguously by whose
// 	time zone a session is scheduled.  To simplify this task for
// 	receivers, we allow the sender to specify the NTP time that a time
// 	zone adjustment happens and the offset from the time when the
// 	session was first scheduled.  The "z" field allows the sender to
// 	specify a list of these adjustment times and offsets from the base
// 	time.
// 	An example might be:

// 	z=2882844526 -1h 2898848070 0

// This specifies that at time 2882844526 the time base by which the
// session's repeat times are calculated is shifted back by 1 hour,
// and that at time 2898848070 the session's original time base is
// restored. Adjustments are always relative to the specified start
// time - they are not cumulative.

// o    If a session is likely to last several years, it is  expected
// that
// the session announcement will be modified periodically rather than
// transmit several years worth of adjustments in one announcement.

// Encryption Keys

// k=<method>
// k=<method>:<encryption key>

// o The session description protocol may be used to convey encryption
// keys.  A key field is permitted before the first media entry (in
// which case it applies to all media in the session), or for each
// media entry as required.

// o The format of keys and their usage is outside the scope of this
// document, but see [3].

// o The method indicates the mechanism to be used to obtain a usable
// key by external means, or from the encoded encryption key given.

// The following methods are defined:

// 	k=clear:<encryption key>
// 	The encryption key (as described in [3] for  RTP  media  streams
// 	under  the  AV  profile)  is  included untransformed in this key
// 	field.

// 	k=base64:<encoded encryption key>
// 	The encryption key (as described in [3] for RTP media streams
// 	under the AV profile) is included in this key field but has been
// 	base64 encoded because it includes characters that are
// 	prohibited in SDP.

// 	k=uri:<URI to obtain key>
// 	A Universal Resource Identifier as used by WWW clients is
// 	included in this key field.  The URI refers to the data
// 	containing the key, and may require additional authentication

// 	before the key can be returned.  When a request is made to the
// 	given URI, the MIME content-type of the reply specifies the
// 	encoding for the key in the reply.  The key should not be
// 	obtained until the user wishes to join the session to reduce
// 	synchronisation of requests to the WWW server(s).

// 	k=prompt
// 	No key is included in this SDP description, but the session or
// 	media stream referred to by this key field is encrypted.  The
// 	user should be prompted for the key when attempting to join the
// 	session, and this user-supplied key should then be used to
// 	decrypt the media streams.

// Attributes

// a=<attribute>
// a=<attribute>:<value>

// Attributes are the primary means for extending SDP.  Attributes may
// be defined to be used as "session-level" attributes, "media-level"
// attributes, or both.

// A media description may have any number of attributes ("a=" fields)
// which are media specific.  These are referred to as "media-level"
// attributes and add information about the media stream.  Attribute
// fields can also be added before the first media field; these
// "session-level" attributes convey additional information that applies
// to the conference as a whole rather than to individual media; an
// example might be the conference's floor control policy.

// Attribute fields may be of two forms:

// o property attributes.  A property attribute is simply of the form
// "a=<flag>".  These are binary attributes, and the presence of the
// attribute conveys that the attribute is a property of the session.
// An example might be "a=recvonly".

// o value attributes.  A value attribute is of the form
// "a=<attribute>:<value>".  An example might be that a whiteboard
// could have the value attribute "a=orient:landscape"

// Attribute interpretation depends on the media tool being invoked.
// Thus receivers of session descriptions should be configurable in
// their interpretation of announcements in general and of attributes in
// particular.

// Attribute names must be in the US-ASCII subset of ISO-10646/UTF-8.
// Attribute values are byte strings, and MAY use any byte value except
// 0x00 (Nul), 0x0A (LF), and 0x0D (CR). By default, attribute values
// are to be interpreted as in ISO-10646 character set with UTF-8
// encoding.  Unlike other text fields, attribute values are NOT
// normally affected by the `charset' attribute as this would make
// comparisons against known values problematic.  However, when an
// attribute is defined, it can be defined to be charset-dependent, in
// which case it's value should be interpreted in the session charset
// rather than in ISO-10646.

// Attributes that will be commonly used can be registered with IANA
// (see Appendix B).  Unregistered attributes should begin with "X-" to
// prevent inadvertent collision with registered attributes.  In either
// case, if an attribute is received that is not understood, it should
// simply be ignored by the receiver.

// Media Announcements

// m=<media> <port> <transport> <fmt list>

// A session description may contain a number of media descriptions.
// Each media description starts with an "m=" field, and is terminated
// by either the next "m=" field or by the end of the session
// description.  A media field also has several sub-fields:

// o The first sub-field is the media type.  Currently defined media are
// "audio", "video", "application", "data" and "control", though this
// list may be extended as new communication modalities emerge (e.g.,
// telepresense).  The difference between "application" and "data" is
// that the former is a media flow such as whiteboard information, and
// the latter is bulk-data transfer such as multicasting of program
// executables which will not typically be displayed to the user.
// "control" is used to specify an additional conference control
// channel for the session.

// o The second sub-field is the transport port to which the media
// stream will be sent.  The meaning of the transport port depends on
// the network being used as specified in the relevant "c" field and
// on the transport protocol defined in the third sub-field.  Other
// ports used by the media application (such as the RTCP port, see
// [2]) should be derived algorithmically from the base media port.

// Note: For transports based on UDP, the value should be in the range
// 1024 to 65535 inclusive.  For RTP compliance it should be an even
// number.
// For applications where hierarchically encoded streams are being
// sent to a unicast address, it may be necessary to specify multiple
// transport ports.  This is done using a similar notation to that
// used for IP multicast addresses in the "c=" field:

// 		m=<media> <port>/<number of ports> <transport> <fmt list>

// In such a case, the ports used depend on the transport protocol.
// For RTP, only the even ports are used for data and the
// corresponding one-higher odd port is used for RTCP.  For example:

// 					m=video 49170/2 RTP/AVP 31

// would specify that ports 49170 and 49171 form one RTP/RTCP pair and
// 49172 and 49173 form the second RTP/RTCP pair.  RTP/AVP is the
// transport protocol and 31 is the format (see below).

// It is illegal for both multiple addresses to be specified in the
// "c=" field and for multiple ports to be specified in the "m=" field
// in the same session description.

// o The third sub-field is the transport protocol.  The transport
// protocol values are dependent on the address-type field in the "c="
// fields.  Thus a "c=" field of IP4 defines that the transport
// protocol runs over IP4.  For IP4, it is normally expected that most
// media traffic will be carried as RTP over UDP.  The following
// transport protocols are preliminarily defined, but may be extended
// through registration of new protocols with IANA:

// - RTP/AVP - the IETF's Realtime Transport Protocol using the
// 	Audio/Video profile carried over UDP.

// - udp - User Datagram Protocol

// If an application uses a single combined proprietary media format
// and transport protocol over UDP, then simply specifying the
// transport protocol as udp and using the format field to distinguish
// the combined protocol is recommended.  If a transport protocol is
// used over UDP to carry several distinct media types that need to be
// distinguished by a session directory, then specifying the transport
// protocol and media format separately is necessary. RTP is an
// example of a transport-protocol that carries multiple payload
// formats that must be distinguished by the session directory for it
// to know how to start appropriate tools, relays, mixers or
// recorders.
// ......
//
// https://www.rfc-editor.org/rfc/rfc4566.html
//
//https://www.rfc-editor.org/rfc/rfc4566.html#section-6
//
// 6.  SDP Attributes
//
// The following attributes are defined.  Since application writers may
// add new attributes as they are required, this list is not exhaustive.
// Registration procedures for new attributes are defined in Section
// 8.2.4.

// a=cat:<category>

// 	This attribute gives the dot-separated hierarchical category of
// 	the session.  This is to enable a receiver to filter unwanted
// 	sessions by category.  There is no central registry of
// 	categories.  It is a session-level attribute, and it is not
// 	dependent on charset.
// a=keywds:<keywords>

// Like the cat attribute, this is to assist identifying wanted
// sessions at the receiver.  This allows a receiver to select
// interesting session based on keywords describing the purpose of
// the session; there is no central registry of keywords.  It is a
// session-level attribute.  It is a charset-dependent attribute,
// meaning that its value should be interpreted in the charset
// specified for the session description if one is specified, or
// by default in ISO 10646/UTF-8.

// a=tool:<name and version of tool>

// This gives the name and version number of the tool used to
// create the session description.  It is a session-level
// attribute, and it is not dependent on charset.

// a=ptime:<packet time>

// This gives the length of time in milliseconds represented by
// the media in a packet.  This is probably only meaningful for
// audio data, but may be used with other media types if it makes
// sense.  It should not be necessary to know ptime to decode RTP
// or vat audio, and it is intended as a recommendation for the
// encoding/packetisation of audio.  It is a media-level
// attribute, and it is not dependent on charset.

// a=maxptime:<maximum packet time>

// This gives the maximum amount of media that can be encapsulated
// in each packet, expressed as time in milliseconds.  The time
// SHALL be calculated as the sum of the time the media present in
// the packet represents.  For frame-based codecs, the time SHOULD
// be an integer multiple of the frame size.  This attribute is
// probably only meaningful for audio data, but may be used with
// other media types if it makes sense.  It is a media-level
// attribute, and it is not dependent on charset.  Note that this
// attribute was introduced after RFC 2327, and non-updated
// implementations will ignore this attribute.

// a=rtpmap:<payload type> <encoding name>/<clock rate> [/<encoding
// parameters>]

// This attribute maps from an RTP payload type number (as used in
// an "m=" line) to an encoding name denoting the payload format
// to be used.  It also provides information on the clock rate and
// encoding parameters.  It is a media-level attribute that is not
// dependent on charset.
//...
//

// Session description
// 	v=  (protocol version)
// 	o=  (owner/creator and session identifier).
// 	s=  (session name)
// 	i=* (session information)
// 	u=* (URI of description)
// 	e=* (email address)
// 	p=* (phone number)
// 	c=* (connection information - not required if included in all media)
// 	b=* (bandwidth information)
// 	One or more time descriptions (see below)
// 	z=* (time zone adjustments)
// 	k=* (encryption key)
// 	a=* (zero or more session attribute lines)
// 	Zero or more media descriptions (see below)

// Time description
// 	t=  (time the session is active)
// 	r=* (zero or more repeat times)

// Media description
// 	m=  (media name and transport address)
// 	i=* (media title)
// 	c=* (connection information - optional if included at session-level)
// 	b=* (bandwidth information)
// 	k=* (encryption key)
// 	a=* (zero or more media attribute lines)
type SdpAttribute struct {
	field  string // a/m/v etc.
	value  string
	source string // source string
	// Cat []byte // Named portion of URI
	// Val []byte // Port number
	// Src []byte // Full source if needed
}

// Media description
// 	m=  (media name and transport address)
// 	i=* (media title)
// 	c=* (connection information - optional if included at session-level)
// 	b=* (bandwidth information)
// 	k=* (encryption key)
// 	a=* (zero or more media attribute lines)
type MediaDescription struct {
	m string // 	m=  (media name and transport address)
	i string // 	i=* (media title)
	c string // 	c=* (connection information - optional if included at session-level)
	b string // 	b=* (bandwidth information)
	k string // 	k=* (encryption key)
	a string // 	a=* (zero or more media attribute lines)
}
