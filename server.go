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
			logf("accept error: %s", err)
			continue
		}
		cconn := NewCryptConn(conn, cfg.flgCryptoMethod, []byte(cfg.flgSecretKey))
		if err := cconn.ServerHandshake(); err != nil {
			logf("server handshake error: %s\n", err)
			cconn.Close()
			continue
		}
		go func() {
			if err := runServer(cfg, cconn); err != nil {
				logf("server run error: %s\n", err)
			}
		}()
	}
}

func runServer(cfg Config, cconn *CryptConn) (err error) {
	if cfg.flgFwdAddr != "" {
		return serverForwardPort(cfg.flgReverse, cfg.flgNetwork, cfg.flgFwdAddr, cconn)
	} else if cfg.flgCommand != "" {
		if !cfg.flgReverse {
			return pipeCmd2Conn(cfg.flgCommand, cconn)
		} else {
			return pipeStdInOut(os.Stdin, os.Stdout, cconn)
		}
	} else {
		if !cfg.flgReverse {
			return pipeOut(os.Stdout, cconn)
		} else {
			return pipeIn(os.Stdin, cconn)
		}
	}
}
