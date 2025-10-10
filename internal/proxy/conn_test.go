/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package proxy

import (
	"io"
	"net"
	"time"
)

type connTester struct {
	deadline time.Time
	buffer   []byte
}

func (c *connTester) Read(b []byte) (n int, err error) {
	copy(b, c.buffer)
	length := len(c.buffer)
	c.buffer = nil
	if length == 0 {
		return 0, io.EOF
	}
	return length, nil
}

func (c *connTester) Write(b []byte) (n int, err error) {
	c.buffer = b
	return len(b), nil
}

func (c *connTester) Close() error {
	return nil
}

func (c *connTester) LocalAddr() net.Addr {
	return nil
}

func (c *connTester) RemoteAddr() net.Addr {
	return nil
}

func (c *connTester) SetDeadline(t time.Time) error {
	c.deadline = t
	return nil
}

func (c *connTester) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *connTester) SetWriteDeadline(t time.Time) error {
	return nil
}
