package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
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
		if cfg.flgFwdAddr != "" {
			if err := pipe2SvrConn(cfg, cconn); err != nil {
				logf("pipe server connection error: %s", err)
			}
			logf("pipe server connection done")
		} else if cfg.flgCommand != "" {
			if err := pipe2SvrCmd(cfg, cconn); err != nil {
				logf("pipe server command error: %s", err)
			}
			logf("pipe server command done")
		} else {
			if err := pipe2SvrConsole(cconn); err != nil {
				logf("pipe server console error: %s", err)
			}
			logf("pipe server console done")
		}
	}
}

func pipe2SvrCmd(cfg Config, cconn *CryptConn) (err error) {
	defer cconn.Close()
	cmdWithArgs := strings.Split(cfg.flgCommand, " ")
	cmd := exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...)
	cmd.Stdin = cconn
	cmd.Stdout = cconn
	cmd.Stderr = cconn
	sigs := make(chan error, 1)
	go monitorProcess(cmd, sigs)
	go func() {
		err := cmd.Run()
		sigs <- err
	}()
	err = <-sigs
	return
}

func monitorProcess(cmd *exec.Cmd, sigs chan error) {
	for {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			sigs <- nil
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func pipe2SvrConn(cfg Config, cconn *CryptConn) (err error) {
	dst, err := net.Dial(cfg.flgNetwork, cfg.flgFwdAddr)
	if err != nil {
		return
	}
	pipe2(cconn, dst)
	return
}

func pipe2SvrConsole(cconn *CryptConn) (err error) {
	defer cconn.Close()
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(os.Stdout, cconn)
		sigs <- err
	}()
	err = <-sigs
	return
}
