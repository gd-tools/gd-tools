package agent

import "testing"

func TestHelloTestNil(t *testing.T) {
	if HelloTest(nil) {
		t.Fatal("expected false for nil request")
	}
}

func TestHelloTestEmpty(t *testing.T) {
	req := &Request{}

	if HelloTest(req) {
		t.Fatal("expected false for empty hello")
	}
}

func TestHelloTestValid(t *testing.T) {
	req := &Request{
		Hello: "Hi",
	}

	if !HelloTest(req) {
		t.Fatal("expected true for hello request")
	}
}

func TestHelloHandlerNil(t *testing.T) {
	err := HelloHandler(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHelloHandlerValid(t *testing.T) {
	req := &Request{
		Hello: "Hi",
	}

	resp := &Response{}

	err := HelloHandler(req, resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Result) == 0 {
		t.Fatal("expected response message")
	}

	if resp.Result[0] != "Hello, World!" {
		t.Fatalf("unexpected response: %s", resp.Result[0])
	}
}
