package main

import (
	"io"
	"net"
)

func serverForwardPort(reverse bool, network, addr string, cconn *CryptConn) (err error) {
	if reverse {
		return listenAndPipe(network, addr, cconn)
	} else {
		return connectAndPipe(network, addr, cconn)
	}
}
func clientForwardPort(reverse bool, network, addr string, cconn *CryptConn) (err error) {
	if reverse {
		return connectAndPipe(network, addr, cconn)
	} else {
		return listenAndPipe(network, addr, cconn)
	}
}

func listenAndPipe(network, addr string, cconn *CryptConn) (err error) {
	ln, err := net.Listen(network, addr)
	if err != nil {
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		if err = pipe2Conn(cconn, conn); err != nil {
			logf("pipe2Conn error: %v", err)
		}
	}
}

func connectAndPipe(network, addr string, cconn *CryptConn) (err error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return
	}
	return pipe2Conn(cconn, conn)
}

func pipe2Conn(dst net.Conn, src net.Conn) (err error) {
	defer func() {
		dst.Close()
		src.Close()
	}()
	logf("pipe to conn: %s -> %s", src.LocalAddr().String(), dst.RemoteAddr().String())
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(dst, src)
		sigs <- err

	}()
	go func() {
		_, err := io.Copy(src, dst)
		sigs <- err
	}()
	err = <-sigs
	logf("pipe to conn done: %v", err)
	return
}
