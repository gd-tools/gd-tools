package agent

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"
)

func makeTLSConn(t *testing.T) (*tls.Conn, *tls.Conn) {
	clientRaw, serverRaw := net.Pipe()

	// generate self-signed cert
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("key gen failed: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("cert gen failed: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("keypair failed: %v", err)
	}

	serverTLS := tls.Server(serverRaw, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	clientTLS := tls.Client(clientRaw, &tls.Config{
		InsecureSkipVerify: true,
	})

	go serverTLS.Handshake()

	if err := clientTLS.Handshake(); err != nil {
		t.Fatalf("TLS handshake failed: %v", err)
	}

	return clientTLS, serverTLS
}

func TestSendCommunication(t *testing.T) {
	clientTLS, serverTLS := makeTLSConn(t)
	defer clientTLS.Close()
	defer serverTLS.Close()

	go func() {
		dec := json.NewDecoder(serverTLS)
		enc := json.NewEncoder(serverTLS)

		var req Request
		if err := dec.Decode(&req); err != nil {
			t.Error(err)
			return
		}

		resp := Response{}
		resp.Say("hello client")

		if err := enc.Encode(resp); err != nil {
			t.Error(err)
		}
	}()

	req := Request{
		Version: ProtocolVersion,
		Conn:    clientTLS,
	}

	if err := req.Send(); err != nil {
		t.Fatalf("send failed: %v", err)
	}
}

func TestSendResponseError(t *testing.T) {
	clientTLS, serverTLS := makeTLSConn(t)
	defer clientTLS.Close()
	defer serverTLS.Close()

	go func() {
		dec := json.NewDecoder(serverTLS)
		enc := json.NewEncoder(serverTLS)

		var req Request
		if err := dec.Decode(&req); err != nil {
			t.Error(err)
			return
		}

		resp := Response{
			Err: "test error",
		}

		enc.Encode(resp)
	}()

	req := Request{
		Version: ProtocolVersion,
		Conn:    clientTLS,
	}

	err := req.Send()
	if err == nil {
		t.Fatal("expected error")
	}
}
