package protocol

type ServiceList struct {
	Services []string `json:"services,omitempty"`
}

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

func (req *Request) HasServiceList() bool {
	if req == nil {
		return false
	}
	return len(req.Services) > 0
}
