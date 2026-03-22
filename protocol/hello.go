package protocol

// Hello is a minimal round-trip test message.
// Convention: "Hello" → "World"
type Hello struct {
	Greeting string `json:"greeting,omitempty"`
}

// HasPing reports whether the request contains a ping message.
func (req *Request) HasHello() bool {
	if req == nil {
		return false
	}
	return req.Greeting != ""
}
