package transport

import "crypto/tls"

// Session contains transport runtime state and is not part of the wire protocol.
type Session struct {
	Conn    *tls.Conn
	Verbose bool
}
