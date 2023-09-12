package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
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
		if cfg.flgReverse {
			if err := pipe2SvrStdinStdout(cconn); err != nil {
				logf("pipe server console error: %s", err)
			}
			logf("pipe server console done")
		} else if cfg.flgFwdAddr != "" {
			if err := pipe2SvrConn(cfg, cconn); err != nil {
				logf("pipe server connection error: %s", err)
			}
			logf("pipe server connection done")
		} else if cfg.flgCommand != "" {
			if err := pipeCmd2Conn(cfg.flgCommand, cconn); err != nil {
				logf("pipe server command error: %s", err)
			}
			logf("pipe server command done")
		} else {
			if err := pipe2SvrStdout(cconn); err != nil {
				logf("pipe server stdout error: %s", err)
			}
			logf("pipe server stdout done")
		}
	}
}

func pipeCmd2Conn(strCmd string, cconn *CryptConn) (err error) {
	defer cconn.Close()
	exe, args := parseCmd(strCmd)
	cmd := exec.Command(exe, args...)
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

func pipe2SvrStdout(cconn *CryptConn) (err error) {
	defer cconn.Close()
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(os.Stdout, cconn)
		sigs <- err
	}()
	err = <-sigs
	return
}

func pipe2SvrStdinStdout(cconn *CryptConn) (err error) {
	defer cconn.Close()
	sigs := make(chan error, 1)
	go func() {
		_, err := io.Copy(os.Stdout, cconn)
		sigs <- err
	}()
	go func() {
		_, err := io.Copy(cconn, os.Stdin)
		sigs <- err
	}()
	err = <-sigs
	return
}
