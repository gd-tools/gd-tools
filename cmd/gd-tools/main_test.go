package main

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/gd-tools/gd-tools/agent"
)

func startTestAgent(t *testing.T) (addr string, stop func()) {

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			go handleConnection(conn)
		}
	}()

	return ln.Addr().String(), func() {
		ln.Close()
	}
}

func sendRequest(t *testing.T, addr string, req agent.Request) agent.Response {

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(2 * time.Second))

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		t.Fatal(err)
	}

	var resp agent.Response

	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp
}

func TestAgentProtocolMismatch(t *testing.T) {

	addr, stop := startTestAgent(t)
	defer stop()

	req := agent.Request{
		Version: 999,
	}

	resp := sendRequest(t, addr, req)

	if resp.Err == "" {
		t.Fatal("expected protocol mismatch error")
	}
}

func TestAgentHandlerPanic(t *testing.T) {

	old := Handlers

	Handlers = []Handler{
		{
			Name: "panic",
			Test: func(*agent.Request) bool { return true },
			Func: func(*agent.Request, *agent.Response) error {
				panic("boom")
			},
		},
	}

	defer func() {
		Handlers = old
	}()

	addr, stop := startTestAgent(t)
	defer stop()

	req := agent.Request{
		Version: agent.ProtocolVersion,
	}

	resp := sendRequest(t, addr, req)

	if resp.Err == "" {
		t.Fatal("expected panic error")
	}
}
