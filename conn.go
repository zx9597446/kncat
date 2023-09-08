package main

import (
	"net"
	"time"
)

type CryptConn struct {
	conn   net.Conn
	cipher *Cipher
}

func NewCryptConn(conn net.Conn, cipherMethod string, secret []byte) *CryptConn {
	cipher := NewCipher(cipherMethod, secret)
	return &CryptConn{
		conn:   conn,
		cipher: cipher,
	}
}

func (c *CryptConn) Read(b []byte) (int, error) {
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Minute))
	if c.cipher == nil {
		return c.conn.Read(b)
	}
	n, err := c.conn.Read(b)
	if n > 0 {
		c.cipher.decrypt(b[0:n], b[0:n])
	}
	return n, err
}

func (c *CryptConn) Write(b []byte) (int, error) {
	if c.cipher == nil {
		return c.conn.Write(b)
	}
	c.cipher.encrypt(b, b)
	return c.conn.Write(b)
}

func (c *CryptConn) Close() {
	c.conn.Close()
}

func (c *CryptConn) CloseRead() {
	if conn, ok := c.conn.(*net.TCPConn); ok {
		conn.CloseRead()
	}
}

func (c *CryptConn) CloseWrite() {
	if conn, ok := c.conn.(*net.TCPConn); ok {
		conn.CloseWrite()
	}
}