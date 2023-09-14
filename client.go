package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func dialAndHandshake(cfg Config) (*CryptConn, error) {
	conn, err := net.Dial(cfg.flgNetwork, cfg.flgConnectAddr)
	if err != nil {
		return nil, err
	}
	cconn := NewCryptConn(conn, cfg.flgCryptoMethod, []byte(cfg.flgSecretKey))
	if err = cconn.ClientHandshake(); err != nil {
		return nil, fmt.Errorf("client handshake error: %s", err)
	}
	return cconn, nil
}

func startClient(cfg Config) {
	cconn, err := dialAndHandshake(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = runClient(cfg, cconn)
	if err != nil {
		logf("run client error: %s", err)
	}
	os.Exit(0)
}

func runClient(cfg Config, cconn *CryptConn) (err error) {
	if cfg.flgFwdAddr != "" {
		return clientForwardPort(cfg.flgReverse, cfg.flgNetwork, cfg.flgFwdAddr, cconn)
	} else if cfg.flgCommand != "" {
		if !cfg.flgReverse {
			return pipeStdInOut(os.Stdin, os.Stdout, cconn)
		} else {
			return pipeCmd2Conn(cfg.flgCommand, cconn)
		}
	} else {
		if !cfg.flgReverse {
			return pipeIn(os.Stdin, cconn)
		} else {
			return pipeOut(os.Stdout, cconn)
		}
	}
}
