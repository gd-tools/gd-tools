package protocol

type Hello struct {
	Greeting string `json:"greeting,omitempty"`
}

func (req *Request) HasHello() bool {
	if req == nil {
		return false
	}
	return req.Greeting != ""
}
