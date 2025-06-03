/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package proxy

import (
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/device-management-toolkit/mps-router/internal/db"
)

// Regular expression to match GUID format
// [a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12} - RFC4122 Search for GUID
// The following guid checks for any uuid/guid format, not following RFC4122 explicitly
var guidRegEx = regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}")

// Server is a TCP server that takes an incoming request and sends it to another
// server, proxying the response back to the client.
type Server struct {
	// TCP address to listen on
	Addr string
	// TCP address of target server
	Target string
	// Database manager
	DB db.Manager
	// Function for serving incoming connections
	serve func(ln net.Listener) error
	// Buffer pool for reducing memory allocations
	bufferPool sync.Pool
}

// NewServer creates a new proxy server with the given address and target
func NewServer(db db.Manager, addr string, target string) *Server {
	if addr == "" {
		addr = ":8003"
	}
	server := &Server{
		Addr:   addr,
		Target: target,
		DB:     db,
		bufferPool: sync.Pool{
			New: func() interface{} {
				buffer := make([]byte, 65535)
				return &buffer
			},
		},
	}
	server.serve = server.serveDefault
	return server
}

// ListenAndServe listens on the TCP network address laddr and then handle packets
// on incoming connections.
func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.serve(listener)
}

// serveDefault is the default serving function that handles incoming connections
func (s *Server) serveDefault(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handleConn(conn)
	}
}

// parseGuid extracts the GUID from the provided content string (the url)
func (s *Server) parseGuid(content string) string {
	idx := strings.IndexByte(content, '\n')
	if idx <= 0 {
		return ""
	}
	return guidRegEx.FindString(content[:idx])
}

// handleConn handles an incoming connection by setting up forward and backward proxies
func (s *Server) handleConn(conn net.Conn) {
	dst, err := s.establishTargetConnection(conn)
	if err != nil {
		log.Printf("Error establishing target connection: %v", err)
		if err := conn.Close(); err != nil {
			log.Printf("Error closing source connection: %v", err)
		}
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Forward direction: client -> target
	go func() {
		defer wg.Done()
		if err := s.proxyData(dst, conn); err != nil && err != io.EOF {
			log.Printf("Forward proxy error: %v", err)
		}
		if tcpConn, ok := dst.(*net.TCPConn); ok {
			tcpConn.CloseWrite()
		}
	}()

	// Backward direction: target -> client
	go func() {
		defer wg.Done()
		if err := s.proxyData(conn, dst); err != nil && err != io.EOF {
			log.Printf("Backward proxy error: %v", err)
		}
		conn.(*net.TCPConn).CloseWrite()
	}()

	wg.Wait()

	if err := conn.Close(); err != nil {
		log.Printf("Error closing source connection: %v", err)
	}
	if err := dst.Close(); err != nil {
		log.Printf("Error closing destination connection: %v", err)
	}
}

// establishTargetConnection reads the initial data from conn, parses GUID, and connects to the target
func (s *Server) establishTargetConnection(conn net.Conn) (net.Conn, error) {
	buffer := s.bufferPool.Get().(*[]byte)
	defer s.bufferPool.Put(buffer)

	n, err := conn.Read(*buffer)
	if err != nil {
		return nil, err
	}

	initialData := (*buffer)[:n]

	// Determine the target destination based on GUID
	destination := s.Target
	guid := s.parseGuid(string(initialData))
	if guid != "" {
		instance := s.DB.Query(guid)
		if instance != "" {
			parts := strings.Split(destination, ":")
			parts[0] = instance
			destination = parts[0] + ":" + parts[1]
		}
	}

	// Connect to target server
	dst, err := net.Dial("tcp", destination)
	if err != nil {
		return nil, err
	}

	// Forward initial data to target
	_, err = dst.Write(initialData)
	if err != nil {
		dst.Close()
		return nil, err
	}

	return dst, nil
}

// proxyData efficiently copies data from src to dst
func (s *Server) proxyData(dst io.Writer, src io.Reader) error {
	buffer := s.bufferPool.Get().(*[]byte)
	defer s.bufferPool.Put(buffer)

	_, err := io.CopyBuffer(dst, src, *buffer)
	return err
}
