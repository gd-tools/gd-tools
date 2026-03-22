package protocol

import (
	"errors"
	"strings"
	"testing"
)

func TestResponseStringNil(t *testing.T) {
	var resp *Response
	if got := resp.String(); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestResponseStringContainsFields(t *testing.T) {
	resp := &Response{
		Version: 1,
		Error:   "boom",
		UserIDs: []UserID{
			{Name: "vmail", UID: "2000", GID: "2000"},
		},
	}
	resp.Info("hello")
	resp.Infof("service: %s", "caddy")
	resp.AddService("caddy")

	got := resp.String()

	if !strings.Contains(got, `"version": 1`) {
		t.Fatalf("expected version in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"error": "boom"`) {
		t.Fatalf("expected error in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"info_lines": [`) {
		t.Fatalf("expected info_lines in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"hello"`) {
		t.Fatalf("expected hello in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"service: caddy"`) {
		t.Fatalf("expected formatted line in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"services": [`) {
		t.Fatalf("expected services in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"caddy"`) {
		t.Fatalf("expected caddy in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"user_ids": [`) {
		t.Fatalf("expected user_ids in output, got:\n%s", got)
	}
	if !strings.Contains(got, `"name": "vmail"`) {
		t.Fatalf("expected vmail user id in output, got:\n%s", got)
	}
}

func TestResponseHasError(t *testing.T) {
	var nilResp *Response
	if nilResp.HasError() {
		t.Fatal("nil response must not have error")
	}

	resp := &Response{}
	if resp.HasError() {
		t.Fatal("empty response must not have error")
	}

	resp.Error = "boom"
	if !resp.HasError() {
		t.Fatal("response with error must have error")
	}
}

func TestResponseSetError(t *testing.T) {
	resp := &Response{}
	resp.SetError(errors.New("failed"))

	if resp.Error != "failed" {
		t.Fatalf("expected error to be set, got %q", resp.Error)
	}
}

func TestResponseSetErrorNil(t *testing.T) {
	resp := &Response{}
	resp.SetError(nil)

	if resp.Error != "" {
		t.Fatalf("expected empty error, got %q", resp.Error)
	}
}

func TestResponseInfo(t *testing.T) {
	resp := &Response{}
	resp.Info("one", 2, []string{"three", "four"})

	want := []string{"one", "2", "three", "four"}
	if len(resp.InfoLines) != len(want) {
		t.Fatalf("expected %d info lines, got %d", len(want), len(resp.InfoLines))
	}
	for i := range want {
		if resp.InfoLines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, resp.InfoLines[i], want[i])
		}
	}
}

func TestResponseInfof(t *testing.T) {
	resp := &Response{}
	resp.Infof("hello %s", "world")

	if len(resp.InfoLines) != 1 {
		t.Fatalf("expected 1 info line, got %d", len(resp.InfoLines))
	}
	if resp.InfoLines[0] != "hello world" {
		t.Fatalf("got %q, want %q", resp.InfoLines[0], "hello world")
	}
}

func TestResponseInfofEmptyFormat(t *testing.T) {
	resp := &Response{}
	resp.Infof("")

	if len(resp.InfoLines) != 0 {
		t.Fatalf("expected no info lines, got %d", len(resp.InfoLines))
	}
}

func TestResponseAddService(t *testing.T) {
	resp := &Response{}
	resp.AddService("caddy")
	resp.AddService("postfix")
	resp.AddService("caddy")

	want := []string{"caddy", "postfix"}
	if len(resp.Services) != len(want) {
		t.Fatalf("expected %d services, got %d", len(want), len(resp.Services))
	}
	for i := range want {
		if resp.Services[i] != want[i] {
			t.Fatalf("service %d = %q, want %q", i, resp.Services[i], want[i])
		}
	}
}

func TestResponseAddServiceIgnoresEmpty(t *testing.T) {
	resp := &Response{}
	resp.AddService("")

	if len(resp.Services) != 0 {
		t.Fatalf("expected no services, got %d", len(resp.Services))
	}
}
