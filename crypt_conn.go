package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

const cstChallengeSize = 32
const cstChallengeTimeout = 120 * time.Second
const cstHandshakeOK = "OKO"
const cstHandshakeOKSize = 3

type CryptConn struct {
	net.Conn
	cipher    *Cipher
	challenge []byte
	secret    []byte
}

func NewCryptConn(conn net.Conn, cipherMethod string, secret []byte) *CryptConn {
	cipher := NewCipher(cipherMethod, secret)
	return &CryptConn{
		Conn:   conn,
		cipher: cipher,
		secret: secret,
	}
}

func (c *CryptConn) ServerHandshake() (err error) {
	c.Conn.SetDeadline(time.Now().Add(cstChallengeTimeout))
	c.challenge = randomBytes(cstChallengeSize)
	_, err = c.Conn.Write(c.challenge)
	if err != nil {
		return err
	}
	reply := make([]byte, cstChallengeSize)
	_, err = io.ReadFull(c.Conn, reply)
	if err != nil {
		return err
	}
	if !validateChallenge(c.challenge, reply, c.secret) {
		c.Conn.Write([]byte("err"))
		return errors.New("invalid secret")
	}
	_, err = c.Conn.Write([]byte(cstHandshakeOK))
	if err != nil {
		return err
	}
	c.Conn.SetDeadline(time.Time{})
	return nil
}

func (c *CryptConn) ClientHandshake() (err error) {
	c.Conn.SetDeadline(time.Now().Add(cstChallengeTimeout))
	c.challenge = make([]byte, cstChallengeSize)
	_, err = io.ReadFull(c.Conn, c.challenge)
	if err != nil {
		return err
	}
	reply := hmacSha256(c.challenge, c.secret)
	_, err = c.Conn.Write(reply)
	if err != nil {
		return err
	}
	ifOK := make([]byte, cstHandshakeOKSize)
	_, err = io.ReadFull(c.Conn, ifOK)
	if err != nil {
		return err
	}
	if !bytes.Equal(ifOK, []byte(cstHandshakeOK)) {
		return errors.New("invalid secret")
	}
	c.Conn.SetDeadline(time.Time{})
	return nil
}

func (c *CryptConn) Read(b []byte) (int, error) {
	if c.cipher == nil {
		return c.Conn.Read(b)
	}
	n, err := c.Conn.Read(b)
	if n > 0 {
		c.cipher.decrypt(b[0:n], b[0:n])
	}
	return n, err
}

func (c *CryptConn) Write(b []byte) (int, error) {
	if c.cipher == nil {
		return c.Conn.Write(b)
	}
	c.cipher.encrypt(b, b)
	return c.Conn.Write(b)
}

func (c *CryptConn) Close() error {
	return c.Conn.Close()
}
