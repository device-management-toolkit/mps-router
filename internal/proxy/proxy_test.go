/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/

package proxy

import (
	"database/sql"
	"log"
	"net"
	"testing"

	"github.com/open-amt-cloud-toolkit/mps-router/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestServer_parseGuid(t *testing.T) {
	type args struct {
		content string
	}
	server := Server{}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty String", args{content: ""}, ""},
		{"Invalid Guid",
			args{
				content: `GET /api/v1/amt/log/audit/12345?startIndex=0 HTTP/1.1 
							User-Agent: PostmanRuntime/7.28.0
							Accept: */*
							Postman-Token: 63f32fee-238e-4f6a-a091-092270d22439
							Host: 127.0.0.1:8003
							Accept-Encoding: gzip, deflate, br
							Connection: keep-alive`},
			""},
		{"No Guid",
			args{
				content: `GET /api/v1/devices HTTP/1.1
						Host: mpslookup:8003
						Connection: keep-alive
						X-Forwarded-For: 10.0.0.2
						X-Forwarded-Proto: https
						X-Forwarded-Host: localhost
						X-Forwarded-Port: 8443
						X-Forwarded-Path: /mps/api/v1/devices
						X-Forwarded-Prefix: /mps
						X-Real-IP: 10.0.0.2
						sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"
						accept: application/json, text/plain, */*`},
			""},
		{"Valid v4 GUID",
			args{
				content: `GET /api/v1/amt/log/audit/63f32fee-238e-4f6a-a091-092270d22439?startIndex=0 HTTP/1.1 
						User-Agent: PostmanRuntime/7.28.0
						Accept: */*
						Postman-Token: 63f32fee-238e-4f6a-a091-092270d22439
						Host: 127.0.0.1:8003
						Accept-Encoding: gzip, deflate, br
						Connection: keep-alive`},
			"63f32fee-238e-4f6a-a091-092270d22439"},
		{"Valid v1 GUID",
			args{
				content: `GET /api/v1/amt/features/63f32fee-238e-1f6a-a091-092270d22439 HTTP/1.1
						Host: mpsrouter:8003
						Connection: keep-alive
						X-Forwarded-For: 10.0.0.2
						X-Forwarded-Proto: https
						X-Forwarded-Host: localhost
						X-Forwarded-Port: 8443
						X-Forwarded-Path: /mps/api/v1/amt/features/63f32fee-238e-4f6a-a091-092270d22439
						X-Forwarded-Prefix: /mps
						X-Real-IP: 10.0.0.2
						sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="90", "Microsoft 
						Edge";v="90"
						accept: application/json, text/plain, */*
						authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI5RW1SSlRiSWlJYjRiSWVTc21nY1dJanJSNkh5RVRxYyIsImV4cCI6MTYyMjMxMzgzM30.A1LRc_smYHP0Oxghf31EYZnWQ7kaszlqv_8mlMZL9ns
						sec-ch-ua-mobile: ?0
						user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36 
						Edg/90.0.818.66
						sec-fetch-site: same-origin
						sec-fetch-mode: cors
						sec-fetch-dest: empty
						referer: https://localhost/devices/dc09fce2-a602-11ea-90a8-90cf9ac8ee00
						accept-encoding: gzip, deflate, br
						accept-language: en-US,en;q=0.9
						cookie: _ga=GA1.1.1282238231.1618252808; _pk_id.1.1fff=fc8749a5a7b04428.1619127233.; pga4_session=add36de3-b001-429e-8d5b-bd052c9483ba!N3lVSJ57PlRcgWxX3QaZFLSGtDU=; PGADMIN_LANGUAGE=en        
						X-Consumer-ID: 393611c3-aea9-510d-9be4-ac429ecc53f4
						X-Consumer-Username: admin
						X-Credential-Identifier: 9EmRJTbIiIb4bIeSsmgcWIjrR6HyETqc`},
			"63f32fee-238e-1f6a-a091-092270d22439"},
		{"Valid GUID Websocket Request",
			args{
				content: `GET /relay/webrelay.ashx?p=2&host=d12428be-9fa1-4226-9784-54b2038beab6&port=16994&tls=0&tls1only=0 HTTP/1.1
						Host: mpslookup:8003
						Upgrade: websocket
						Connection: keep-alive, Upgrade
						X-Forwarded-For: 10.0.0.2
						X-Forwarded-Proto: https
						X-Forwarded-Host: localhost
						X-Forwarded-Port: 8443
						X-Forwarded-Path: /mps/ws/relay/webrelay.ashx
						X-Forwarded-Prefix: /mps/ws
						X-Real-IP: 10.0.0.2
						Pragma: no-cache
						Cache-Control: no-cache
						User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36`},
			"d12428be-9fa1-4226-9784-54b2038beab6"},
		{"Invalid GUID Websocket Request",
			args{
				content: `GET /relay/webrelay.ashx?p=2&host=d12428be-9fa1-4226-9784&port=16994&tls=0&tls1only=0 HTTP/1.1
						Host: mpslookup:8003
						Upgrade: websocket
						Connection: keep-alive, Upgrade
						X-Forwarded-For: 10.0.0.2
						X-Forwarded-Proto: https
						X-Forwarded-Host: localhost
						X-Forwarded-Port: 8443
						X-Forwarded-Path: /mps/ws/relay/webrelay.ashx
						X-Forwarded-Prefix: /mps/ws
						X-Real-IP: 10.0.0.2
						Pragma: no-cache
						Cache-Control: no-cache
						User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36`},
			""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.parseGuid(tt.args.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListenAndServe(t *testing.T) {
	server := Server{Addr: ":8010"}
	hasBeenServed := false
	server.serve = func(ln net.Listener) error {
		hasBeenServed = true
		return nil
	}
	_ = server.ListenAndServe()

	assert.True(t, hasBeenServed)
}
func TestForwardNoGUID(t *testing.T) {
	mockDB := &test.MockSQLDBManager{
		ConnectResult:     &sql.DB{},
		ConnectError:      nil,
		ConnectionStr:     "",
		MPSInstanceResult: "",
		MPSInstanceError:  nil,
		HealthResult:      false,
		QueryResult:       "",
	}
	testServer := NewServer(mockDB, ":8009", ":3000")
	var server net.Conn = &connTester{}
	destChannel := make(chan net.Conn)
	complete := make(chan string)
	serverReady := make(chan bool)

	go func() {
		ln, err := net.Listen("tcp", ":3000")
		if err != nil {
			log.Fatal(err.Error())
		}
		for {
			serverReady <- true
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}
			buff := make([]byte, 65535)
			<-destChannel
			for {
				n, err := conn.Read(buff)
				if err != nil {
					log.Println(err)
					return
				}
				b := buff[:n]
				defer func() {
					if err := conn.Close(); err != nil {
						log.Printf("Error closing connection: %v", err)
					}
				}()
				if string(b) != "" {
					complete <- string(b)
				}
			}
		}
	}()

	<-serverReady

	go func() {
		_, _ = server.Write([]byte("original request"))
		testServer.forward(server, destChannel)
		println("got the connection")
	}()

	result := <-complete
	assert.Equal(t, "original request", result)

}

func TestBackwardNoGUID(t *testing.T) {
	mockDB := &test.MockSQLDBManager{
		ConnectResult:     &sql.DB{},
		ConnectError:      nil,
		ConnectionStr:     "",
		MPSInstanceResult: "",
		MPSInstanceError:  nil,
		HealthResult:      false,
		QueryResult:       "",
	}
	testServer := NewServer(mockDB, ":8009", ":3000")

	var server net.Conn = &connTester{}
	var destination net.Conn = &connTester{}

	complete := make(chan string)
	ready := make(chan bool)
	go func() {
		_, _ = destination.Write([]byte("upstream data"))
		testServer.backward(server, destination)
		ready <- true
	}()
	<-ready
	go func() {
		for {
			buff := make([]byte, 65535)
			n, err := server.Read(buff)
			if err != nil {
				log.Println(err)
				return
			}
			b := buff[:n]
			if string(b) != "" {
				complete <- string(b)
			}
		}
	}()

	result := <-complete
	assert.Equal(t, "upstream data", result)
}
