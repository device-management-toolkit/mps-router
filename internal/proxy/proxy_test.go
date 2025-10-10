/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package proxy

import (
	"database/sql"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/device-management-toolkit/mps-router/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestServer_parseGuid(t *testing.T) {
	type args struct{ content string }
	server := Server{}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty String", args{content: "\n"}, ""},
		{"Invalid Guid", args{content: "GET /api/v1/amt/log/audit/12345?startIndex=0 HTTP/1.1\nHost: example"}, ""},
		{"No Guid", args{content: "GET /api/v1/devices HTTP/1.1\nHost: example"}, ""},
		{"Valid v4 GUID", args{content: "GET /api/v1/amt/log/audit/63f32fee-238e-4f6a-a091-092270d22439?startIndex=0 HTTP/1.1\nHost: example"}, "63f32fee-238e-4f6a-a091-092270d22439"},
		{"Valid v1 GUID", args{content: "GET /api/v1/amt/features/63f32fee-238e-1f6a-a091-092270d22439 HTTP/1.1\nHost: example"}, "63f32fee-238e-1f6a-a091-092270d22439"},
		{"Valid GUID Websocket Request", args{content: "GET /relay/webrelay.ashx?p=2&host=d12428be-9fa1-4226-9784-54b2038beab6&port=16994 HTTP/1.1\nHost: example"}, "d12428be-9fa1-4226-9784-54b2038beab6"},
		{"Invalid GUID Websocket Request", args{content: "GET /relay/webrelay.ashx?p=2&host=d12428be-9fa1-4226-9784&port=16994 HTTP/1.1\nHost: example"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.parseGuid(tt.args.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListenAndServe(t *testing.T) {
	server := Server{Addr: ":0"}
	served := false
	server.serve = func(ln net.Listener) error {
		served = true
		_ = ln.Close()
		return nil
	}
	_ = server.ListenAndServe()
	assert.True(t, served)
}

func TestListenAndServeError(t *testing.T) {
	server := Server{Addr: "localhost:99999"}
	err := server.ListenAndServe()
	assert.Error(t, err)
}

func TestBackwardNoGUID(t *testing.T) {
	mockDB := &test.MockSQLDBManager{}
	testServer := NewServer(mockDB, ":0", "127.0.0.1:0")

	var serverConn net.Conn = &connTester{}
	var destConn net.Conn = &connTester{}

	complete := make(chan string, 1)
	ready := make(chan bool, 1)
	go func() {
		_, _ = destConn.Write([]byte("upstream data"))
		testServer.backward(serverConn, destConn)
		ready <- true
	}()
	<-ready
	go func() {
		buf := make([]byte, 65535)
		n, _ := serverConn.Read(buf)
		if n > 0 {
			complete <- string(buf[:n])
		}
	}()
	result := <-complete
	assert.Equal(t, "upstream data", result)
}

func TestBackwardEOF(t *testing.T) {
	mockDB := &test.MockSQLDBManager{}
	srv := NewServer(mockDB, ":0", "127.0.0.1:0")
	c := &connTester{}
	pr, pw := net.Pipe()
	_ = pw.Close()
	srv.backward(c, pr)
}

func TestHandleConnEndToEnd(t *testing.T) {
	mockDB := &test.MockSQLDBManager{QueryResult: "127.0.0.1"}
	lst, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer func() { _ = lst.Close() }()
	_, port, _ := net.SplitHostPort(lst.Addr().String())
	srv := NewServer(mockDB, ":0", "127.0.0.1:"+port)

	ready := make(chan bool, 1)
	done := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go func() {
		ready <- true
		conn, err := lst.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer func() { _ = conn.Close() }()
		buf := make([]byte, 65535)
		n, _ := conn.Read(buf)
		_, _ = conn.Write([]byte("echo:" + string(buf[:n])))
		done <- struct{}{}
	}()

	<-ready
	client, app := net.Pipe()
	defer func() { _ = client.Close() }()
	defer func() { _ = app.Close() }()
	go srv.handleConn(app)

	req := "GET /x/63f32fee-238e-4f6a-a091-092270d22439 HTTP/1.1\r\n\r\nhello"
	_, _ = client.Write([]byte(req))

	_ = client.SetReadDeadline(time.Now().Add(2 * time.Second))
	rb := make([]byte, 65535)
	n, err := client.Read(rb)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	resp := string(rb[:n])
	assert.True(t, strings.HasPrefix(resp, "echo:"))
	select {
	case <-done:
	case err := <-errCh:
		t.Fatalf("server error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for echo server completion")
	}
}

func TestForwardWriteErrorHandling(t *testing.T) {
	mockDB := &test.MockSQLDBManager{QueryResult: "127.0.0.1"}
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer func() { _ = ln.Close() }()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	srv := NewServer(mockDB, ":0", "127.0.0.1:"+port)
	destChannel := make(chan net.Conn)

	// Accept one connection on the server side
	go func() {
		conn, _ := ln.Accept()
		// keep briefly, then close
		time.Sleep(50 * time.Millisecond)
		_ = conn.Close()
	}()
	// Drain destChannel so forward() doesn't block on sending dst
	go func() {
		if dst := <-destChannel; dst != nil {
			_ = dst.Close()
		}
	}()

	clientConn := &connTester{}
	_, _ = clientConn.Write([]byte("GET /x/63f32fee-238e-4f6a-a091-092270d22439 HTTP/1.1\r\n\r\n"))

	done := make(chan struct{})
	go func() { srv.forward(clientConn, destChannel); close(done) }()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for forward to handle write error")
	}
}

func TestForwardWithGUIDUsesDBInstance(t *testing.T) {
	mockDB := &test.MockSQLDBManager{
		ConnectResult:     &sql.DB{},
		ConnectError:      nil,
		ConnectionStr:     "",
		MPSInstanceResult: "127.0.0.1",
		MPSInstanceError:  nil,
		HealthResult:      false,
		QueryResult:       "127.0.0.1",
	}
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer func() { _ = ln.Close() }()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	testServer := NewServer(mockDB, ":0", "mps:"+port)
	var clientConn net.Conn = &connTester{}

	destChannel := make(chan net.Conn)
	complete := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer func() { _ = conn.Close() }()
		buf := make([]byte, 65535)
		<-destChannel
		n, err := conn.Read(buf)
		if err != nil {
			errCh <- err
			return
		}
		complete <- string(buf[:n])
	}()

	req := "GET /api/v1/amt/log/audit/63f32fee-238e-4f6a-a091-092270d22439?startIndex=0 HTTP/1.1\r\nHost: example\r\n\r\nbody"
	go func() { _, _ = clientConn.Write([]byte(req)); testServer.forward(clientConn, destChannel) }()

	select {
	case got := <-complete:
		assert.Equal(t, req, got)
	case err := <-errCh:
		t.Fatalf("server error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for forwarded data")
	}
}

func TestForwardDialFailureReturnsGracefully(t *testing.T) {
	mockDB := &test.MockSQLDBManager{QueryResult: "127.0.0.1"}
	testServer := NewServer(mockDB, ":0", "127.0.0.1:0")
	var clientConn net.Conn = &connTester{}
	destChannel := make(chan net.Conn)

	_, _ = clientConn.Write([]byte("GET /x/63f32fee-238e-4f6a-a091-092270d22439 HTTP/1.1\r\n\r\n"))
	testServer.forward(clientConn, destChannel)
}

func TestForwardNoGUID_UsesDefaultTarget(t *testing.T) {
	mockDB := &test.MockSQLDBManager{QueryResult: ""}
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer func() { _ = ln.Close() }()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	srv := NewServer(mockDB, ":0", "127.0.0.1:"+port)

	destChannel := make(chan net.Conn)
	got := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer func() { _ = conn.Close() }()
		buf := make([]byte, 65535)
		<-destChannel
		n, err := conn.Read(buf)
		if err != nil {
			errCh <- err
			return
		}
		got <- string(buf[:n])
	}()

	clientConn := &connTester{}
	_, _ = clientConn.Write([]byte("original request"))
	go srv.forward(clientConn, destChannel)

	select {
	case s := <-got:
		assert.Equal(t, "original request", s)
	case err := <-errCh:
		t.Fatalf("server error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for forwarded data without GUID")
	}
}

func TestNewServerDefaultAddr(t *testing.T) {
	mockDB := &test.MockSQLDBManager{}
	s := NewServer(mockDB, "", "target:1234")
	assert.Equal(t, ":8003", s.Addr)
	assert.Equal(t, "target:1234", s.Target)
}
