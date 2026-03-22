package protocol

import (
	"encoding/json"
	"fmt"
)

// Response is the shared response message between prod and dev.
//
// Execution stops at the first error. Therefore a single Error field is used
// instead of a list of errors.
type Response struct {
	Version int    `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`

	// Informational output
	InfoLines []string `json:"info_lines,omitempty"`

	// Services that have been (re)started
	ServiceList

	// User IDs that Dev needs to know (e.g. vmail)
	UserIDs []UserID `json:"user_ids,omitempty"`

	// Dev needs the RustDesk credentials for crash recovery
	RustDeskApp
}

func (resp *Response) String() string {
	if resp == nil {
		return ""
	}
	content, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "failed to marshal Response"
	}
	return string(content)
}

func (resp *Response) HasError() bool {
	return resp != nil && resp.Error != ""
}

func (resp *Response) SetError(err error) {
	if resp == nil || err == nil {
		return
	}
	resp.Error = err.Error()
}

func (resp *Response) Info(args ...any) {
	if resp == nil {
		return
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			resp.InfoLines = append(resp.InfoLines, v...)
		default:
			resp.InfoLines = append(resp.InfoLines, fmt.Sprint(v))
		}
	}
}

func (resp *Response) Infof(format string, args ...any) {
	if resp == nil || format == "" {
		return
	}
	resp.InfoLines = append(resp.InfoLines, fmt.Sprintf(format, args...))
}

func (resp *Response) AddService(service string) {
	if resp == nil || service == "" {
		return
	}
	for _, check := range resp.Services {
		if check == service {
			return
		}
	}
	resp.Services = append(resp.Services, service)
}
