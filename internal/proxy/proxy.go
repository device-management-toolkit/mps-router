/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package proxy

import (
	"bytes"
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
}

// NewServer creates a new proxy server with the given address and target
func NewServer(db db.Manager, addr, target string) Server {
	if addr == "" {
		addr = ":8003"
	}
	server := Server{
		Addr:   addr,
		Target: target,
		DB:     db,
	}
	server.serve = server.serveDefault
	return server
}

// ListenAndServe listens on the TCP network address laddr and then handle packets
// on incoming connections.
func (s Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.serve(listener)
}

// serveDefault is the default serving function that handles incoming connections
func (s Server) serveDefault(ln net.Listener) error {
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
func (s Server) parseGuid(content string) string {
	guid := ""
	splitString := strings.Split(content, "\n")
	if len(splitString) < 2 {
		return guid
	}
	guid = guidRegEx.FindString(splitString[0])
	return guid
}

// handleConn handles an incoming connection by setting up forward and backward proxies
func (s Server) handleConn(conn net.Conn) {
	destChannel := make(chan net.Conn, 1)
	defer close(destChannel)

	go s.forward(conn, destChannel)
	dst := <-destChannel
	go s.backward(conn, dst)
}

// forward proxies data from the source connection to the destination server
func (s Server) forward(conn net.Conn, destChannel chan net.Conn) {
	var dst net.Conn
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	var once sync.Once
	buff := make([]byte, 65535)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		b := buff[:n]

		once.Do(func() {
			destination := s.Target
			guid := s.parseGuid(string(b))
			if guid != "" {
				// call to database to get the mps instance
				instance := s.DB.Query(guid)
				if instance != "" {
					parts := strings.Split(destination, ":")
					parts[0] = instance
					destination = parts[0] + ":" + parts[1]
				}
			}
			// connects to target server
			dst, err = net.Dial("tcp", destination)
			if err != nil {
				log.Println(err.Error())
				return
			}
			destChannel <- dst
		})

		if dst == nil {
			return
		}
		_, err = io.Copy(dst, bytes.NewReader(b))
		if err != nil {
			log.Println(err)
			if err := dst.Close(); err != nil {
				log.Printf("Error closing dst: %v", err)
			}
			return
		}
	}
}

// backward proxies data from the destination server back to the source connection
func (s Server) backward(conn, dst net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing conn: %v", err)
		}
		if err := dst.Close(); err != nil {
			log.Printf("Error closing dst: %v", err)
		}
	}()
	_, err := io.Copy(conn, dst)
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}
}
