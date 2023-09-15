package main

import (
	"log"
	"net"
	"os"
)

func startServer(cfg Config) {
	ln, err := net.Listen(cfg.flgNetwork, cfg.flgListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			logf("accept error: %v", err)
			continue
		}
		cconn := NewCryptConn(conn, cfg.flgCryptoMethod, []byte(cfg.flgSecretKey))
		if err := cconn.ServerHandshake(); err != nil {
			logf("server handshake error: %v", err)
			cconn.Close()
			continue
		}
		go func() {
			if err := runServer(cfg, cconn); err != nil {
				logf("server run error: %v", err)
			}
		}()
	}
}

func runServer(cfg Config, cconn *CryptConn) (err error) {
	if cfg.flgFwdAddr != "" {
		return serverForwardPort(cfg.flgReverse, cfg.flgNetwork, cfg.flgFwdAddr, cconn)
	}
	//server non reverse: exec cmd mode
	if cfg.flgCommand != "" && !cfg.flgReverse {
		return pipeCmd2Conn(cfg.flgCommand, cconn)
	}
	return pipeStdInOut(os.Stdin, os.Stdout, cconn)
}
