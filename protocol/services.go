package protocol

// ServiceList contains system services that should be restarted or reloaded.
type ServiceList struct {
	Services []string `json:"services,omitempty"`
}

// AddService adds a service to the request if it is non-empty and not already present.
func (req *Request) AddService(service string) {
	if req == nil || service == "" {
		return
	}
	for _, check := range req.Services {
		if check == service {
			return
		}
	}
	req.Services = append(req.Services, service)
}

// HasServiceList reports whether the request contains at least one service entry.
func (req *Request) HasServiceList() bool {
	if req == nil {
		return false
	}
	return len(req.Services) > 0
}
